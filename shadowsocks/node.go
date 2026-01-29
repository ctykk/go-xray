package shadowsocks

import (
	"context"
	"fmt"

	"github.com/ctykk/go-xray/common"
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

// Node shadowsocks
type Node struct {
	host     string // Server host
	port     uint16 // Server port
	cipher   Cipher // Encryption method
	password string // Encryption password

	Name string // Node display name
}

// New creates a Shadowsocks node from raw config values.
func New(host string, port uint16, cipher Cipher, password string, name string) (*Node, error) {
	node := Node{
		host:     host,
		port:     port,
		cipher:   cipher,
		password: password,

		Name: name,
	}
	return &node, nil
}

func (n *Node) DialContext(ctx context.Context) (common.DialContext, error) {
	config := core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),

			// disable log
			serial.ToTypedMessage(&log.Config{
				AccessLogType: log.LogType_None,
				ErrorLogType:  log.LogType_None,
			}),
		},
		Outbound: []*core.OutboundHandlerConfig{{
			ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
				Server: &protocol.ServerEndpoint{
					Address: net.NewIPOrDomain(net.ParseAddress(n.host)),
					Port:    uint32(n.port),
					User: &protocol.User{Account: serial.ToTypedMessage(&shadowsocks.Account{
						CipherType: n.cipher,
						Password:   n.password,
					})},
				},
			}),
		}},
	}

	return common.NewDialContext(ctx, &config)
}

func (n *Node) HTTPProxy(ctx context.Context, port uint16) error {
	config := core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),

			// disable log
			serial.ToTypedMessage(&log.Config{
				AccessLogType: log.LogType_None,
				ErrorLogType:  log.LogType_None,
			}),
		},
		Outbound: []*core.OutboundHandlerConfig{{
			ProxySettings: serial.ToTypedMessage(&shadowsocks.ClientConfig{
				Server: &protocol.ServerEndpoint{
					Address: net.NewIPOrDomain(net.ParseAddress(n.host)),
					Port:    uint32(n.port),
					User: &protocol.User{Account: serial.ToTypedMessage(&shadowsocks.Account{
						CipherType: n.cipher,
						Password:   n.password,
					})},
				},
			}),
		}},
		Inbound: []*core.InboundHandlerConfig{{
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortList: &net.PortList{Range: []*net.PortRange{net.SinglePortRange(net.Port(port))}},
				Listen:   net.NewIPOrDomain(net.LocalHostIP),
			}),
			ProxySettings: serial.ToTypedMessage(&http.ServerConfig{UserLevel: 0}),
		}},
	}

	instance, err := core.NewWithContext(ctx, &config)
	if err != nil {
		return fmt.Errorf("init instance: %w", err)
	}

	err = instance.Start()
	if err != nil {
		return fmt.Errorf("start instance: %w", err)
	}
	go func() {
		<-ctx.Done()
		_ = instance.Close()
	}()
	return nil
}
