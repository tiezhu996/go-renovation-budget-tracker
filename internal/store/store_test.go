package store

import (
	"testing"

	"renobudget/internal/model"
)

func TestPutGetItem(t *testing.T) {
	s := New()
	if err := s.PutItem(&model.BudgetItem{ID: "i1", Name: "水电", Budget: 1000}); err != nil {
		t.Fatal(err)
	}
	if err := s.PutItem(&model.BudgetItem{ID: "i1", Name: "x", Budget: 1}); err != ErrItemExists {
		t.Fatalf("dup err=%v", err)
	}
	if _, err := s.GetItem("i1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.GetItem("nope"); err != ErrItemNotFound {
		t.Fatalf("missing err=%v", err)
	}
}

func TestItemIDsFresh(t *testing.T) {
	s := New()
	_ = s.PutItem(&model.BudgetItem{ID: "i2", Name: "x", Budget: 1})
	_ = s.PutItem(&model.BudgetItem{ID: "i1", Name: "x", Budget: 1})
	ids := s.ItemIDs()
	ids[0] = "zz"
	if s.ItemIDs()[0] != "i2" {
		t.Fatal("ItemIDs aliased")
	}
}

func TestListItemsFresh(t *testing.T) {
	s := New()
	_ = s.PutItem(&model.BudgetItem{ID: "i1", Name: "水电", Budget: 1000})
	items := s.ListItems()
	items[0] = &model.BudgetItem{ID: "zz"}
	if it, _ := s.GetItem("i1"); it.Name != "水电" {
		t.Fatal("ListItems mutated store")
	}
}

func TestRecordExpense(t *testing.T) {
	s := New()
	_ = s.PutItem(&model.BudgetItem{ID: "i1", Name: "水电", Budget: 1000})
	if err := s.RecordExpense(&model.ExpenseRecord{ID: "e1", ItemID: "i1", Amount: 300}); err != nil {
		t.Fatal(err)
	}
	it, _ := s.GetItem("i1")
	if it.Spent != 300 {
		t.Fatalf("spent=%v", it.Spent)
	}
	if err := s.RecordExpense(&model.ExpenseRecord{ID: "e2", ItemID: "nope", Amount: 1}); err != ErrItemNotFound {
		t.Fatalf("missing err=%v", err)
	}
}

func TestAddListAlert(t *testing.T) {
	s := New()
	id, err := s.AddAlert(&model.Alert{ItemID: "i1", Ratio: 0.95})
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Fatal("empty alert id")
	}
	if _, err := s.AddAlert(&model.Alert{ItemID: "i1", Ratio: 0.96}); err != ErrAlertExists {
		t.Fatalf("dup alert err=%v", err)
	}
	if len(s.ListAlerts()) != 1 {
		t.Fatal("alert list len")
	}
}

func TestMarkAlertSent(t *testing.T) {
	s := New()
	id, _ := s.AddAlert(&model.Alert{ItemID: "i1", Ratio: 0.95})
	if err := s.MarkAlertSent(id); err != nil {
		t.Fatal(err)
	}
	if !s.ListAlerts()[0].Sent {
		t.Fatal("alert not marked sent")
	}
	if err := s.MarkAlertSent("nope"); err != ErrAlertNotFound {
		t.Fatalf("missing err=%v", err)
	}
}
