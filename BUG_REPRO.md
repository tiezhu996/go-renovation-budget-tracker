# BUG_REPRO

## Bug 是什么
Store.New 没有初始化 items map，首次 RecordExpense 往 nil map 读/写时崩溃。

## 如何触发
`go test ./...`

## 错误信息
- panic: assignment to entry in nil map
