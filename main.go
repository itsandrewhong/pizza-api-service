package main

func main() {
	a := App{}
	port := "8000"

	a.Initialize()

	a.Run(":" + port)
}
