package instance

import (
	"errors"
	"time"

	"github.com/xtls/xray-core/features/policy"
	"github.com/xtls/xray-core/features/stats"
)

// Stats
// uplink and downlink traffic (measure: bytes)
type Stats struct {
	Time time.Time

	OutboundUplink   int64
	OutboundDownlink int64
}

// NOTE: https://deepwiki.com/search/-xraycore-go-go-import-githubc_0921d0e3-36b5-4bf3-aaef-a31dd6df0664

const (
	outboundUplinkName   = "outbound>>>>>>traffic>>>uplink"
	outboundDownlinkName = "outbound>>>>>>traffic>>>downlink"
)

func (i *Instance) Stats() (Stats, error) {
	if !i.IsStarted() {
		return Stats{}, errors.New("not started")
	}

	// 获取统计管理器
	statsManager, ok := i.inst.GetFeature(stats.ManagerType()).(stats.Manager)
	if !ok || statsManager == nil {
		return Stats{}, errors.New("no stats manager")
	}

	// 获取策略管理器
	policyManager, ok := i.inst.GetFeature(policy.ManagerType()).(policy.Manager)
	if !ok || policyManager == nil {
		return Stats{}, errors.New("no policy manager")
	}

	// 检查是否启用对应统计策略
	if !policyManager.ForSystem().Stats.OutboundUplink {
		return Stats{}, errors.New("outbound uplink disabled")
	}
	if !policyManager.ForSystem().Stats.OutboundDownlink {
		return Stats{}, errors.New("outbound downlink disabled")
	}

	result := Stats{
		Time:             time.Now(),
		OutboundUplink:   statsManager.GetCounter(outboundUplinkName).Value(),
		OutboundDownlink: statsManager.GetCounter(outboundDownlinkName).Value(),
	}

	return result, nil
}
