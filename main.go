package main

import (
	"fmt"

	util "fnet.cz/x/generics/internal"
)

func stackDemo() {
	var s = []string{"one"}
	util.Push(&s, "two")
	util.Push(&s, "three")
	fmt.Println(s)
	for i := 1; i <= 5; i++ {
		val, ok := util.Pop(&s)
		fmt.Println(val, ok)
		fmt.Println(s)
	}
}

func main() {
	fmt.Println("Hello world")

	util.PubSubDemo()
}
