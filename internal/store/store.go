package store

import (
	"errors"
	"strconv"
	"sync"

	"renobudget/internal/model"
)

var (
	ErrItemNotFound  = errors.New("item not found")
	ErrItemExists    = errors.New("item already exists")
	ErrAlertNotFound = errors.New("alert not found")
	ErrAlertExists   = errors.New("alert already exists")
)

type Store struct {
	mu          sync.RWMutex
	items       map[string]*model.BudgetItem
	alerts      map[string]*model.Alert
	itemOrder   []string
	alertOrder  []string
	nextAlertID int
}

func New() *Store {
	return &Store{
		items:       make(map[string]*model.BudgetItem),
		alerts:      make(map[string]*model.Alert),
		itemOrder:   []string{},
		alertOrder:  []string{},
		nextAlertID: 1,
	}
}

func (s *Store) PutItem(it *model.BudgetItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[it.ID]; ok {
		return ErrItemExists
	}
	s.items[it.ID] = it
	s.itemOrder = append(s.itemOrder, it.ID)
	return nil
}

func (s *Store) GetItem(id string) (*model.BudgetItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	if !ok {
		return nil, ErrItemNotFound
	}
	return it, nil
}

func (s *Store) ListItems() []*model.BudgetItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.BudgetItem, 0, len(s.itemOrder))
	for _, id := range s.itemOrder {
		out = append(out, s.items[id])
	}
	return out
}

func (s *Store) ItemIDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.itemOrder
}

func (s *Store) RecordExpense(e *model.ExpenseRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	it, ok := s.items[e.ItemID]
	if !ok {
		return ErrItemNotFound
	}
	it.Spent += e.Amount
	return nil
}

func (s *Store) AddAlert(a *model.Alert) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range s.alertOrder {
		if s.alerts[id].ItemID == a.ItemID {
			return "", ErrAlertExists
		}
	}
	a.ID = "alert-" + strconv.Itoa(s.nextAlertID)
	s.nextAlertID++
	s.alerts[a.ID] = a
	s.alertOrder = append(s.alertOrder, a.ID)
	return a.ID, nil
}

func (s *Store) ListAlerts() []*model.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Alert, 0, len(s.alertOrder))
	for _, id := range s.alertOrder {
		out = append(out, s.alerts[id])
	}
	return out
}

func (s *Store) MarkAlertSent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.alerts[id]
	if !ok {
		return ErrAlertNotFound
	}
	a.Sent = true
	return nil
}
