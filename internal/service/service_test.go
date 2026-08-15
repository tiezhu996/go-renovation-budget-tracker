package service

import (
	"errors"
	"testing"

	"renobudget/internal/config"
	"renobudget/internal/model"
	"renobudget/internal/store"
)

func newSvc() (*store.Store, *Service) {
	s := store.New()
	return s, New(s, config.Load())
}

func TestAddItemRejectsInvalid(t *testing.T) {
	_, svc := newSvc()
	if err := svc.AddItem(&model.BudgetItem{ID: "i1", Name: "水电", Budget: 1000}); err != nil {
		t.Fatal(err)
	}
	if err := svc.AddItem(&model.BudgetItem{}); err == nil {
		t.Fatal("invalid item accepted")
	}
	if err := svc.AddItem(&model.BudgetItem{ID: "i1", Name: "x", Budget: 1}); !errors.Is(err, store.ErrItemExists) {
		t.Fatalf("dup err=%v", err)
	}
}

func TestRecordExpenseRejectsInvalid(t *testing.T) {
	_, svc := newSvc()
	_ = svc.AddItem(&model.BudgetItem{ID: "i1", Name: "水电", Budget: 1000})
	if err := svc.RecordExpense(&model.ExpenseRecord{ID: "e1", ItemID: "i1", Amount: 100}); err != nil {
		t.Fatal(err)
	}
	if err := svc.RecordExpense(&model.ExpenseRecord{ID: "e2", ItemID: "i1"}); err == nil {
		t.Fatal("invalid expense accepted")
	}
	if err := svc.RecordExpense(&model.ExpenseRecord{ID: "e3", ItemID: "nope", Amount: 1}); !errors.Is(err, store.ErrItemNotFound) {
		t.Fatalf("missing err=%v", err)
	}
}

func TestEvaluateItem(t *testing.T) {
	_, svc := newSvc()
	_ = svc.AddItem(&model.BudgetItem{ID: "i1", Name: "水电", Budget: 1000})
	_ = svc.RecordExpense(&model.ExpenseRecord{ID: "e1", ItemID: "i1", Amount: 950})
	it, _ := svc.store.GetItem("i1")
	triggered, err := svc.EvaluateItem(it)
	if err != nil {
		t.Fatal(err)
	}
	if !triggered {
		t.Fatal("should trigger at 95%")
	}
	if len(svc.ListAlerts()) != 1 {
		t.Fatal("alert not created")
	}
	// second evaluation should not duplicate alert
	_, _ = svc.EvaluateItem(it)
	if len(svc.ListAlerts()) != 1 {
		t.Fatal("duplicate alert created")
	}
}

func TestListItemsAndBatches(t *testing.T) {
	_, svc := newSvc()
	_ = svc.AddItem(&model.BudgetItem{ID: "i2", Name: "木工", Budget: 1000})
	_ = svc.AddItem(&model.BudgetItem{ID: "i1", Name: "水电", Budget: 1000})
	bs := svc.AlertBatches()
	if len(bs) == 0 || bs[0][0].ID != "i1" {
		t.Fatalf("batches=%v", bs)
	}
}
