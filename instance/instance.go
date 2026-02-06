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
	Inst *core.Instance
}

func New(ctx context.Context, cfg *core.Config) (*Instance, error) {
	inst, err := core.NewWithContext(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &Instance{
		ctx:  ctx,
		Inst: inst,
	}, nil
}

func (i *Instance) Start() error    { return i.Inst.Start() }
func (i *Instance) IsStarted() bool { return i.Inst.IsRunning() }
func (i *Instance) Stop() error     { return i.Inst.Close() }
