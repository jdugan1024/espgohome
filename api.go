package espgohome

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ESPHomeConnection represents a connection to a device that speaks the ESPHome protocol
type ESPHomeConnection struct {
	Hostname   string
	Port       int
	Password   string
	ClientInfo string
	conn       net.Conn
	reader     *bufio.Reader
	receivers  map[chan proto.Message]map[MessageID]bool
	closed bool
}

func (c *ESPHomeConnection) Dial() error {
	address, err := net.LookupIP(c.Hostname)
	if err != nil {
		return err
	}

	connstr := fmt.Sprintf("%s:%d", address, c.Port)
	conn, err := net.Dial("tcp", connstr)
	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.closed = false
	go c.receiveLoop()

	return nil
}

func (c *ESPHomeConnection) sendMessage(m proto.Message, msgType MessageID) error {
	b, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	log.Printf(">>> SENDING %d %d %x\n", len(b), msgType, b)

	buf := bytes.Buffer{}
	ibuf := make([]byte, binary.MaxVarintLen64)

	buf.Write([]byte{0})
	nb := binary.PutUvarint(ibuf, uint64(len(b)))
	buf.Write(ibuf[:nb])

	nb = binary.PutUvarint(ibuf, uint64(msgType))
	buf.Write(ibuf[:nb])
	buf.Write(b)
	c.conn.Write(buf.Bytes())

	return nil
}

func (c *ESPHomeConnection) receiveLoop() {
	for {
		preamble, err := c.reader.ReadByte()
		if err != nil {
			if c.closed {
				break
			}
			log.Printf("Unable to read preamble: %s", err)
			continue
		}
		if preamble != 0 {
			log.Printf("Invalid preabmle: %b", preamble)
			continue
		}
		size, err := binary.ReadUvarint(c.reader)
		if err != nil {
			log.Printf("unable to decode size: %s", err)
			continue
		}
		msgTypeRaw, err := binary.ReadUvarint(c.reader)
		if err != nil {
			log.Printf("unable to decode message type: %s", err)
			continue
		}
		msgType := MessageID(msgTypeRaw)

		respBytes := make([]byte, size)
		n, err := c.reader.Read(respBytes)
		if n != int(size) {
			log.Printf("didn't read the right number of bytes! %d != %d", n, size)
			continue
		}

		resp, err := decodeMessage(respBytes, msgType)

		for r, filter := range c.receivers {
			ok := filter[msgType]
			if ok {
				r <- resp
			}
		}
	}
}

func (c *ESPHomeConnection) AddReceiver(r chan proto.Message, filters ...MessageID) {
	if c.receivers == nil {
		fmt.Printf("making da map")
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

func (c *ESPHomeConnection) RemoveReceiver(r chan proto.Message) {
	delete(c.receivers, r)
}

func (c *ESPHomeConnection) sendMessageGetResponse(m proto.Message, msgType MessageID, respTypes ...MessageID) chan proto.Message {
	r := make(chan proto.Message)
	c.AddReceiver(r, respTypes...)
	c.sendMessage(m, msgType)
	return r
}

func (c *ESPHomeConnection) Hello() error {
	req := HelloRequest{ClientInfo: c.ClientInfo}
	receiver := c.sendMessageGetResponse(&req, HelloRequestID, HelloResponseID)
	resp := (<-receiver).(*HelloResponse)
	defer c.RemoveReceiver(receiver)

	j, err := protojson.Marshal(resp)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("GOT: %s\n", string(j))

	return nil
}

func (c *ESPHomeConnection) Connect() error {
	req := ConnectRequest{Password: c.Password}
	receiver := c.sendMessageGetResponse(&req, ConnectRequestID, ConnectResponseID)
	resp := (<-receiver).(*ConnectResponse)
	defer c.RemoveReceiver(receiver)

	j, err := protojson.Marshal(resp)
	if err != nil {
		return err
	}

	log.Printf("GOT: %s %v\n", string(j), resp.InvalidPassword)

	return nil
}

func (c *ESPHomeConnection) Disconnect() error {
	req := DisconnectRequest{}
	receiver := c.sendMessageGetResponse(&req, DisconnectRequestID, DisconnectResponseID)
	resp := (<-receiver).(*DisconnectResponse)
	defer c.RemoveReceiver(receiver)

	c.closed = true

	j, err := protojson.Marshal(resp)
	if err != nil {
		return err
	}

	log.Printf("GOT: %s", string(j))

	c.conn.Close()

	return nil
}

func (c *ESPHomeConnection) DeviceInfo() error {
	req := DeviceInfoRequest{}
	receiver := c.sendMessageGetResponse(&req, DeviceInfoRequestID, DeviceInfoResponseID)
	resp := (<-receiver).(*DeviceInfoResponse)
	defer c.RemoveReceiver(receiver)

	j, err := protojson.Marshal(resp)
	if err != nil {
		return err
	}

	log.Printf("GOT: %s\n", string(j))

	return nil
}

func (c *ESPHomeConnection) decodeListEntitesSwitchResponse(respBytes []byte) error {
	resp := ListEntitiesSwitchResponse{}
	err := proto.Unmarshal(respBytes, &resp)
	if err != nil {
		return fmt.Errorf("unable to unmarshal deviceinfo response: %s", err)
	}

	j, err := protojson.Marshal(&resp)
	if err != nil {
		return err
	}

	log.Printf("GOT: %s\n", string(j))

	return nil
}

func (c *ESPHomeConnection) decodeListEntitesBinarySensor(respBytes []byte) error {
	resp := ListEntitiesBinarySensorResponse{}
	err := proto.Unmarshal(respBytes, &resp)
	if err != nil {
		return fmt.Errorf("unable to unmarshal binarysensor response: %s", err)
	}

	j, err := protojson.Marshal(&resp)
	if err != nil {
		return err
	}

	log.Printf("GOT: %s\n", string(j))

	return nil
}

// func (c *ESPHomeConnection) ListEntities() error {
// 	req := ListEntitiesRequest{}
// 	c.sendMessage(&req, ListEntitiesRequestID)

// 	for done := false; done != true; {
// 		respBytes, msgType, err := c.receiveMessage()
// 		if err != nil {
// 			return err
// 		}

// 		switch msgType {
// 		case ListEntitiesDoneResponseID:
// 			done = true
// 		case ListEntitiesSwitchResponseID:
// 			c.decodeListEntitesSwitchResponse(respBytes)
// 		case ListEntitiesBinarySensorResponseID:
// 			c.decodeListEntitesBinarySensor(respBytes)
// 		default:
// 			log.Printf("unsupported entity: %d\n", msgType)
// 		}
// 	}

// 	return nil
// }

func (c *ESPHomeConnection) SwitchCommand(key uint32, state bool) error {
	log.Printf("Switch %d %t\n", key, state)

	req := SwitchCommandRequest{Key: key, State: state}
	err := c.sendMessage(&req, SwitchCommandRequestID)
	if err != nil {
		log.Printf("ERR: %s", err)
	}

	return err
}

func (c *ESPHomeConnection) Ping() error {
	req := PingRequest{}
	receiver := c.sendMessageGetResponse(&req, PingRequestID, PingResponseID)
	resp := (<-receiver).(*PingResponse)
	defer c.RemoveReceiver(receiver)

	j, err := protojson.Marshal(resp)
	if err != nil {
		return err
	}

	log.Printf("GOT: %s\n", string(j))

	return nil
}

func (c *ESPHomeConnection) SubscribeStates() error {
	return nil
}
