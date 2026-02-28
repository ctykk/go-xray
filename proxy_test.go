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
	"github.com/ctykk/go-xray/proxy/vless"
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

	// vless
	// vl, err := vless.New("localhost", 39573, "433d0d23-d443-5212-a0df-cdbdd2d53015", "none", "")
	vl, err := vless.New("localhost", 39573, "433d0d23-d443-5212-a0df-cdbdd2d53015", "kjiwehh2tmsGrdtPjTPG03qPAEJqLeAJR4ClA_G4J0QAz7NyY2o0paRJo4B14TNavpIlbeUPjFRHIkhebjBOKxU8v3shXPUrHdha9ATNg2ILi1u16uS27JVBW6umvLh5ulqxq4JfFiR8qfY45WPEPIIAb4lpORkcqWJQmiMDfoxA_ccx0lnI5dWysiZHL6QvETiIKuhwwgYxCeyTxuFCRlBRwXnCPdsEg1THxaa1H6iB8hiXDMK0p9wz__FenZZsVAFTwBJ538pVFWeluLePCncJD0JEY8gKF2BF8bAzODwFwtM9zrlZciaXc2K9QHRS_GCc0EcWbjoAO5CMgPAbkHHNNboXzPy0RJeVKoguEUWWe-vAd9JwTrl0PiJZu3fNvZwGRQRdvlReWVk5FgBDYoGFbNi3rEuVvGdaxbVUm3YhkmhSSxerHbqeMsDPOkYhzeychiQClFUMyNW629cu2hjE0YNbspiYvpuY2uFVuIMQ76pROxhTrbibP0gPdIKD2wqAFSlIlqA1I8FnsrcSRuCI3bUyXexdE4hRuRGPXsggdUMPYwqk4axyN5RP9NS57kpM6Yu4bnNkzVeF1MqNxJZQcHlsusUV5JoxVMFxzbgxq2yY4lJB5Jxu5pAHBoZEh8JYB5xvd7h1ZugglqSS6KROa-YAivuvDXrBXBdwuTaBDLtgUEc4_wqbaaC8VIYK4Ui_E9xDcmLDM7WbxXBuGmeA1tiUfoREY7S4f_NF8YNeHAvPFSlRFgNNgkkySoMUp1sW6euDP5SMm6M52SOZirx_bdp2F_plDfGBNcyE6yez6AUkYXAVeZG9NHHM_VZbIAkRV6GU3MoFZFRExgk853QQqGzPXexp0AKqI8M0HrNvdhcdiDo3qOiAtDOOkZUz9PEPT0tsq1NUsxMSc_rDHXlqsaKiJxiQSKyS7mZscKQJINpqLzXAIcuPd0OsaNiWQOUd2ZzItjh6rbJCoqFTRiEV9bazMBGyWrbBnNMvJstkVQJLWxSG8DXO1ICewXlIicFM_DyAA8hSDFG7_3AJ6mS6n5J2L8I2hxtnWEkVeeK0RZQdEdCZkAeVihu3CnQjnxg80ZaLIjwzYUeBcgUNxPaGa1auU3FofcEcCcwqzmRtAfeS86N3I-hSplRnsrcifUUUEKUVHOV7M1icKKoNXFUw-rZKiCXJq6a77yoREsM2OtBzsmAoFIkmp0AYBZegUGNZ5ber57SixXZiQWx26KE2VXXPbDFLkZMN9oy_WWCHDRA5C0gMWtFw4OAmXrlys_CyTKl9YGA3JJVBawRfZdnDpPIIv-sZxAAiDVYPgwt3cdKtKhEdpSW1Q_oo4bC5Pesc7paU6SmfEoyRX8of3JmZHPa0b2miTQSz57c-1tNkHpev72EplKkxnCItRYRpsooO_QB1gxdEQtOVVNnMOZJN6Ih30tW8HMHAwYh24VAbgiTElDBUZKpVJbLDhDIb9chG5TIqOCxeEAcTFBy3HQEkjcBAdyGOm9x3yTnK3IQpOAaaKelpEwICcuoNbxxcFdanSW2EkfzpL7LaA8kDr9t0QF3T1hSBnwkSnJzItVE", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: vl, Name: "vless"})

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

	// vless
	// vl, err := vless.New("localhost", 39573, "433d0d23-d443-5212-a0df-cdbdd2d53015", "none", "")
	vl, err := vless.New("localhost", 39573, "433d0d23-d443-5212-a0df-cdbdd2d53015", "kjiwehh2tmsGrdtPjTPG03qPAEJqLeAJR4ClA_G4J0QAz7NyY2o0paRJo4B14TNavpIlbeUPjFRHIkhebjBOKxU8v3shXPUrHdha9ATNg2ILi1u16uS27JVBW6umvLh5ulqxq4JfFiR8qfY45WPEPIIAb4lpORkcqWJQmiMDfoxA_ccx0lnI5dWysiZHL6QvETiIKuhwwgYxCeyTxuFCRlBRwXnCPdsEg1THxaa1H6iB8hiXDMK0p9wz__FenZZsVAFTwBJ538pVFWeluLePCncJD0JEY8gKF2BF8bAzODwFwtM9zrlZciaXc2K9QHRS_GCc0EcWbjoAO5CMgPAbkHHNNboXzPy0RJeVKoguEUWWe-vAd9JwTrl0PiJZu3fNvZwGRQRdvlReWVk5FgBDYoGFbNi3rEuVvGdaxbVUm3YhkmhSSxerHbqeMsDPOkYhzeychiQClFUMyNW629cu2hjE0YNbspiYvpuY2uFVuIMQ76pROxhTrbibP0gPdIKD2wqAFSlIlqA1I8FnsrcSRuCI3bUyXexdE4hRuRGPXsggdUMPYwqk4axyN5RP9NS57kpM6Yu4bnNkzVeF1MqNxJZQcHlsusUV5JoxVMFxzbgxq2yY4lJB5Jxu5pAHBoZEh8JYB5xvd7h1ZugglqSS6KROa-YAivuvDXrBXBdwuTaBDLtgUEc4_wqbaaC8VIYK4Ui_E9xDcmLDM7WbxXBuGmeA1tiUfoREY7S4f_NF8YNeHAvPFSlRFgNNgkkySoMUp1sW6euDP5SMm6M52SOZirx_bdp2F_plDfGBNcyE6yez6AUkYXAVeZG9NHHM_VZbIAkRV6GU3MoFZFRExgk853QQqGzPXexp0AKqI8M0HrNvdhcdiDo3qOiAtDOOkZUz9PEPT0tsq1NUsxMSc_rDHXlqsaKiJxiQSKyS7mZscKQJINpqLzXAIcuPd0OsaNiWQOUd2ZzItjh6rbJCoqFTRiEV9bazMBGyWrbBnNMvJstkVQJLWxSG8DXO1ICewXlIicFM_DyAA8hSDFG7_3AJ6mS6n5J2L8I2hxtnWEkVeeK0RZQdEdCZkAeVihu3CnQjnxg80ZaLIjwzYUeBcgUNxPaGa1auU3FofcEcCcwqzmRtAfeS86N3I-hSplRnsrcifUUUEKUVHOV7M1icKKoNXFUw-rZKiCXJq6a77yoREsM2OtBzsmAoFIkmp0AYBZegUGNZ5ber57SixXZiQWx26KE2VXXPbDFLkZMN9oy_WWCHDRA5C0gMWtFw4OAmXrlys_CyTKl9YGA3JJVBawRfZdnDpPIIv-sZxAAiDVYPgwt3cdKtKhEdpSW1Q_oo4bC5Pesc7paU6SmfEoyRX8of3JmZHPa0b2miTQSz57c-1tNkHpev72EplKkxnCItRYRpsooO_QB1gxdEQtOVVNnMOZJN6Ih30tW8HMHAwYh24VAbgiTElDBUZKpVJbLDhDIb9chG5TIqOCxeEAcTFBy3HQEkjcBAdyGOm9x3yTnK3IQpOAaaKelpEwICcuoNbxxcFdanSW2EkfzpL7LaA8kDr9t0QF3T1hSBnwkSnJzItVE", "")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	pps = append(pps, ProxyPort{Proxy: vl, Name: "vless"})

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
