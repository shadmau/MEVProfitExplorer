package main

import (
	"flag"
	"fmt"
)

func main() {

	startUnixTs := flag.Int("startUnixTs", 0, "blub")
	//entUnixTs 	:= flag.Int()
	flag.Parse()
	fmt.Println(startUnixTs)

}
