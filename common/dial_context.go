package common

import (
	"context"
	"fmt"
	"net"

	net2 "github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/core"
)

// DialContext is a function type matching net.Dialer.DialContext.
type DialContext = func(ctx context.Context, network, address string) (net.Conn, error)

func NewDialContext(ctx context.Context, config *core.Config) (DialContext, error) {
	instance, err := core.NewWithContext(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("init instance: %w", err)
	}

	dialer := func(ctx context.Context, network, address string) (net.Conn, error) {
		dest, err := net2.ParseDestination(network + ":" + address)
		if err != nil {
			return nil, fmt.Errorf("parse destinaton %s:%s: %w", network, address, err)
		}

		conn, err := core.Dial(ctx, instance, dest)
		if err != nil {
			return nil, fmt.Errorf("dial %s: %w", dest, err)
		}

		return conn, nil
	}

	return dialer, nil
}
