package model

import "testing"

func TestValidExpense(t *testing.T) {
	if !ValidExpense(&ExpenseRecord{ID: "e1", ItemID: "i1", Amount: 10}) {
		t.Fatal("valid expense rejected")
	}
	for _, e := range []*ExpenseRecord{nil, {}, {ID: "e1", ItemID: "i1"}, {ID: "e1", ItemID: "", Amount: 10}} {
		if ValidExpense(e) {
			t.Fatalf("invalid expense accepted %+v", e)
		}
	}
}

func TestValidItem(t *testing.T) {
	if !ValidItem(&BudgetItem{ID: "i1", Name: "水电", Budget: 1000}) {
		t.Fatal("valid item rejected")
	}
	if ValidItem(nil) || ValidItem(&BudgetItem{ID: "i1"}) || ValidItem(&BudgetItem{ID: "", Name: "x", Budget: 0}) {
		t.Fatal("invalid item accepted")
	}
}

func TestRatioAndTrigger(t *testing.T) {
	if Ratio(900, 1000) != 0.9 {
		t.Fatalf("ratio=%v", Ratio(900, 1000))
	}
	if Ratio(100, 0) != 0 {
		t.Fatal("zero budget ratio should be 0")
	}
	if !AlertTriggered(0.9, 0.9) || AlertTriggered(0.89, 0.9) {
		t.Fatal("threshold boundary wrong")
	}
}

func TestBuildAlertBatchesFresh(t *testing.T) {
	items := []*BudgetItem{{ID: "i1"}, {ID: "i2"}, {ID: "i3"}}
	bs := BuildAlertBatches(items, 2)
	if len(bs) != 2 {
		t.Fatalf("batches=%d", len(bs))
	}
	bs[0][0] = &BudgetItem{ID: "z"}
	if items[0].ID != "i1" {
		t.Fatal("mutating batch corrupted input")
	}
}

func TestMergeSummary(t *testing.T) {
	got := MergeSummary(Summary{Checked: 1, Failed: 1}, Summary{Checked: 2, Alerted: 3, Failed: 4})
	if got.Checked != 3 || got.Alerted != 3 || got.Failed != 5 {
		t.Fatalf("merge=%+v", got)
	}
}

func TestCategories(t *testing.T) {
	if len(Categories()) != 7 {
		t.Fatalf("categories=%v", Categories())
	}
}
