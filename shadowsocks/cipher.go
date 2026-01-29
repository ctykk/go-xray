package shadowsocks

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xtls/xray-core/proxy/shadowsocks"
)

type Cipher = shadowsocks.CipherType

const (
	CipherUnknown = shadowsocks.CipherType_UNKNOWN

	CipherNone = shadowsocks.CipherType_NONE

	CipherAES128GCM = shadowsocks.CipherType_AES_128_GCM
	CipherAES256GCM = shadowsocks.CipherType_AES_256_GCM

	CipherChacha20Poly1305  = shadowsocks.CipherType_CHACHA20_POLY1305
	CipherXchacha20Poly1305 = shadowsocks.CipherType_XCHACHA20_POLY1305
)

var strCipherRe = regexp.MustCompile(`[-_]`)

// ParseCipher parse Cipher from string
func ParseCipher(s string) (Cipher, error) {
	switch strings.ToLower(strCipherRe.ReplaceAllString(s, "")) {
	case "aes128gcm":
		return CipherAES128GCM, nil
	case "aes256gcm":
		return CipherAES256GCM, nil
	case "chacha20poly1305":
		return CipherChacha20Poly1305, nil
	case "xchacha20poly1305":
		return CipherXchacha20Poly1305, nil
	case "none":
		return CipherNone, nil
	}
	return CipherUnknown, fmt.Errorf("unknown cipher: %s", s)
}
