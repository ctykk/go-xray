package vmess

import (
	"context"

	"github.com/ctykk/go-xray/common"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/vmess"
	"github.com/xtls/xray-core/proxy/vmess/outbound"
)

// Node vmess
type Node struct {
	host   string
	port   uint16
	cipher Cipher
	uuid   string

	Name string

	config *core.Config
}

func New(host string, port uint16, cipher Cipher, uuid string, name string) (*Node, error) {
	config := common.NewConfig()
	config.Outbound = []*core.OutboundHandlerConfig{{
		ProxySettings: serial.ToTypedMessage(&outbound.Config{
			Receiver: &protocol.ServerEndpoint{
				Address: net.NewIPOrDomain(net.ParseAddress(host)),
				Port:    uint32(port),
				User: &protocol.User{Account: serial.ToTypedMessage(&vmess.Account{
					Id:               uuid,
					SecuritySettings: &protocol.SecurityConfig{Type: cipher},
				})},
			},
		}),
	}}

	node := Node{
		host:   host,
		port:   port,
		cipher: cipher,
		uuid:   uuid,

		Name: name,

		config: config,
	}

	return &node, nil
}

func (n *Node) DialContext(ctx context.Context) (common.DialContext, error) {
	return common.NewDialContext(ctx, n.config)
}

func (n *Node) HTTPProxy(ctx context.Context, port uint16) error {
	return common.NewHTTPProxy(ctx, n.config, port)
}
