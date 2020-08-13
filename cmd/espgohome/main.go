package main

import (
	"log"
	"time"

	"github.com/jdugan1024/espgohome"
)

func main() {
	c := &espgohome.ESPHomeConnection{
		ClientInfo: "esphomego",
		Hostname:   "10.37.0.196",
		Password:   "foobar",
		Port:       6053}
	err := c.Dial()
	if err != nil {
		log.Fatal(err)
	}
	c.Hello()
	c.Connect()
	c.DeviceInfo()
	c.Ping()
	//c.ListEntities()
	c.SwitchCommand(714259650, true)
	// c.Ping()
	time.Sleep(2 * time.Second)
	// TODO: add pings to keep the connection alive.
	c.SwitchCommand(714259650, false)
	c.Disconnect()
}
