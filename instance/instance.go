package instance

import (
	"context"

	"github.com/xtls/xray-core/core"
)

// TODO: 如何编排
//  - Proxy 的 DialContext, HTTPProxy
//  - Instance 的 Stats
//  的关系？

type Instance struct {
	ctx  context.Context
	inst *core.Instance
}

func New(ctx context.Context, cfg *core.Config) (*Instance, error) {
	inst, err := core.NewWithContext(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Instance{
		ctx:  ctx,
		inst: inst,
	}, nil
}

func (i *Instance) Start() error    { return i.inst.Start() }
func (i *Instance) IsStarted() bool { return i.inst.IsRunning() }
func (i *Instance) Stop() error     { return i.inst.Close() }
