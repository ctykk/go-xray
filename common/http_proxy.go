package common

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/http"
)

func NewHTTPProxy(ctx context.Context, config *core.Config, port uint16) error {
	// add inbound config
	config.Inbound = []*core.InboundHandlerConfig{{
		ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
			PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(net.Port(port))}},
			Listen:   net.NewIPOrDomain(net.LocalHostIP),
		}),
		ProxySettings: serial.ToTypedMessage(&http.ServerConfig{UserLevel: 0}),
	}}

	// init instance
	instance, err := core.NewWithContext(ctx, config)
	if err != nil {
		return fmt.Errorf("init instance: %w", err)
	}

	// start instance
	err = instance.Start()
	if err != nil {
		return fmt.Errorf("start instance: %w", err)
	}

	// close instance when ctx closed
	go func() {
		<-ctx.Done()
		_ = instance.Close()
	}()

	return nil
}
