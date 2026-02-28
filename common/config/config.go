package config

import (
	"github.com/xtls/xray-core/app/dispatcher"
	"github.com/xtls/xray-core/app/log"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
	"google.golang.org/protobuf/proto"
)

type Config struct {
	CoreConfig *core.Config
}

func DefaultConfig() *Config {
	cfg := &Config{
		CoreConfig: &core.Config{
			App: []*serial.TypedMessage{
				serial.ToTypedMessage(&dispatcher.Config{}),
				serial.ToTypedMessage(&proxyman.InboundConfig{}),
				serial.ToTypedMessage(&proxyman.OutboundConfig{}),

				serial.ToTypedMessage(&log.Config{
					AccessLogType: log.LogType_None,
					ErrorLogType:  log.LogType_None,
				}),
			},
		},
	}
	return cfg
}

func (c *Config) Clone() *Config {
	newC := &Config{}
	newC.CoreConfig = proto.Clone(c.CoreConfig).(*core.Config)
	return newC
}
