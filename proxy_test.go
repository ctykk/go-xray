package xray_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/ctykk/go-xray"
	"github.com/ctykk/go-xray/proxy/shadowsocks"
	"github.com/ctykk/go-xray/proxy/trojan"
	"github.com/ctykk/go-xray/proxy/vmess"
	"github.com/go-resty/resty/v2"
)

type ProxyPort struct {
	Proxy xray.Proxy
	Name  string
}

func TestDialContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pps := make([]ProxyPort, 0)

	// shadowsocks
	ss, err := shadowsocks.New("localhost", 39571, shadowsocks.CipherAES256GCM, "123456", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: ss, Name: "shadowsocks"})

	// trojan
	tr, err := trojan.New("localhost", 39572, "123456", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: tr, Name: "trojan"})

	// vmess
	vm, err := vmess.New("localhost", 39574, vmess.CipherAES128GCM, "433d0d23-d443-5212-a0df-cdbdd2d53015", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: vm, Name: "vmess"})

	for _, p := range pps {
		dialer, err := p.Proxy.DialContext(ctx)
		if err != nil {
			t.Errorf("%s: %+v", p.Name, err)
		}

		client := resty.New().SetTransport(&http.Transport{DialContext: dialer})

		_, err = client.R().Get("https://bing.com")
		if err != nil {
			t.Errorf("%s: %+v", p.Name, err)
		}
		t.Logf("%s: success", p.Name)
	}
}

func TestHTTPProxy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pps := make([]ProxyPort, 0)

	// shadowsocks
	ss, err := shadowsocks.New("localhost", 39571, shadowsocks.CipherAES256GCM, "123456", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: ss, Name: "shadowsocks"})

	// trojan
	tr, err := trojan.New("localhost", 39572, "123456", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: tr, Name: "trojan"})

	// vmess
	vm, err := vmess.New("localhost", 39574, vmess.CipherAES128GCM, "433d0d23-d443-5212-a0df-cdbdd2d53015", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: vm, Name: "vmess"})

	for _, p := range pps {
		port := uint16(10000 + rand.UintN(55535))
		err = p.Proxy.HTTPProxy(ctx, port)
		if err != nil {
			t.Errorf("%s: %+v", p.Name, err)
		}

		client := resty.New().SetProxy(fmt.Sprintf("http://localhost:%d", port))

		_, err = client.R().Get("https://bing.com")
		if err != nil {
			t.Errorf("%s: %+v", p.Name, err)
		}
		t.Logf("%s: success", p.Name)
	}
}
