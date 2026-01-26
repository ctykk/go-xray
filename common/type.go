package common

import (
	"context"
	"net"
)

// DialContext is a function type matching net.Dialer.DialContext.
type DialContext = func(ctx context.Context, network, address string) (net.Conn, error)
