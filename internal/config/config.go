package config

import (
	"github.com/xtls/xray-core/app/dispatcher"
	"github.com/xtls/xray-core/app/log"
	"github.com/xtls/xray-core/app/policy"
	"github.com/xtls/xray-core/app/proxyman"
	"github.com/xtls/xray-core/app/stats"
	"github.com/xtls/xray-core/common/serial"
	"github.com/xtls/xray-core/core"
)

type Configurer interface {
	App() []*serial.TypedMessage
	Inbound() []*core.InboundHandlerConfig
	Outbound() []*core.OutboundHandlerConfig
}

func ToCoreConfig(cfg Configurer) core.Config {
	return core.Config{
		App:      cfg.App(),
		Inbound:  cfg.Inbound(),
		Outbound: cfg.Outbound(),
	}
}

func CommonApp() []*serial.TypedMessage {
	return []*serial.TypedMessage{
		// essential
		serial.ToTypedMessage(&dispatcher.Config{}),
		serial.ToTypedMessage(&proxyman.InboundConfig{}),
		serial.ToTypedMessage(&proxyman.OutboundConfig{}),

		// disable log
		serial.ToTypedMessage(&log.Config{
			AccessLogType: log.LogType_None,
			ErrorLogType:  log.LogType_None,
		}),

		// enable traffic stats
		serial.ToTypedMessage(&stats.Config{}),
		serial.ToTypedMessage(&policy.Config{
			System: &policy.SystemPolicy{
				Stats: &policy.SystemPolicy_Stats{
					InboundUplink:    true,
					InboundDownlink:  true,
					OutboundUplink:   true,
					OutboundDownlink: true,
				},
			},
		}),
	}
}
