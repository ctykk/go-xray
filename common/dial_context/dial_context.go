package dial_context

import (
	"context"
	"fmt"

	"github.com/ctykk/go-xray/common/config"
	"github.com/xtls/xray-core/common/net"
	"github.com/xtls/xray-core/core"
)

func DialContext(ctx context.Context, cfg *config.Config) (func(context.Context, string, string) (net.Conn, error), error) {
	cfg = cfg.Clone()

	inst, err := core.NewWithContext(ctx, cfg.CoreConfig)
	if err != nil {
		return nil, fmt.Errorf("new instance: %w", err)
	}

	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
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

	go func() {
		<-ctx.Done()
		_ = inst.Close()
	}()

	return dialer, nil
}
