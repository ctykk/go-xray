package shadowsocks_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ctykk/go-xray/shadowsocks"
	"github.com/go-resty/resty/v2"
	_ "github.com/xtls/xray-core/main/distro/all"
)

func TestNode_DialContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, err := shadowsocks.New("localhost", 42665, shadowsocks.CipherAES256GCM, "123456", "TestNode_DialContext")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	dialer, err := node.DialContext(ctx)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	client := resty.New().SetTransport(&http.Transport{DialContext: dialer})

	resp, err := client.R().Get("https://bing.com")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("status code: %d", resp.StatusCode())
}

func TestNode_HTTPProxy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	node, err := shadowsocks.New("localhost", 42665, shadowsocks.CipherAES256GCM, "123456", "TestNode_DialContext")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	port := uint16(37415)
	err = node.HTTPProxy(ctx, port)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	client := resty.New().SetProxy(fmt.Sprintf("http://localhost:%d", port))

	resp, err := client.R().Get("https://bing.com")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	t.Logf("status code: %d", resp.StatusCode())
}
