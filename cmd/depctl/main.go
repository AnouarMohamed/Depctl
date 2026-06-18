package main

import (
	"github.com/AnouarMohamed/Depctl/cmd"
	"github.com/AnouarMohamed/Depctl/internal/brand"
)

func main() {
	cmd.SetBanner(brand.Banner())
	cmd.Execute()
}
