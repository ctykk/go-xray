package http_proxy

import (
	"context"
	"fmt"

	"github.com/ctykk/go-xray/common/config"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/http"
)

func HTTPProxy(ctx context.Context, cfg *config.Config, port uint16) error {
	cfg = cfg.Clone()

	cfg.CoreConfig.Inbound = []*core.InboundHandlerConfig{{
		ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
			PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(net.Port(port))}},
			Listen:   net.NewIPOrDomain(net.LocalHostIP),
		}),
		ProxySettings: serial.ToTypedMessage(&http.ServerConfig{}),
	}}

	inst, err := core.NewWithContext(ctx, cfg.CoreConfig)
	if err != nil {
		return fmt.Errorf("new instance: %w", err)
	}

	err = inst.Start()
	if err != nil {
		return fmt.Errorf("start instance: %w", err)
	}

	go func() {
		<-ctx.Done()
		_ = inst.Close()
	}()

	return nil
}
