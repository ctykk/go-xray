package trojan

import (
	"context"

	"github.com/ctykk/go-xray/common/config"
	"github.com/ctykk/go-xray/common/dial_context"
	"github.com/ctykk/go-xray/common/http_proxy"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/trojan"
)

type Trojan struct {
	Name   string
	Config *config.Config
}

func New(host string, port uint16, password string, name string) (*Trojan, error) {
	cfg := config.DefaultConfig()
	cfg.CoreConfig.Outbound = []*core.OutboundHandlerConfig{{
		ProxySettings: serial.ToTypedMessage(&trojan.ClientConfig{
			Server: &protocol.ServerEndpoint{
				Address: net.NewIPOrDomain(net.ParseAddress(host)),
				Port:    uint32(port),
				User: &protocol.User{Account: serial.ToTypedMessage(&trojan.Account{
					Password: password,
				})},
			},
		}),
	}}

	return &Trojan{
		Name:   name,
		Config: cfg,
	}, nil
}

func (n *Trojan) DialContext(ctx context.Context) (func(context.Context, string, string) (net.Conn, error), error) {
	return dial_context.DialContext(ctx, n.Config)
}

func (n *Trojan) HTTPProxy(ctx context.Context, port uint16) error {
	return http_proxy.HTTPProxy(ctx, n.Config, port)
}
