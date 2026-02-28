package vless

import (
	"context"

	"github.com/ctykk/go-xray/common/config"
	"github.com/ctykk/go-xray/common/dial_context"
	"github.com/ctykk/go-xray/common/http_proxy"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/vless"
	"github.com/xtls/xray-core/proxy/vless/outbound"
)

type Vless struct {
	Name   string
	Config *config.Config
}

func New(host string, port uint16, uuid string, encryption string, name string) (*Vless, error) {
	cfg := config.DefaultConfig()
	cfg.CoreConfig.Outbound = []*core.OutboundHandlerConfig{{
		ProxySettings: serial.ToTypedMessage(&outbound.Config{
			Vnext: &protocol.ServerEndpoint{
				Address: net.NewIPOrDomain(net.ParseAddress(host)),
				Port:    uint32(port),
				User: &protocol.User{Account: serial.ToTypedMessage(&vless.Account{
					Id:         uuid,
					Encryption: encryption,
				})},
			},
		}),
	}}

	return &Vless{
		Name:   name,
		Config: cfg,
	}, nil
}

func (n *Vless) DialContext(ctx context.Context) (dial_context.DialContext, error) {
	return dial_context.CommonDialContext(ctx, n.Config)
}

func (n *Vless) HTTPProxy(ctx context.Context, port uint16) error {
	return http_proxy.CommonHTTPProxy(ctx, n.Config, port)
}
