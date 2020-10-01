package main

import (
	"log"
	"time"

	"github.com/jdugan1024/espgohome"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func states(c chan protoreflect.ProtoMessage) {
	for m := range c {
		switch state := m.(type) {
		case *espgohome.BinarySensorStateResponse:
			log.Printf("Binary sensor %d %t", state.Key, state.State)
		case *espgohome.CoverStateResponse:
			log.Printf("Cover state %d", state.Key)
		case *espgohome.FanStateResponse:
			log.Printf("Fan state")
		case *espgohome.LightStateResponse:
			log.Printf("Light state")
		case *espgohome.SensorStateResponse:
			log.Printf("Sensor state")
		case *espgohome.SwitchStateResponse:
			log.Printf("Switch state %d %t", state.Key, state.State)
		case *espgohome.TextSensorStateResponse:
			log.Printf("TextSensor state")
		case *espgohome.HomeAssistantStateResponse:
			log.Printf("HomeAssistant state")
		case *espgohome.ClimateStateResponse:
			log.Printf("Climate state")
		}
	}
}

func logs(c chan protoreflect.ProtoMessage) {
	for m := range c {
		l := m.(*espgohome.SubscribeLogsResponse)
		log.Printf("L|%s %s %s", l.Level, l.Tag, l.Message)
	}
}
func main() {
	c := &espgohome.ESPHomeConnection{
		ClientInfo: "client-info",
		Password:   "foobar",
		Debug:      true,
	}
	err := c.Dial("10.37.0.187:6053")
	if err != nil {
		log.Fatal(err)
	}
	c.Hello()
	c.Connect()
	info, err := c.DeviceInfo()
	if err != nil {
		panic(err)
	}
	log.Printf("INFO: %v\n", info)
	c.Ping()
	entities, err := c.ListEntities()
	for _, e := range entities {
		log.Printf("|%d %s %s %s\n", e.GetKey(), espgohome.GetEntityType(e), e.GetName(), e.GetObjectId())
	}
	receiver, err := c.SubscribeStates()
	go states(receiver)
	receiver, err = c.SubscribeLogs(espgohome.LogLevel_LOG_LEVEL_VERBOSE)
	go logs(receiver)

	c.SwitchCommand(714259650, true)
	// c.Ping()
	time.Sleep(20 * time.Second)
	// TODO: add pings to keep the connection alive.
	c.SwitchCommand(714259650, false)
	c.Disconnect()
}
