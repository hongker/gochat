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

	username  string
	channelID string
}

func (suite *ClientSuite) SetupTest() {
	suite.codec = codec.Default()
	var err error
	suite.conn, err = net.Dial("tcp", "localhost:8081")
	suite.Nil(err)

	go suite.receive(512)

	suite.login("foo")
	suite.channelID = "91df11cc-6cd6-4009-a855-fb3603568582"

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

	time.Sleep(time.Second * 1)

	go func() {
		for {
			time.Sleep(time.Second * 10)
			suite.SendHeartbeat()
		}
	}()

}

func (suite *ClientSuite) SendHeartbeat() {
	log.Println("heartbeat")
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
		Target:      "f2a6a816-f6ed-403f-8db8-16ef279cfd39",
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

func (suite *ClientSuite) TestCreateChannel() {
	msg, err := suite.encode(api.OperateCreateChannel, api.ChannelCreateRequest{Name: "world"})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) TestJoinChannel() {
	log.Println("TestJoinChannel")
	msg, err := suite.encode(api.OperateJoinChannel, api.ChannelLeaveRequest{ID: suite.channelID})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) TestLeaveChannel() {
	log.Println("TestLeaveChannel")
	msg, err := suite.encode(api.OperateLeaveChannel, api.ChannelLeaveRequest{ID: suite.channelID})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) TestBroadcastChannel() {
	log.Println("TestBroadcastChannel")
	suite.TestJoinChannel()
	msg, err := suite.encode(api.OperateBroadcastChannel, api.ChannelBroadcastRequest{
		Content:     "some message",
		ContentType: "text",
		Target:      suite.channelID,
	})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	log.Printf("write message:%s, n=%d\n", string(msg), n)
	suite.Nil(err)
	suite.Equal(len(msg), n)
}

func (suite *ClientSuite) TestQueryMessage() {
	msg, err := suite.encode(api.OperateQueryMessage, api.MessageQueryRequest{
		SessionID: suite.channelID,
	})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	log.Printf("write message:%s, n=%d\n", string(msg), n)
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
