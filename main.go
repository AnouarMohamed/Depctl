package main

import (
	_ "embed"

	"github.com/AnouarMohamed/Depctl/cmd"
)

//go:embed HHHQ
var banner string

func main() {
	cmd.SetBanner(banner)
	cmd.Execute()
}
