package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "redis-cluster-proxy-log.msp:6380"
	aArr := strings.Split(s, ":")
	s1 := "sadsadasd"
	aArr1 := strings.Split(s1, ":")
	fmt.Println(aArr, aArr1)
}
