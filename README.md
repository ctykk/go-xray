# go-xray

A proxy client library built on [Xray-core](https://github.com/xtls/xray-core).

## Features

- Currently supported proxy types:
  - Shadowsocks
  - Vmess
  - Trojan
  - Vless
- Usage modes:
  - DialContext
  - HTTPProxy

## Installation

```bash
go get github.com/ctykk/go-xray
```

## Usage

### 1. Using DialContext Mode

The `DialContext` method creates a custom dialer function that can be integrated into HTTP clients.

```go
// Create a Shadowsocks node
proxy, _ := shadowsocks.New("your-server.com", 1234, shadowsocks.CipherAES256GCM, "your-password", "proxy-name")

// Get a dialer function
dialer, _ := proxy.DialContext(ctx)

// Make requests through the proxy
resp, _ := resty.New().SetTransport(&http.Transport{DialContext: dialer}).
	R().Get("https://example.com")
```

### 2. Using HTTPProxy Mode

The `HTTPProxy` method starts a local HTTP proxy server that forwards traffic through the proxy.

```go
// Create a Shadowsocks node
proxy, _ := shadowsocks.New("your-server.com", 1234, shadowsocks.CipherAES256GCM, "your-password", "proxy-name")

// Start HTTP proxy on local port 11234
_ = proxy.HTTPProxy(ctx, uint16(11234))

// Make requests through the proxy
resp, _ := resty.New().SetProxy(fmt.Sprintf("http://127.0.0.1:%d", port)).
    R().Get("https://example.com")
```

## Notice

- You must import `_ "github.com/ctykk/go-xray"` in your application, otherwise the proxies will not start properly.

## [MIT Licence](LICENSE)
