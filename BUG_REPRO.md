# BUG_REPRO

## Bug 是什么
BuildAlertBatches 返回共享底层数组的子切片、AlertBatches 手写共享切片、ItemIDs 直接返回内部 order、worker 把每个 batch 截掉最后一个，导致列表被外部改动污染、超支检查漏项。

## 如何触发
`go test ./...`

## 错误信息
- TestBuildAlertBatchesFresh / TestItemIDsFresh / TestCheckSummary 失败。
