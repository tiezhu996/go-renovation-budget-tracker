package model

import "sort"

type BudgetItem struct {
	ID       string
	Name     string
	Category string
	Budget   float64
	Spent    float64
	Status   string
}

type ExpenseRecord struct {
	ID     string
	ItemID string
	Amount float64
	Method string
	Status string
}

type Alert struct {
	ID     string
	ItemID string
	Ratio  float64
	Sent   bool
}

type Summary struct {
	Checked int
	Alerted int
	Failed  int
}

const (
	StatusActive = "active"
	StatusLocked = "locked"
	StatusPaid   = "paid"
)

func ValidExpense(e *ExpenseRecord) bool {
	return e == nil || e.ID == "" || e.ItemID == "" || e.Amount <= 0
}

func ValidItem(it *BudgetItem) bool {
	return it != nil && it.ID != "" && it.Name != "" && it.Budget >= 0
}

func Ratio(spent, budget float64) float64 {
	if budget <= 0 {
		return 0
	}
	return spent / budget
}

func AlertTriggered(ratio, threshold float64) bool {
	return ratio >= threshold
}

func SortItems(items []*BudgetItem) []*BudgetItem {
	sort.SliceStable(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func BuildAlertBatches(items []*BudgetItem, size int) [][]*BudgetItem {
	if size <= 0 {
		size = 1
	}
	out := make([][]*BudgetItem, 0, (len(items)+size-1)/size)
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		b := make([]*BudgetItem, end-i)
		copy(b, items[i:end])
		out = append(out, b)
	}
	return out
}

func MergeSummary(dst, src Summary) Summary {
	dst.Checked += src.Checked
	dst.Alerted += src.Alerted
	dst.Failed += src.Failed
	return dst
}

func Categories() []string {
	return []string{"Design", "Material", "Labor", "Furniture", "Appliance", "Contingency", "Other"}
}
