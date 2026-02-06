package shadowsocks

import (
	"context"

	"github.com/ctykk/go-xray/common"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/shadowsocks"
)

// Node shadowsocks
type Node struct {
	host     string // Server host
	port     uint16 // Server port
	cipher   Cipher // Encryption method
	password string // Encryption password

	Name string // Node display name

	config *core.Config
}

// New creates a Shadowsocks node from raw config values.
func New(host string, port uint16, cipher Cipher, password string, name string) (*Node, error) {
	config := common.NewConfig()
	config.Outbound = []*core.OutboundHandlerConfig{{
		// NOTE: 不添加 Tag 的话，Stats 拿不到 Counter
		Tag: "outbound-tag",
		ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
			Server: &protocol.ServerEndpoint{
				Address: net.NewIPOrDomain(net.ParseAddress(host)),
				Port:    uint32(port),
				User: &protocol.User{Account: serial.ToTypedMessage(&shadowsocks.Account{
					CipherType: cipher,
					Password:   password,
				})},
			},
		}),
	}}

	node := Node{
		host:     host,
		port:     port,
		cipher:   cipher,
		password: password,

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
