package main

import (
	"fmt"
	"github.com/ibuymovie/qmp"
	"github.com/ibuymovie/qmp/Message"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("Launching server...")
	server := qmp.NewServer(4571)

	go func() {
		err := server.Run()
		if err != nil {
			panic(err)
			return
		}
	}()
	for {
		timeNow := time.Now()
		fmt.Println("Send message", timeNow)
		server.SendMessageToAll(Message.NewMessage(Message.Json, map[string]interface{}{
			"Hi":   "Hello",
			"Bye":  "Goodbye",
			"time": timeNow,
		}))

		time.Sleep(time.Second * 2)

	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
