# Field Sensor Orchestrator

Field Sensor Orchestrator 是面向野外生态观测团队的本地 Go CLI。现场技术员可登记与校准设备，观测协调员可安排并运行采集周期、确认告警，数据管理员可导出指定站点的观测摘要。

## 目录结构
- `cmd/sensorctl`：命令入口与演示工作流。
- `internal/domain`：传感器、站点、任务、租约、读数、校准、告警与审计模型。
- `internal/store`：线程安全的内存状态仓库及快照。
- `internal/scheduler`：采集任务选择、排程和租约创建。
- `internal/executor`：任务领取、上下文传播、并发采集和资源收尾。
- `internal/device`：设备会话、模拟读数与故障注入。
- `internal/calibration`：校准流程与设备状态更新。
- `internal/alerts`：告警创建、确认和站点隔离。
- `internal/report`：站点观测摘要与健康度聚合。
- `internal/audit`：结构化审计事件生成。
- `internal/app`：跨模块工作流编排。

## 运行
```bash
go run ./cmd/sensorctl demo
go run ./cmd/sensorctl register-sensor
go run ./cmd/sensorctl schedule-cycle
go run ./cmd/sensorctl run-cycle
go run ./cmd/sensorctl calibrate-sensor
go run ./cmd/sensorctl ack-alert
go run ./cmd/sensorctl export-summary
```

## 构建与测试
```bash
go build ./...
go test ./...
```

CLI 使用内存状态并在每次进程启动时装载演示站点与设备，无需数据库或外部网络。可选环境变量 `SENSOR_SITE` 用于覆盖演示站点标识，默认值为 `ridge-north`。
