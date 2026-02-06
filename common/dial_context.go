package common

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ctykk/go-xray/instance"
	net2 "github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/core"
)

// DialContext is a function type matching net.Dialer.DialContext.
type DialContext = func(ctx context.Context, network, address string) (net.Conn, error)

func NewDialContext(ctx context.Context, config *core.Config) (DialContext, error) {
	inst, err := instance.New(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("init instance: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Second):
				stats, err := inst.Stats()
				fmt.Printf("stats: %+v, err: %+v\n", stats, err)
			}
		}
	}()

	dialer := func(ctx context.Context, network, address string) (net.Conn, error) {
		dest, err := net2.ParseDestination(network + ":" + address)
		if err != nil {
			return nil, fmt.Errorf("parse destinaton %s:%s: %w", network, address, err)
		}

		conn, err := core.Dial(ctx, inst.Inst, dest)
		if err != nil {
			return nil, fmt.Errorf("dial %s: %w", dest, err)
		}

		return conn, nil
	}

	return dialer, nil
}
