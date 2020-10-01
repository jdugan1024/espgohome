package espgohome

// NEXT:
// get a basic test working, possible approaches:
//  * create mock service that sends specific set of responses
//  * inject a bufio.Reader with packet data? generated from proto bufs + framing

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ErrorClosed indicates that an operation has been attempted on a Connection that has been closed.
var ErrorClosed = errors.New("Connection closed")

// ESPHomeConnection represents a connection to a device that speaks the ESPHome protocol
type ESPHomeConnection struct {
	Password   string
	ClientInfo string
	conn       net.Conn
	reader     *bufio.Reader
	receivers  map[chan proto.Message]map[MessageID]bool
	closed     bool
	Debug      bool
}

// Dial creates a new ESPHomeConnection over TCP
func (c *ESPHomeConnection) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.closed = false
	go c.receiveLoop()

	return nil
}

// Pipe creates a new ESPHomeConnection connected via net.Pipe() to another net.Conn
// this is primarily used for testing
func (c *ESPHomeConnection) Pipe() net.Conn {
	client, server := net.Pipe()

	c.conn = client
	c.reader = bufio.NewReader(c.conn)
	c.closed = false
	go c.receiveLoop()

	return server
}

func encodeMessage(m proto.Message, msgType MessageID) (*bytes.Buffer, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	ibuf := make([]byte, binary.MaxVarintLen64)

	buf.Write([]byte{0})
	nb := binary.PutUvarint(ibuf, uint64(len(b)))
	buf.Write(ibuf[:nb])

	nb = binary.PutUvarint(ibuf, uint64(msgType))
	buf.Write(ibuf[:nb])
	buf.Write(b)

	return &buf, nil
}

func (c *ESPHomeConnection) sendMessage(m proto.Message, msgType MessageID) error {
	buf, err := encodeMessage(m, msgType)
	if c.Debug {
		log.Printf(">>> SENDING %d %d %x\n", buf.Len(), msgType, buf.Bytes())
	}
	if err != nil {
		return err
	}

	c.conn.Write(buf.Bytes())

	return nil
}

func (c *ESPHomeConnection) receiveLoop() {
	for {
		msgType, respBytes, err := receiveMessage(c.reader)

		if err == io.EOF || c.closed {
			c.closed = true
			break
		}

		if err != nil {
			// log.Printf("receive failed: %v", err)
		}

		resp, err := decodeMessage(respBytes, msgType)

		for r, filter := range c.receivers {
			ok := filter[msgType]
			if ok {
				r <- resp
			}
		}
	}

	c.conn.Close()
	for r := range c.receivers {
		close(r)
	}
}

func receiveMessage(r *bufio.Reader) (MessageID, []byte, error) {
	preamble, err := r.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	if preamble != 0 {
		return 0, nil, fmt.Errorf("Invalid preamble: %b", preamble)
	}

	size, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, nil, err
	}

	msgTypeRaw, err := binary.ReadUvarint(r)
	if err != nil {
		return 0, nil, err
	}

	msgType := MessageID(msgTypeRaw)

	respBytes := make([]byte, size)
	n, err := r.Read(respBytes)
	if n != int(size) {
		return 0, nil, fmt.Errorf("didn't read the right number of bytes! %d != %d", n, size)
	}

	return msgType, respBytes, nil
}

// AddReceiver registers a channel used to receive events for given message types
func (c *ESPHomeConnection) AddReceiver(r chan proto.Message, filters ...MessageID) {
	if c.receivers == nil {
		c.receivers = make(map[chan proto.Message]map[MessageID]bool)
	}
	// warn if changing filter?
	for _, f := range filters {
		if c.receivers[r] == nil {
			c.receivers[r] = make(map[MessageID]bool)
		}
		c.receivers[r][f] = true
	}
}

// RemoveReceiver removes a channel from this list of receivers
func (c *ESPHomeConnection) RemoveReceiver(r chan proto.Message) {
	delete(c.receivers, r)
}

func (c *ESPHomeConnection) sendMessageGetResponse(m proto.Message, msgType MessageID, respTypes ...MessageID) (chan proto.Message, error) {
	if c.closed {
		return nil, ErrorClosed
	}
	r := make(chan proto.Message)
	c.AddReceiver(r, respTypes...)
	c.sendMessage(m, msgType)
	return r, nil
}

func (c *ESPHomeConnection) logMessage(name string, msg protoreflect.ProtoMessage) {
	j, err := protojson.Marshal(msg)
	if err != nil {
		// log.Fatal(err)
		return
	}

	log.Printf("%s : msg : %s\n", name, string(j))
}

// Hello sends the Hello message
func (c *ESPHomeConnection) Hello() error {
	req := HelloRequest{ClientInfo: c.ClientInfo}
	receiver, err := c.sendMessageGetResponse(&req, HelloRequestID, HelloResponseID)
	if err != nil {
		return err
	}

	raw, ok := <-receiver
	c.RemoveReceiver(receiver)
	if !ok {
		return ErrorClosed
	}
	resp := raw.(*HelloResponse)

	if c.Debug {
		c.logMessage("Hello", resp)
	}

	return nil
}

// Connect sends the Connect message
func (c *ESPHomeConnection) Connect() error {
	req := ConnectRequest{Password: c.Password}

	receiver, err := c.sendMessageGetResponse(&req, ConnectRequestID, ConnectResponseID)
	if err != nil {
		return err
	}

	raw, ok := <-receiver
	c.RemoveReceiver(receiver)
	if !ok {
		return ErrorClosed
	}

	resp := raw.(*ConnectResponse)
	if c.Debug {
		c.logMessage("Connect", resp)
	}

	if resp.InvalidPassword {
		c.conn.Close()
		c.closed = true
		return fmt.Errorf("invalid password")
	}

	return nil
}

