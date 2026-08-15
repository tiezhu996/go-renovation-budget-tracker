package service

import (
	"errors"
	"fmt"

	"renobudget/internal/config"
	"renobudget/internal/model"
	"renobudget/internal/store"
)

type Service struct {
	store     *store.Store
	batchSize int
	threshold float64
}

func New(s *store.Store, cfg *config.Config) *Service {
	b := cfg.BatchSize
	if b <= 0 {
		b = 1
	}
	t := cfg.Threshold
	if t <= 0 {
		t = 0.9
	}
	return &Service{store: s, batchSize: b, threshold: t}
}

func (svc *Service) AddItem(it *model.BudgetItem) error {
	if !model.ValidItem(it) {
		return errors.New("invalid item")
	}
	if err := svc.store.PutItem(it); err != nil {
		return fmt.Errorf("add item %s: %w", it.ID, err)
	}
	return nil
}

func (svc *Service) RecordExpense(e *model.ExpenseRecord) error {
	if !model.ValidExpense(e) {
		return errors.New("invalid expense")
	}
	if err := svc.store.RecordExpense(e); err != nil {
		return fmt.Errorf("record expense %s: %v", e.ID, err)
	}
	return nil
}

func (svc *Service) ListItems() []*model.BudgetItem {
	return svc.store.ListItems()
}

func (svc *Service) ListAlerts() []*model.Alert {
	return svc.store.ListAlerts()
}

func (svc *Service) EvaluateItem(it *model.BudgetItem) (bool, error) {
	ratio := model.Ratio(it.Spent, it.Budget)
	if !model.AlertTriggered(ratio, svc.threshold) {
		return false, nil
	}
	_, err := svc.store.AddAlert(&model.Alert{ItemID: it.ID, Ratio: ratio})
	if err != nil {
		return false, fmt.Errorf("add alert %s: %v", it.ID, err)
	}
	return true, nil
}

func (svc *Service) AlertBatches() [][]*model.BudgetItem {
	items := svc.store.ListItems()
	model.SortItems(items)
	return model.BuildAlertBatches(items, svc.batchSize)
}
