package main

import (
	"math/rand"
	"time"

	"github.com/767829413/normal-frame/internal/apiserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	apiserver.NewApp("apiserver", "config").Run()
}
