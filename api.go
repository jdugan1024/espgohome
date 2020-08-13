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
	receivers  map[chan proto.Message]map[uint64]bool
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
	go c.receiveLoop()

	return nil
}

func (c *ESPHomeConnection) sendMessage(m proto.Message, msgType uint64) error {
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
	// TODO: handle shutdown probably via a shutdown channel
	for {
		preamble, err := c.reader.ReadByte()
		if err != nil {
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
		msgType, err := binary.ReadUvarint(c.reader)
		if err != nil {
			log.Printf("unable to decode message type: %s", err)
			continue
		}

		respBytes := make([]byte, size)
		n, err := c.reader.Read(respBytes)
		if n != int(size) {
			log.Printf("didn't read the right number of bytes! %d != %d", n, size)
			continue
		}

		resp, err := c.decodeMessage(respBytes, msgType)

		for r, filter := range c.receivers {
			ok := filter[msgType]
			if ok {
				r <- resp
			}
		}
	}
}

// TODO: generate this?
func (c *ESPHomeConnection) decodeMessage(raw []byte, msgType uint64) (proto.Message, error) {
	log.Printf("decode %d", msgType)
	switch msgType {
	case HelloResponse_ID:
		resp := &HelloResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ConnectResponse_ID:
		resp := &ConnectResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case DeviceInfoResponse_ID:
		resp := &DeviceInfoResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesDoneResponse_ID:
		resp := &ListEntitiesDoneResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesSwitchResponse_ID:
		resp := &ListEntitiesSwitchResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case ListEntitiesBinarySensorResponse_ID:
		resp := &ListEntitiesBinarySensorResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case PingResponse_ID:
		resp := &PingResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	case DisconnectResponse_ID:
		resp := &DisconnectResponse{}
		err := proto.Unmarshal(raw, resp)
		return resp, err
	default:
		err := fmt.Errorf("unsupported message: %d", msgType)
		return nil, err
	}
}

func (c *ESPHomeConnection) AddReceiver(r chan proto.Message, filters []uint64) {
	if c.receivers == nil {
		fmt.Printf("making da map")
		c.receivers = make(map[chan proto.Message]map[uint64]bool)
	}
	// warn if changing filter?
	for _, f := range filters {
		if c.receivers[r] == nil {
			c.receivers[r] = make(map[uint64]bool)
		}
		c.receivers[r][f] = true
	}
}

func (c *ESPHomeConnection) RemoveReceiver(r chan proto.Message) {
	delete(c.receivers, r)
}

func (c *ESPHomeConnection) sendMessageGetResponse(m proto.Message, msgType, respType uint64) chan proto.Message {
	r := make(chan proto.Message)
	c.AddReceiver(r, []uint64{respType})
	c.sendMessage(m, msgType)
	return r
}

func (c *ESPHomeConnection) Hello() error {
	req := HelloRequest{ClientInfo: c.ClientInfo}
	receiver := c.sendMessageGetResponse(&req, HelloRequest_ID, HelloResponse_ID)
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
	receiver := c.sendMessageGetResponse(&req, ConnectRequest_ID, ConnectResponse_ID)
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
	receiver := c.sendMessageGetResponse(&req, DisconnectRequest_ID, DisconnectResponse_ID)
	resp := (<-receiver).(*DisconnectResponse)
	defer c.RemoveReceiver(receiver)

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
	receiver := c.sendMessageGetResponse(&req, DeviceInfoRequest_ID, DeviceInfoResponse_ID)
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
// 	c.sendMessage(&req, ListEntitiesRequest_ID)

// 	for done := false; done != true; {
// 		respBytes, msgType, err := c.receiveMessage()
// 		if err != nil {
// 			return err
// 		}

// 		switch msgType {
// 		case ListEntitiesDoneResponse_ID:
// 			done = true
// 		case ListEntitiesSwitchResponse_ID:
// 			c.decodeListEntitesSwitchResponse(respBytes)
// 		case ListEntitiesBinarySensorResponse_ID:
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
	err := c.sendMessage(&req, SwitchCommandRequest_ID)
	if err != nil {
		log.Printf("ERR: %s", err)
	}

	return err
}

func (c *ESPHomeConnection) Ping() error {
	req := PingRequest{}
	receiver := c.sendMessageGetResponse(&req, PingRequest_ID, PingResponse_ID)
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
