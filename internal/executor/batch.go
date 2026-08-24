package executor

import (
	"context"
	"errors"
	"example.com/field-sensor-orchestrator/internal/domain"
	"sync"
)

func (r *Runner) RunSites(ctx context.Context, sites []string, actor string) ([]domain.Cycle, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	out := make([]domain.Cycle, len(sites))
	errs := make(chan error, len(sites))
	var wg sync.WaitGroup
	for i, site := range sites {
		i, site := i, site
		wg.Add(1)
		go func() {
			defer wg.Done()
			cycle, err := r.Run(ctx, site, actor)
			if err != nil {
				errs <- err
				cancel()
				return
			}
			out[i] = cycle
		}()
	}
	wg.Wait()
	close(errs)
	var joined error
	for err := range errs {
		joined = errors.Join(joined, err)
	}
	return out, joined
}
