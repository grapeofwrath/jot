package main

import (
	"github.com/grapeofwrath/jot/cmd"
	"github.com/grapeofwrath/jot/helper"
)

func main() {
	err := cmd.Execute()
	helper.Check(err)
}
