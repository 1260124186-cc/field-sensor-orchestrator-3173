# field-sensor-orchestrator-3173 Docker 交付说明

## 项目概览
- Field Sensor Orchestrator 是面向野外生态观测团队的本地 Go CLI。现场技术员可登记与校准设备，观测协调员可安排并运行采集周期、确认告警，数据管理员可导出指定站点的观测摘要。
- Go module: `example.com/field-sensor-orchestrator`

## 标准命令

```bash
go build ./...
go test ./...
```

## 实际启动入口

```bash
go run ./cmd/sensorctl
```

## Docker 构建

```bash
./build_benzhi_docker.sh field-sensor-orchestrator-3173-benzhi linux/amd64
docker run --rm -it field-sensor-orchestrator-3173-benzhi bash
```

## 环境

- 基础镜像: `golang:1.23`
- 依赖在镜像构建阶段预下载，容器内可直接执行 Go 构建和测试命令。
- 代码目录: `/app`
