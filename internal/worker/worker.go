package worker

import (
	"context"
	"sync"

	"renobudget/internal/model"
	"renobudget/internal/service"
)

type Pool struct {
	svc     *service.Service
	workers int
}

func New(svc *service.Service, workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}
	return &Pool{svc: svc, workers: workers}
}

func (p *Pool) Check(ctx context.Context) model.Summary {
	if err := ctx.Err(); err != nil {
		return model.Summary{}
	}

	batches := p.svc.AlertBatches()

	var wg sync.WaitGroup
	ch := make(chan []*model.BudgetItem, len(batches))

	go func() {
		defer close(ch)
		for _, b := range batches {
			select {
			case <-ctx.Done():
				return
			case ch <- b:
			}
		}
	}()

	results := make(chan model.Summary, p.workers)

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var local model.Summary
			for batch := range ch {
				if ctx.Err() != nil {
					break
				}
				for _, it := range batch {
					triggered, err := p.svc.EvaluateItem(it)
					local.Checked++
					if err != nil {
						local.Failed++
						continue
					}
					if triggered {
						local.Alerted++
					}
				}
			}
			results <- local
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var sum model.Summary
	for local := range results {
		sum = model.MergeSummary(sum, local)
	}
	return sum
}
