# kapp

kapp 是一个使用 Go、MyGo 和 SQLite 构建的跨平台桌面及移动应用。

## 环境要求

- Go（版本以 `go.mod` 为准）
- [Task](https://taskfile.dev)
- 平台构建所需的 MyGo 工具链

## 开发

运行桌面应用：

```bash
task dev
```

运行 Android 应用：

```bash
task dev-android
```

构建当前平台：

```bash
task build
```

构建 Android ARM64 应用：

```bash
task build-android
```

## 测试

```bash
CGO_ENABLED=0 go test ./...
```

## 发布构建

构建所有支持平台的产物：

```bash
task release
```

该命令会：

1. 清空并重新创建 `build/` 与 `dist/`
2. 构建 macOS、Windows、Linux 和 Android ARM64 产物
3. 将安装包和可执行文件复制或打包到 `dist/`

发布产物包含：

| 平台    | 架构      | 文件格式                  |
| ------- | --------- | ------------------------- |
| macOS   | Universal | DMG、APP `.tar.gz`        |
| macOS   | ARM64     | DMG、APP `.tar.gz`        |
| macOS   | AMD64     | DMG、APP `.tar.gz`        |
| Windows | AMD64     | `.exe`                    |
| Windows | ARM64     | `.exe`                    |
| Linux   | AMD64     | `.deb`、`.rpm`、`.tar.gz` |
| Linux   | ARM64     | `.deb`、`.rpm`、`.tar.gz` |
| Android | ARM64     | `.apk`、`.aab`            |

版本号由 `mygo.json` 的 `version` 字段提供，产物文件名会包含该版本号和平台架构。

## 自动发布

GitHub Actions 在以下情况触发发布：

- 推送 `dev` 分支：创建预发布 Beta Release
- 推送 `v*` 标签：创建正式 Release

工作流会运行测试、执行 `task release`、生成 `dist/checksums.txt`，并上传 `dist/` 中的全部产物。
