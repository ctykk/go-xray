package shadowsocks

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/ctykk/go-xray/common"
	"github.com/xtls/xray-core/app/dispatcher"
	"github.com/xtls/xray-core/app/log"
	"github.com/xtls/xray-core/app/proxyman"
	net2 "github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/common/protocol"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"github.com/xtls/xray-core/proxy/http"
	"github.com/xtls/xray-core/proxy/shadowsocks"
)

type Node struct {
	host       string                 // Server host
	port       uint16                 // Server port
	cipherType shadowsocks.CipherType // Encryption method
	password   string                 // Encryption password

	Name string // Node display name
}

const (
	CipherNameAES128GCM = "AES-128-GCM"
	CipherNameAES256GCM = "AES-256-GCM"
)

// New creates a Shadowsocks node from raw config values.
func New(host string, port uint16, cipherType string, password string, name string) (*Node, error) {
	// cipher type
	var ct shadowsocks.CipherType
	switch strings.ToUpper(cipherType) {
	case CipherNameAES128GCM:
		ct = shadowsocks.CipherType_AES_128_GCM
	case CipherNameAES256GCM:
		ct = shadowsocks.CipherType_AES_256_GCM
	default:
		return nil, fmt.Errorf("unknown cipher type: %s", cipherType)
	}

	node := Node{
		host:       host,
		port:       port,
		cipherType: ct,
		password:   password,

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
					Address: net2.NewIPOrDomain(net2.ParseAddress(n.host)),
					Port:    uint32(n.port),
					User: &protocol.User{Account: serial.ToTypedMessage(&shadowsocks.Account{
						CipherType: n.cipherType,
						Password:   n.password,
					})},
				},
			}),
		}},
	}

	instance, err := core.NewWithContext(ctx, &config)
	if err != nil {
		return nil, fmt.Errorf("init instance: %w", err)
	}

	dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
		dest, err := net2.ParseDestination(network + ":" + addr)
		if err != nil {
			return nil, err
		}
		return core.Dial(ctx, instance, dest)
	}
	return dialContext, nil
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
					Address: net2.NewIPOrDomain(net2.ParseAddress(n.host)),
					Port:    uint32(n.port),
					User: &protocol.User{Account: serial.ToTypedMessage(&shadowsocks.Account{
						CipherType: n.cipherType,
						Password:   n.password,
					})},
				},
			}),
		}},
		Inbound: []*core.InboundHandlerConfig{{
			ReceiverSettings: serial.ToTypedMessage(&proxyman.ReceiverConfig{
				PortList: &net2.PortList{Range: []*net2.PortRange{net2.SinglePortRange(net2.Port(port))}},
				Listen:   net2.NewIPOrDomain(net2.LocalHostIP),
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
