package common

import (
	"github.com/xtls/xray-core/app/dispatcher"
	"github.com/xtls/xray-core/app/log"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
)

func NewConfig() *core.Config {
	config := core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),

			// disable log
			serial.ToTypedMessage(&log.Config{
				AccessLogType: log.LogType_None,
				ErrorLogType:  log.LogType_None,
			}),
		},
	}

	return &config
}
