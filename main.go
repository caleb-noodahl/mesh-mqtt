package main

import (
	"fmt"
	"log"
	"mesh-mqtt/clients/noaa"
	"mesh-mqtt/clients/radio"
	"mesh-mqtt/config"
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
	conf := config.Config()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	cache, err := pebble.Open(conf.CacheDB, nil)
	if err != nil {
		log.Panic(err)
	}

	rad := radio.NewRadioClient(
		radio.WithPort(conf.RadioPort),
	)

	noaa := noaa.NewNOAAClient()

	server := server.NewServer(
		server.WithPebbleHook(conf.HookDB),
		server.WithPebbleDB(cache),
		server.WithRadioClient(rad),
		server.WithNOAAClient(noaa),
		//todo: setup auth hook
		server.WithOpenAuthHook(),
		server.WithUserHook(),
		server.WithWeatherHook(),
		server.WithDefaultListener(),
	)
	defer server.Close()

	server.Broadcast(fmt.Sprintf("t:%v mtqq server start", time.Now().Unix()), radio.BroadcastOpt{
		Channel: 0,
	})

	go func() {
		err := server.Serve()
		if err != nil {
			server.Broadcast(fmt.Sprintf("t:%v mtqq server shutdown", time.Now().Unix()), radio.BroadcastOpt{
				Channel: 0,
			})
			log.Fatal(err)
		}
	}()

	<-done
}
