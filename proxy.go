package xray

import (
	"context"

	"github.com/ctykk/go-xray/common"
	_ "github.com/xtls/xray-core/main/distro/all"
)

// Proxy is the interface for a proxy service.
type Proxy interface {
	// DialContext returns a function to establish connections through the proxy.
	DialContext(ctx context.Context) (common.DialContext, error)

	// HTTPProxy starts an HTTP proxy server on the given port.
	HTTPProxy(ctx context.Context, port uint16) error
}
