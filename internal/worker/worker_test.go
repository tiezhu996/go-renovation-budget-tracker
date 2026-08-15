package worker

import (
	"context"
	"testing"

	"renobudget/internal/config"
	"renobudget/internal/model"
	"renobudget/internal/service"
	"renobudget/internal/store"
)

func newPool() (*service.Service, *Pool) {
	s := store.New()
	svc := service.New(s, config.Load())
	return svc, New(svc, 4)
}

func TestCheckSummary(t *testing.T) {
	svc, p := newPool()
	for i := 0; i < 10; i++ {
		_ = svc.AddItem(&model.BudgetItem{ID: string(rune('a' + i)), Name: "x", Budget: 1000})
		_ = svc.RecordExpense(&model.ExpenseRecord{ID: "e" + string(rune('a'+i)), ItemID: string(rune('a' + i)), Amount: 950})
	}
	sum := p.Check(context.Background())
	if sum.Checked != 10 || sum.Alerted != 10 || sum.Failed != 0 {
		t.Fatalf("summary=%+v", sum)
	}
}

func TestCheckNoAlert(t *testing.T) {
	svc, p := newPool()
	_ = svc.AddItem(&model.BudgetItem{ID: "i1", Name: "x", Budget: 1000})
	_ = svc.RecordExpense(&model.ExpenseRecord{ID: "e1", ItemID: "i1", Amount: 100})
	sum := p.Check(context.Background())
	if sum.Alerted != 0 {
		t.Fatalf("summary=%+v", sum)
	}
}

func TestCheckCancel(t *testing.T) {
	svc, p := newPool()
	_ = svc.AddItem(&model.BudgetItem{ID: "i1", Name: "x", Budget: 1000})
	_ = svc.RecordExpense(&model.ExpenseRecord{ID: "e1", ItemID: "i1", Amount: 950})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	sum := p.Check(ctx)
	if sum.Alerted != 0 {
		t.Fatalf("alerted=%d want 0", sum.Alerted)
	}
}
