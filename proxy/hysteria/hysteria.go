package hysteria

import (
	"context"

	"github.com/ctykk/go-xray/common/config"
	"github.com/ctykk/go-xray/common/dial_context"
	"github.com/ctykk/go-xray/common/http_proxy"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/hysteria"
	"github.com/xtls/xray-core/proxy/hysteria/account"
)

type Hysteria struct {
	Name   string
	Config *config.Config
}

func New(host string, port uint16, password string, name string) (*Hysteria, error) {
	// fixme: support mihomo: ports mport udp skip-cert-verify sni
	cfg := config.DefaultConfig()
	cfg.CoreConfig.Outbound = []*core.OutboundHandlerConfig{{
		ProxySettings: serial.ToTypedMessage(&hysteria.ClientConfig{
			Server: &protocol.ServerEndpoint{
				Address: net.NewIPOrDomain(net.ParseAddress(host)),
				Port:    uint32(port),
				User: &protocol.User{Account: serial.ToTypedMessage(&account.Account{
					Auth: password,
				})},
			},
			Version: 2,
		}),
	}}

	return &Hysteria{
		Name:   name,
		Config: cfg,
	}, nil
}

func (h *Hysteria) DialContext(ctx context.Context) (dial_context.DialContext, error) {
	return dial_context.CommonDialContext(ctx, h.Config)
}

func (h *Hysteria) HTTPProxy(ctx context.Context, port uint16) error {
	return http_proxy.CommonHTTPProxy(ctx, h.Config, port)
}
