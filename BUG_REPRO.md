# BUG_REPRO

## Bug 是什么
ValidExpense 校验写反、RecordExpense 丢掉了入参校验、MarkAlertSent 去掉锁、worker 把 wg.Add 放进 goroutine 且汇总无锁，导致无效支出被接受、告警重复、统计错误并存在 data race。

## 如何触发
`go test -race ./...`

## 错误信息
- TestValidExpense / TestRecordExpenseRejectsInvalid / TestCheckSummary 失败。
- `go test -race` 报 data race。
