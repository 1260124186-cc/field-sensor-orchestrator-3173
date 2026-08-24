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
docker run --rm field-sensor-orchestrator-repro bash -c "go test -count=1 -run '^TestFailedCycleReleasesLeaseForNextSchedule$' ./internal/app"
```

## 实际错误输出

```text
--- FAIL: TestFailedCycleReleasesLeaseForNextSchedule (0.00s)
    bug005_test.go:20: failed cycle retained its lease
FAIL
FAIL	example.com/field-sensor-orchestrator/internal/app	3.053s
FAIL
```

## 期望行为

失败周期应释放租约和设备会话，使同站点下一次排程能够正常开始并完成。
