package vmess

import (
	"context"

	"github.com/ctykk/go-xray/common/config"
	"github.com/ctykk/go-xray/common/dial_context"
	"github.com/ctykk/go-xray/common/http_proxy"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/vmess"
	"github.com/xtls/xray-core/proxy/vmess/outbound"
)

type Vmess struct {
	Name   string
	Config *config.Config
}

func New(host string, port uint16, cipher Cipher, uuid string, name string) (*Vmess, error) {
	cfg := config.DefaultConfig()
	cfg.CoreConfig.Outbound = []*core.OutboundHandlerConfig{{
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

	return &Vmess{
		Name:   name,
		Config: cfg,
	}, nil
}

func (v *Vmess) DialContext(ctx context.Context) (dial_context.DialContext, error) {
	return dial_context.CommonDialContext(ctx, v.Config)
}

func (v *Vmess) HTTPProxy(ctx context.Context, port uint16) error {
	return http_proxy.HTTPProxy(ctx, v.Config, port)
}
