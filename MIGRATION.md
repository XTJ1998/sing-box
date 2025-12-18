# 迁移说明

## 变更概述

项目已从 Wails 桌面应用重构为可编译为共享库（DLL/dylib/so）的核心加速模块，供 Electron 应用使用。

## 主要变更

### 1. 移除的文件和目录
- `frontend/` - 前端代码目录（已删除）
- `wails.json` - Wails 配置文件（已删除）
- `build.bat` - Wails 构建脚本（已删除）
- `main.go` - 原 Wails 主程序（已备份为 `main_wails.go.bak`）
- `app.go` - 原 Wails 应用逻辑（已备份为 `app_wails.go.bak`）

### 2. 新增的文件
- `bridge.go` - CGO 接口文件，导出 C 函数供 Electron 调用
- `main_lib.go` - 共享库构建入口文件
- `build_dll.bat` - Windows DLL 构建脚本
- `build_dylib.sh` - macOS dylib 构建脚本
- `build_so.sh` - Linux .so 构建脚本
- `README_LIB.md` - 库使用说明
- `electron_example.js` - Electron 集成示例代码

### 3. 保留的核心功能
- `internal/core/` - 核心加速功能（sing-box 集成）
- `internal/node/` - 代理节点管理
- `internal/api/` - API 配置
- `internal/http-client/` - HTTP 客户端
- `internal/path/` - 路径管理
- `utils/` - 工具函数（IP、路由等）

### 4. 不再使用的模块（保留但未使用）
- `internal/dialog/` - 对话框功能（依赖 Wails，不再使用）
- `internal/systray/` - 系统托盘功能（不再使用）

## 构建说明

### Windows
```bash
build_dll.bat
```
生成 `playfast.dll` 和 `playfast.h`

### macOS
```bash
chmod +x build_dylib.sh
./build_dylib.sh
```
生成 `playfast.dylib` 和 `playfast.h`

### Linux
```bash
chmod +x build_so.sh
./build_so.sh
```
生成 `playfast.so` 和 `playfast.h`

## 导出的 C 函数

所有函数都在 `bridge.go` 中定义：

1. `Init()` - 初始化核心模块
2. `Switch(status, proxy, route)` - 启动/停止加速
3. `GetProxyList()` - 获取代理列表（JSON）
4. `GetAnnouncement()` - 获取公告
5. `GetVersion()` - 获取版本号
6. `GetRouteInfo()` - 获取路由配置信息（JSON）
7. `Stop()` - 停止加速
8. `Cleanup()` - 清理资源
9. `FreeString(str)` - 释放 C 字符串内存

## Electron 集成

参考 `electron_example.js` 和 `README_LIB.md` 了解如何在 Electron 中使用。

需要安装依赖：
```bash
npm install ffi-napi ref-napi
```

## 注意事项

1. **依赖清理**：虽然 `go.mod` 中仍包含 Wails 依赖，但代码中已不再使用。可以运行 `go mod tidy` 清理，但某些间接依赖可能仍需要 Wails 的依赖项。

2. **内存管理**：所有返回的 C 字符串都需要使用 `FreeString()` 释放，避免内存泄漏。

3. **初始化顺序**：必须先调用 `Init()` 才能使用其他函数。

4. **资源清理**：程序退出前应调用 `Cleanup()` 清理资源。

5. **权限要求**：Windows 下修改路由需要管理员权限。

## 恢复 Wails 版本

如果需要恢复 Wails 版本，可以：
1. 恢复备份文件：`main_wails.go.bak` -> `main.go`，`app_wails.go.bak` -> `app.go`
2. 恢复 `wails.json`（如果有备份）
3. 恢复 `frontend/` 目录（如果有备份）

