package main

import (
	"os"

	"github.com/Melidee/gogo/hopeful"
)

func main() {
	example := hopeful.Example()
	example.Apply(os.Args)
}