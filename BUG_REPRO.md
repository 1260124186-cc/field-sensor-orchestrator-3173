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
docker run --rm field-sensor-orchestrator-repro bash -c "go test -count=1 -run '^TestCancelledRunCycleStopsCollection$' ./internal/app"
```

## 实际错误输出

```text
--- FAIL: TestCancelledRunCycleStopsCollection (0.00s)
    bug002_test.go:18: run error=<nil> want context canceled
FAIL
FAIL	example.com/field-sensor-orchestrator/internal/app	2.550s
FAIL
```

## 期望行为

取消应沿采集调用链生效，返回 context canceled；周期应失败、不写入读数并释放租约。
