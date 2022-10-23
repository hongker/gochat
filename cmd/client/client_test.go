package client

import (
	"github.com/ebar-go/ego/utils/runtime/signal"
	"github.com/ebar-go/znet/codec"
	"github.com/stretchr/testify/suite"
	"gochat/api"
	"log"
	"net"
	"testing"
	"time"
)

type ClientSuite struct {
	suite.Suite

	conn  net.Conn
	codec codec.Codec

	username string
}

func (suite *ClientSuite) SetupTest() {
	suite.codec = codec.Default()
	var err error
	suite.conn, err = net.Dial("tcp", "localhost:8081")
	suite.Nil(err)

	go suite.receive(512)

	suite.login("foo")
}

func (suite *ClientSuite) receive(maxReadBufferSize int) {
	for {
		buf := make([]byte, maxReadBufferSize)
		n, err := suite.conn.Read(buf)
		if err != nil {
			log.Println("read error: ", err)
			return
		}

		packet := &codec.Packet{}
		err = suite.codec.Unpack(packet, buf[:n])
		if err != nil {
			log.Println("Unpack failed: ", err)
			return
		}
		log.Printf("receive: operate=%d,seq=%d,body=%s", packet.Header.Operate, packet.Header.Seq, string(packet.Body))
	}
}

func (suite *ClientSuite) login(username string) {
	log.Println("login")
	msg, err := suite.encode(api.OperateLogin, api.LoginRequest{Name: username})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)

	time.Sleep(time.Second * 3)

	go func() {
		for {
			time.Sleep(time.Second * 30)
			suite.SendHeartbeat()
		}
	}()

}

func (suite *ClientSuite) SendHeartbeat() {
	msg, err := suite.encode(api.OperateHeartbeat, api.HeartbeatRequest{})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) TestListSession() {
	msg, err := suite.encode(api.OperateListSession, api.SessionListRequest{})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) TestSendMessage() {
	msg, err := suite.encode(api.OperateSendMessage, api.MessageSendRequest{
		Content:     "some message",
		ContentType: "text",
		Target:      "1d96b25d-62d1-43d1-9ae7-8bca5bb6e59c",
	})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) TearDownSuite() {
	log.Println("TearDownSuite")
	<-signal.SetupSignalHandler()
}

func (suite *ClientSuite) AfterTest(_, testName string) {
	log.Println("after: ", testName)
}

func (suite *ClientSuite) TestCreateChannel() {
	msg, err := suite.encode(api.OperateCreateChannel, api.ChannelCreateRequest{Name: "world"})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) encode(operate int16, data any) ([]byte, error) {
	packet := &codec.Packet{Header: codec.Header{Operate: operate, ContentType: codec.ContentTypeJSON, Seq: 1}}

	return suite.codec.Pack(packet, data)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
