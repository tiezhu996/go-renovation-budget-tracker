# renovation-budget

一个用 Go 写的内存装修预算与支出服务，演示预算项、支出记账、超支告警与并发检查 worker。

## 功能
- 预算项管理、支出记账
- 超支比例计算与告警
- 并发告警检查 worker 池，支持 context 取消

## 目录结构
```
cmd/renobudget/      程序入口
internal/config/     环境配置
internal/model/      模型与纯工具函数
internal/store/      内存存储（预算项 + 告警 + 锁）
internal/service/    业务逻辑
internal/worker/     告警检查 worker 池
```

## 运行与测试
```bash
go build ./...
go test ./...
```
