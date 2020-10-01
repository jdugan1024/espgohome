package espgohome

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestBasics(t *testing.T) {
	c := ESPHomeConnection{ClientInfo: "test-client"}
	conn := c.Pipe()
	server := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	x := make(chan bool)
	go func(w chan bool) {
		err := c.Hello()
		if err != nil {
			t.Errorf("error in hello: %v", err)
		}
		log.Printf("sent hello\n")
		x <- true
	}(x)

	_, err := server.ReadByte()
	size, err := binary.ReadUvarint(server)
	msgTypeRaw, err := binary.ReadUvarint(server)
	if err != nil {
		log.Printf("unable to decode message type: %s", err)
		panic("decode")
	}
	msgType := MessageID(msgTypeRaw)
	respBytes := make([]byte, size)
	n, err := server.Read(respBytes)
	if n != int(size) {
		log.Printf("didn't read the right number of bytes! %d != %d", n, size)
		panic("read")
	}

	_, err = decodeMessage(respBytes, msgType)
	if err != nil {
		panic(err)
	}

	m := &HelloResponse{ApiVersionMajor: 1, ApiVersionMinor: 3, ServerInfo: "fake-server"}
	buf, err := encodeMessage(m, HelloResponseID)
	if err != nil {
		t.Errorf("error in hello response: %v", err)
	}
	_, err = server.Write(buf.Bytes())
	server.Flush()

	<-x

}

type MockServer struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	t        *testing.T
	Close    chan bool
	closed   bool
	Password string
}

func NewMockServer(conn net.Conn) *MockServer {
	return &MockServer{
		conn:     conn,
		reader:   bufio.NewReader(conn),
		Close:    make(chan bool),
		Password: "********",
	}
}

func (s *MockServer) ReceiveLoop() {
	// TODO: needed?
	go func() {
		<-s.Close
		s.closed = true
	}()

	for {
		msgType, msgBytes, err := receiveMessage(s.reader)
		msg, err := decodeMessage(msgBytes, msgType)
		if err == io.EOF || s.closed {
			s.closed = true
			break
		}

		if err != nil {
			// log.Printf("MockServer error: %v", err)
		}

		// log.Printf("GOT: %s", msgType)
		switch msgType {
		case HelloRequestID:
			s.SendHelloResponse(msg.(*HelloRequest))
		case ConnectRequestID:
			s.SendConnectResponse(msg.(*ConnectRequest))
		default:
			log.Printf("Unsupported message type: %s", msgType)
		}
	}
	s.conn.Close()
}

func (s *MockServer) sendMessage(m proto.Message, msgType MessageID) error {
	buf, err := encodeMessage(m, msgType)
	if err != nil {
		return err
	}

	s.conn.Write(buf.Bytes())

	return nil
}

func (s *MockServer) SendHelloResponse(msg *HelloRequest) {
	resp := &HelloResponse{}
	s.sendMessage(resp, HelloResponseID)
}

func (s *MockServer) SendConnectResponse(msg *ConnectRequest) {
	invalid := msg.Password != s.Password
	resp := &ConnectResponse{InvalidPassword: invalid}
	s.sendMessage(resp, ConnectResponseID)

	if invalid {
		s.conn.Close()
		s.closed = true
	}
}

func TestBadPassword(t *testing.T) {
	client := ESPHomeConnection{ClientInfo: "test-client"}
	serverConn := client.Pipe()
	server := NewMockServer(serverConn)
	go server.ReceiveLoop()

	err := client.Hello()
	if err != nil {
		t.Errorf("hello failed: %v", err)
	}

	err = client.Connect()
	if err == nil {
		t.Error("succeeded with the wrong password!")
	}

	// if the first Connect fails the spec says we must close the connection
	// on both sides immediately
	err = client.Connect()
	if err != ErrorClosed {
		t.Errorf("expected ErrorClosed, got: %v", err)
	}
}

func TestGoodPassword(t *testing.T) {
	client := ESPHomeConnection{ClientInfo: "test-client"}
	serverConn := client.Pipe()
	server := NewMockServer(serverConn)
	go server.ReceiveLoop()

	err := client.Hello()
	if err != nil {
		t.Errorf("hello failed: %v", err)
	}
	log.Printf("Got HelloResponse")

	client.Password = server.Password
	err = client.Connect()
	if err != nil {
		t.Errorf("connect failed: %v", err)
	}

	client.conn.Close()
}

func TestUnexpectedClose(t *testing.T) {
	client := ESPHomeConnection{ClientInfo: "test-client"}
	serverConn := client.Pipe()
	server := NewMockServer(serverConn)
	go server.ReceiveLoop()

	err := client.Hello()
	if err != nil {
		t.Errorf("hello failed: %v", err)
	}

	server.Close <- true

	client.Password = server.Password
	err = client.Connect()
	if err != ErrorClosed {
		t.Errorf("expected ErrorClosed, got %v", err)
	}
}
