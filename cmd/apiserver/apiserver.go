package main

import (
	"github.com/767829413/normal-frame/internal/apiserver"
)

func main() {
	apiserver.NewApp("apiserver", "config").Run()
}
