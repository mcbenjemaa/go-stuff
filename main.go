package main

import (
	"fmt"

	"github.com/mcbenjemaa/go-stuff/rest"
)

func main() {

	rest.HttpGet()
	fmt.Println("##############################\n")
	rest.HttpPost()
	fmt.Println("##############################\n")
	rest.HttpError()

}
