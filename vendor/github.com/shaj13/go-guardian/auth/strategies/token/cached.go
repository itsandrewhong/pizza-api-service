package token

import (
	"context"
	"net/http"

	"github.com/shaj13/go-guardian/auth"
	"github.com/shaj13/go-guardian/errors"
	"github.com/shaj13/go-guardian/store"
)

// CachedStrategyKey export identifier for the cached bearer strategy,
// commonly used when enable/add strategy to go-guardian authenticator.
const CachedStrategyKey = auth.StrategyKey("Token.Cached.Strategy")

// AuthenticateFunc declare custom function to authenticate request using token.
// The authenticate function invoked by Authenticate Strategy method when
// The token does not exist in the cahce and the invocation result will be cached, unless an error returned.
// Use NoOpAuthenticate instead to refresh/mangae token directly using cache or Append function.
type AuthenticateFunc func(ctx context.Context, r *http.Request, token string) (auth.Info, error)

// New return new auth.Strategy.
// The returned strategy, caches the invocation result of authenticate function, See AuthenticateFunc.
// Use NoOpAuthenticate to refresh/mangae token directly using cache or Append function, See NoOpAuthenticate.
func New(auth AuthenticateFunc, c store.Cache, opts ...auth.Option) auth.Strategy {
	if auth == nil {
		panic("Authenticate Function required and can't be nil")
	}

	if c == nil {
		panic("Cache object required and can't be nil")
	}

	cached := &cachedToken{
		authFunc: auth,
		cache:    c,
		typ:      Bearer,
		parser:   AuthorizationParser(string(Bearer)),
	}

	for _, opt := range opts {
		opt.Apply(cached)
	}

	return cached
}

type cachedToken struct {
	parser   Parser
	typ      Type
	cache    store.Cache
	authFunc AuthenticateFunc
}

func (c *cachedToken) Authenticate(ctx context.Context, r *http.Request) (auth.Info, error) {
	token, err := c.parser.Token(r)
	if err != nil {
		return nil, err
	}

	info, ok, err := c.cache.Load(token, r)

	if err != nil {
		return nil, err
	}

	// if token not found invoke user authenticate function
	if !ok {
		info, err = c.authFunc(ctx, r, token)
		if err == nil {
			// cache result
			err = c.cache.Store(token, info, r)
		}
	}

	if err != nil {
		return nil, err
	}

	if _, ok := info.(auth.Info); !ok {
		return nil, errors.NewInvalidType((*auth.Info)(nil), info)
	}

	return info.(auth.Info), nil
}

func (c *cachedToken) Append(token string, info auth.Info, r *http.Request) error {
	return c.cache.Store(token, info, r)
}

func (c *cachedToken) Revoke(token string, r *http.Request) error {
	return c.cache.Delete(token, r)
}

func (c *cachedToken) Challenge(realm string) string { return challenge(realm, c.typ) }

// NoOpAuthenticate implements Authenticate function, it return nil, auth.ErrNOOP,
// commonly used when token refreshed/mangaed directly using cache or Append function,
// and there is no need to parse token and authenticate request.
func NoOpAuthenticate(ctx context.Context, r *http.Request, token string) (auth.Info, error) {
	return nil, auth.ErrNOOP
}
