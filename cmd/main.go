package main

import (
	"fmt"

	"github.com/sebastian-nunez/golang-search-engine/config"
)

func main() {
	port := ":" + config.Envs.Port
	fmt.Println("Port -> " + port)
}
