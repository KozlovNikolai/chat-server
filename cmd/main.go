package main

import (
	"fmt"
	"log"
)

func main() {
	//n, err := fmt.Println(color.GreenString("Hello, world!"))
	n, err := fmt.Println("Hello, world!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(n)
}
