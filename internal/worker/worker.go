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

	var sum model.Summary

	for i := 0; i < p.workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for batch := range ch {
				var local model.Summary
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
				sum = model.MergeSummary(sum, local)
			}
		}()
	}

	wg.Wait()
	return sum
}
