# 修复前故障复现（Docker）

## 项目与标准命令

项目为 Go 1.23 CLI。标准检查命令：

```bash
go build ./...
go test ./...
```

## 环境构建与编译

当前验证平台为 `linux/arm64`，容器内 Go 版本为 `go1.23.12`。执行：

```bash
docker build -f benzhi.Dockerfile -t field-sensor-orchestrator-repro .
docker run --rm field-sensor-orchestrator-repro bash -c 'go build ./...'
```

镜像构建与容器内 `go build ./...` 均成功。

## 故障触发步骤

在包含随附测试的镜像中执行：

```bash
docker run --rm -e GORACE=halt_on_error=1 field-sensor-orchestrator-repro bash -c "go test -race -count=1 -run '^TestConcurrentCycleScheduling$' ./internal/app"
```

## 实际错误输出

```text
==================
WARNING: DATA RACE
Write at 0x0000003b3e20 by goroutine 12:
  example.com/field-sensor-orchestrator/internal/app.(*App).Schedule()
      /app/internal/app/app.go:69 +0x58
  example.com/field-sensor-orchestrator/internal/app.TestConcurrentCycleScheduling.func1()
      /app/internal/app/bug001_test.go:19 +0x110
  example.com/field-sensor-orchestrator/internal/app.TestConcurrentCycleScheduling.gowrap1()
      /app/internal/app/bug001_test.go:22 +0x44

Previous write at 0x0000003b3e20 by goroutine 10:
  example.com/field-sensor-orchestrator/internal/app.(*App).Schedule()
      /app/internal/app/app.go:69 +0x58
  example.com/field-sensor-orchestrator/internal/app.TestConcurrentCycleScheduling.func1()
      /app/internal/app/bug001_test.go:19 +0x110
  example.com/field-sensor-orchestrator/internal/app.TestConcurrentCycleScheduling.gowrap1()
      /app/internal/app/bug001_test.go:22 +0x44

Goroutine 12 (running) created at:
  example.com/field-sensor-orchestrator/internal/app.TestConcurrentCycleScheduling()
      /app/internal/app/bug001_test.go:17 +0xa0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1690 +0x184
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1743 +0x40

Goroutine 10 (running) created at:
  example.com/field-sensor-orchestrator/internal/app.TestConcurrentCycleScheduling()
      /app/internal/app/bug001_test.go:17 +0xa0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1690 +0x184
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1743 +0x40
==================
FAIL	example.com/field-sensor-orchestrator/internal/app	0.005s
FAIL
```

## 期望行为

并发排程不应触发数据竞争；每个请求都应得到唯一周期编号，全部排程都应成功落库。
