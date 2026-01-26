package xray

import (
	"context"

	"github.com/ctykk/go-xray/common"
	_ "github.com/xtls/xray-core/main/distro/all"
)

// Proxy is the interface for a proxy service.
type Proxy interface {
	// DialContext returns a function to establish connections through the proxy.
	//
	// Example:
	//  dialer, _ := proxy.DialContext(ctx)
	//  transport := &http.Transport{DialContext: dialer}
	//  client := &http.Client{Transport: transport}
	//  resp, _ := client.Get("https://www.bing.com")
	DialContext(ctx context.Context) (common.DialContext, error)

	// HTTPProxy starts an HTTP proxy server on the given port.
	//
	// Example:
	//  _ = proxy.HTTPProxy(ctx, 37415)
	//  proxyURL, _ := url.Parse("http://localhost:37415")
	//  transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	//  client := &http.Client{Transport: transport}
	//  resp, _ := client.Get("https://www.bing.com")
	HTTPProxy(ctx context.Context, port uint16) error
}
