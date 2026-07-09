package main

import (
	"os"

	"github.com/Melidee/gogo/cli"
)

func main() {
	example := cli.Example()
	example.Apply(os.Args)
}
