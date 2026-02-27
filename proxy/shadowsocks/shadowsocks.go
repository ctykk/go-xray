package shadowsocks

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/app/dispatcher"
	"github.com/xtls/xray-core/app/log"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/http"
	"github.com/xtls/xray-core/proxy/shadowsocks"
)

type Shadowsocks struct {
	host     string // Server host
	port     uint16 // Server port
	cipher   Cipher // Encryption method
	password string // Encryption password

	Name string // Shadowsocks display name

	ConfigBuilder func() core.Config
}

// New creates a Shadowsocks node from raw config values.
func New(host string, port uint16, cipher Cipher, password string, name string) (*Shadowsocks, error) {
	node := Shadowsocks{
		host:     host,
		port:     port,
		cipher:   cipher,
		password: password,

		Name: name,
	}

	configBuilder := func() core.Config {
		return core.Config{
			App: []*serial.TypedMessage{
				serial.ToTypedMessage(&dispatcher.Config{}),
				serial.ToTypedMessage(&proxyman.InboundConfig{}),
				serial.ToTypedMessage(&proxyman.OutboundConfig{}),

				serial.ToTypedMessage(&log.Config{
					AccessLogType: log.LogType_None,
					ErrorLogType:  log.LogType_None,
				}),
			},

			Outbound: []*core.OutboundHandlerConfig{{
				ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
					Server: &protocol.ServerEndpoint{
						Address: net.NewIPOrDomain(net.ParseAddress(node.host)),
						Port:    uint32(node.port),
						User: &protocol.User{Account: serial.ToTypedMessage(&shadowsocks.Account{
							CipherType: node.cipher,
							Password:   node.password,
						})},
					},
				}),
			}},
		}

	}
	node.ConfigBuilder = configBuilder

	return &node, nil
}

func (n *Shadowsocks) DialContext(ctx context.Context) (func(context.Context, string, string) (net.Conn, error), error) {
	cfg := n.ConfigBuilder()

	inst, err := core.NewWithContext(ctx, &cfg)
	if err != nil {
		return nil, fmt.Errorf("new instance: %w", err)
	}

	dc := func(ctx context.Context, network, addr string) (net.Conn, error) {
		dest, err := net.ParseDestination(network + ":" + addr)
		if err != nil {
			return nil, err
		}

		conn, err := core.Dial(ctx, inst, dest)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	return dc, nil
}

func (n *Shadowsocks) HTTPProxy(ctx context.Context, port uint16) error {
	cfg := n.ConfigBuilder()

	cfg.Inbound = []*core.InboundHandlerConfig{{
		ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
			PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(net.Port(port))}},
			Listen:   net.NewIPOrDomain(net.LocalHostIP),
		}),
		ProxySettings: serial.ToTypedMessage(&http.ServerConfig{}),
	}}

	inst, err := core.NewWithContext(ctx, &cfg)
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
