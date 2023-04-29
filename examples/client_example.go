package main

import (
	"fmt"
	"qmp"
)

func main() {
	fmt.Println("Connecting to server...")
	client := qmp.NewClient("localhost:4571")

	go func() {
		err := client.Run()
		if err != nil {
			panic(err)
			return
		}
	}()

	for {
		mes := <-client.Messages
		fmt.Println(mes)
	}
}