// Disconnect sends the Disconnect message and shuts down
func (c *ESPHomeConnection) Disconnect() error {
	req := DisconnectRequest{}
	receiver, err := c.sendMessageGetResponse(&req, DisconnectRequestID, DisconnectResponseID)
	if err != nil {
		return err
	}
	raw, ok := <-receiver
	c.RemoveReceiver(receiver)
	if !ok {
		return ErrorClosed
	}
	resp := raw.(*DisconnectResponse)

	c.closed = true

	if c.Debug {
		c.logMessage("Disconnect", resp)
	}

	c.conn.Close()

	return nil
}

// DeviceInfo sends the DeviceInfo message
func (c *ESPHomeConnection) DeviceInfo() (*DeviceInfoResponse, error) {
	req := DeviceInfoRequest{}
	receiver, err := c.sendMessageGetResponse(&req, DeviceInfoRequestID, DeviceInfoResponseID)
	if err != nil {
		return nil, err
	}
	raw, ok := <-receiver
	c.RemoveReceiver(receiver)
	if !ok {
		return nil, ErrorClosed
	}
	resp := raw.(*DeviceInfoResponse)

	if c.Debug {
		c.logMessage("DeviceInfo", resp)
	}

	return resp, nil
}

// Entity provides an interface that represents a generic entity
type Entity interface {
	GetKey() uint32
	GetName() string
	GetObjectId() string
}

type XEnt struct {
	Entity
}

type EntityID int32

//go:generate stringer -type=EntityID

const (
	UndefinedEntity EntityID = 0
	BinarySensor    EntityID = 1
	Cover           EntityID = 2
	Fan             EntityID = 3
	Light           EntityID = 4
	Sensor          EntityID = 5
	Switch          EntityID = 6
	TextSensor      EntityID = 7
	Camera          EntityID = 8
	Climate         EntityID = 9
)

func GetEntityType(e Entity) EntityID {
	switch e.(type) {
	case *ListEntitiesBinarySensorResponse:
		return BinarySensor
	case *ListEntitiesCoverResponse:
		return Cover
	case *ListEntitiesFanResponse:
		return Fan
	case *ListEntitiesLightResponse:
		return Light
	case *ListEntitiesSensorResponse:
		return Sensor
	case *ListEntitiesSwitchResponse:
		return Switch
	case *ListEntitiesTextSensorResponse:
		return TextSensor
	case *ListEntitiesCameraResponse:
		return Camera
	case *ListEntitiesClimateResponse:
		return Climate
	default:
		return UndefinedEntity
	}
}

func (c *ESPHomeConnection) ListEntities() ([]Entity, error) {
	req := ListEntitiesRequest{}
	receiver, err := c.sendMessageGetResponse(&req, ListEntitiesRequestID,
		ListEntitiesBinarySensorResponseID,
		ListEntitiesCameraResponseID,
		ListEntitiesClimateResponseID,
		ListEntitiesCoverResponseID,
		ListEntitiesDoneResponseID,
		ListEntitiesFanResponseID,
		ListEntitiesLightResponseID,
		ListEntitiesSensorResponseID,
		ListEntitiesServicesResponseID,
		ListEntitiesSwitchResponseID,
		ListEntitiesTextSensorResponseID)
	if err != nil {
		return []Entity{}, err
	}
	defer c.RemoveReceiver(receiver)

	entities := []Entity{}

	done := false
	for done != true {
		resp, ok := <-receiver
		if !ok {
			return entities, ErrorClosed
		}
		switch m := resp.(type) {
		case *ListEntitiesDoneResponse:
			done = true
			continue
		case Entity:
			entities = append(entities, m)
		default:
			log.Printf("Unsupported message: %t", resp)
		}
	}

	return entities, nil
}

func (c *ESPHomeConnection) SwitchCommand(key uint32, state bool) error {
	req := SwitchCommandRequest{Key: key, State: state}
	err := c.sendMessage(&req, SwitchCommandRequestID)

	return err
}

func (c *ESPHomeConnection) Ping() error {
	req := PingRequest{}
	receiver, err := c.sendMessageGetResponse(&req, PingRequestID, PingResponseID)
	if err != nil {
		return err
	}
	raw, ok := <-receiver
	c.RemoveReceiver(receiver)
	if !ok {
		return ErrorClosed
	}
	resp := raw.(*PingResponse)

	if c.Debug {
		c.logMessage("Ping", resp)
	}

	return nil
}

func (c *ESPHomeConnection) SubscribeStates() (chan protoreflect.ProtoMessage, error) {
	req := SubscribeStatesRequest{}
	receiver := make(chan protoreflect.ProtoMessage)
	c.AddReceiver(receiver,
		BinarySensorStateResponseID,
		CoverStateResponseID,
		FanStateResponseID,
		LightStateResponseID,
		SensorStateResponseID,
		SwitchStateResponseID,
		TextSensorStateResponseID,
		HomeAssistantStateResponseID,
		ClimateStateResponseID)
	err := c.sendMessage(&req, SubscribeStatesRequestID)
	return receiver, err
}

func (c *ESPHomeConnection) SubscribeLogs(level LogLevel) (chan protoreflect.ProtoMessage, error) {
	req := SubscribeLogsRequest{Level: level, DumpConfig: true}
	receiver := make(chan protoreflect.ProtoMessage)
	c.AddReceiver(receiver, SubscribeLogsResponseID)
	err := c.sendMessage(&req, SubscribeLogsRequestID)
	return receiver, err
}
