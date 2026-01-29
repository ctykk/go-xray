package common

import (
	"context"
	"fmt"

	"github.com/xtls/xray-core/core"
)

func NewHTTPProxy(ctx context.Context, config *core.Config) error {
	instance, err := core.NewWithContext(ctx, config)
	if err != nil {
		return fmt.Errorf("init instance: %w", err)
	}

	err = instance.Start()
	if err != nil {
		return fmt.Errorf("start instance: %w", err)
	}
	go func() {
		<-ctx.Done()
		_ = instance.Close()
	}()
	return nil
}
