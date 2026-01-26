package shadowsocks

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// `ss://<base64-encoded-auth>@<host>:<port>#<name>`
var ssURLRe = regexp.MustCompile(`^ss://([A-Za-z0-9]+)@([a-z0-9.]+?\.[a-z]+?):(\d{1,5})#(.*?)$`)

// FromBase64 parses shadowsocks nodes from a base64 encoded list
func FromBase64(b64 string) ([]Node, error) {
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("b64 decode: %w", err)
	}

	// key as "host:port"
	nodes := make(map[string]Node)

	for _, line := range strings.Split(string(decoded), "\n") {
		line = strings.TrimRight(line, "\r")

		match := ssURLRe.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		// port
		port, err := strconv.Atoi(match[3])
		if err != nil {
			continue
		}

		// cipher-type, password
		encodedAuth, err := base64.StdEncoding.DecodeString(match[1])
		if err != nil {
			return nil, err
		}
		authSplit := strings.SplitN(string(encodedAuth), ":", 2)
		if len(authSplit) != 2 {
			continue
		}

		// name
		name, err := url.QueryUnescape(match[4])
		if err != nil {
			continue
		}

		node, err := New(match[2], uint16(port), authSplit[0], authSplit[1], name)
		if err != nil {
			continue
		}
		nodes[fmt.Sprintf("%s:%d", strings.ToLower(node.host), node.port)] = *node
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found")
	}

	result := make([]Node, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, node)
	}
	return result, nil
}
