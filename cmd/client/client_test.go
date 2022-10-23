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
}

func (suite *ClientSuite) SetupTest() {
	suite.codec = codec.Default()
	var err error
	suite.conn, err = net.Dial("tcp", "localhost:8081")
	suite.Nil(err)

	go suite.receive(signal.SetupSignalHandler())
}

func (suite *ClientSuite) receive(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}

		buf := make([]byte, 512)
		n, err := suite.conn.Read(buf)
		if err != nil {
			return
		}

		log.Println("receive:", string(buf[:n]))
	}
}

func (suite *ClientSuite) TestLogin() {
	msg, err := suite.encode(api.OperateLogin, api.LoginRequest{Name: "foo"})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)

	time.Sleep(time.Second * 5)

	suite.SendHeartbeat()
}

func (suite *ClientSuite) SendHeartbeat() {
	msg, err := suite.encode(api.OperateHeartbeat, api.HeartbeatRequest{})
	suite.Nil(err)

	n, err := suite.conn.Write(msg)
	suite.Nil(err)
	suite.Equal(len(msg), n)

	time.Sleep(time.Second * 5)
}

func (suite *ClientSuite) encode(operate int16, data any) ([]byte, error) {
	packet := &codec.Packet{Header: codec.Header{Operate: operate, ContentType: codec.ContentTypeJSON, Seq: 1}}

	return suite.codec.Pack(packet, data)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
