package main

import (
	"fmt"
	"log"
	"mesh-mqtt/clients/noaa"
	"mesh-mqtt/clients/radio"
	"mesh-mqtt/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cockroachdb/pebble"
)

func main() {
	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	serverDB, err := pebble.Open("server.db", nil)
	if err != nil {
		log.Panic(err)
	}

	rad := radio.NewRadioClient(
		radio.WithPort("/dev/tty.usbmodemF412FA6DC03C1"),
	)

	noaa := noaa.NewNOAAClient()

	server := server.NewServer(
		server.WithPebbleHook("cache.db"),
		server.WithPebbleDB(serverDB),
		server.WithRadioClient(rad),
		server.WithNOAAClient(noaa),
		server.WithOpenAuthHook(),
		server.WithUserHook(),
		server.WithWeatherHook(),
		server.WithDefaultListener(),
	)
	defer server.Close()

	server.Broadcast(fmt.Sprintf("t:%v mtqq server start", time.Now().Unix()), radio.BroadcastOpt{
		Channel: 2,
	})

	go func() {
		err := server.Serve()
		if err != nil {
			server.Broadcast(fmt.Sprintf("t:%v mtqq server shutdown", time.Now().Unix()), radio.BroadcastOpt{
				Channel: 2,
			})
			log.Fatal(err)
		}
	}()

	<-done
}
