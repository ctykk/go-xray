package vmess

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xtls/xray-core/common/protocol"
)

type Cipher = protocol.SecurityType

const (
	CipherUnknown = protocol.SecurityType_UNKNOWN

	CipherNone = protocol.SecurityType_NONE
	CipherZero = protocol.SecurityType_ZERO

	CipherAuto = protocol.SecurityType_AUTO

	CipherAES128GCM        = protocol.SecurityType_AES128_GCM
	CipherChacha20Poly1305 = protocol.SecurityType_CHACHA20_POLY1305
)

var strCipherRe = regexp.MustCompile(`[-_]`)

// StrCipher parse Cipher from string
func StrCipher(s string) (Cipher, error) {
	s = strCipherRe.ReplaceAllString(s, "")
	s = strings.ToLower(s)
	switch s {
	case "auto":
		return CipherAuto, nil
	case "aes128gcm":
		return CipherAES128GCM, nil
	case "chacha20poly1305":
		return CipherChacha20Poly1305, nil
	case "none":
		return CipherNone, nil
	case "zero":
		return CipherZero, nil
	}
	return CipherUnknown, fmt.Errorf("unknown cipher: %s", s)
}
