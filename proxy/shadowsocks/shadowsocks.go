package shadowsocks

import (
	"context"

	"github.com/ctykk/go-xray/common/config"
	"github.com/ctykk/go-xray/common/dial_context"
	"github.com/ctykk/go-xray/common/http_proxy"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/shadowsocks"
)

type Shadowsocks struct {
	Name   string
	Config *config.Config
}

func New(host string, port uint16, cipher Cipher, password string, name string) (*Shadowsocks, error) {
	cfg := config.DefaultConfig()
	cfg.CoreConfig.Outbound = []*core.OutboundHandlerConfig{{
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

	return &Shadowsocks{
		Name:   name,
		Config: cfg,
	}, nil
}

func (s *Shadowsocks) DialContext(ctx context.Context) (dial_context.DialContext, error) {
	return dial_context.CommonDialContext(ctx, s.Config)
}

func (s *Shadowsocks) HTTPProxy(ctx context.Context, port uint16) error {
	return http_proxy.CommonHTTPProxy(ctx, s.Config, port)
}
