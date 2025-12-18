# PlayFast Node.js 调用示例

使用 `koffi` 包调用 PlayFast 共享库的完整示例。

## 安装依赖

```bash
npm install
```

或者只安装 koffi：

```bash
npm install koffi
```

## 使用方法

### 方法 1: 使用 npm 脚本

```bash
npm start
```

或者：

```bash
npm run demo
```

### 方法 2: 直接运行

```bash
node node_playfast_koffi.js
```

## 文件说明

- `node_playfast_koffi.js` - 使用 koffi 调用共享库的完整示例
- `package.json` - 项目配置和依赖
- `playfast_amd64.dll` - Windows 64 位共享库（需要放在同一目录）
- `playfast_386.dll` - Windows 32 位共享库（可选）

## 功能说明

示例脚本会执行以下操作：

1. **Init()** - 初始化核心模块
2. **GetVersion()** - 获取版本号
3. **GetProxyList()** - 获取代理列表（JSON 格式）
4. **Switch()** - 启动加速（使用第一个节点）
5. **Stop()** - 停止加速
6. **GetAnnouncement()** - 获取公告（可选）
7. **Cleanup()** - 清理资源

## koffi 的优势

- ✅ **无需编译** - 纯 JavaScript 实现，不需要 C++ 编译环境
- ✅ **跨平台** - 支持 Windows、macOS、Linux
- ✅ **简单易用** - API 简洁明了
- ✅ **性能优秀** - 使用 JIT 编译优化
- ✅ **类型安全** - 支持 TypeScript

## 函数签名说明

```javascript
// 初始化
Init(): void

// 启动/停止加速
Switch(status: int, proxy: string, route: int): string
// status: 1=启动, 0=停止
// proxy: 代理节点名称
// route: 1=启用路由模式, 0=禁用
// 返回: 错误信息（成功返回空字符串）

// 获取代理列表（返回 JSON 字符串数组）
GetProxyList(): string

// 获取版本号
GetVersion(): string

// 获取公告
GetAnnouncement(): string

// 获取路由配置信息（JSON 格式）
GetRouteInfo(): string

// 停止加速
Stop(): void

// 清理资源
Cleanup(): void

// 释放字符串内存
FreeString(str: string): void
```

## 注意事项

1. **内存管理**: 所有返回字符串的函数都需要使用 `FreeString()` 释放内存
2. **初始化顺序**: 必须先调用 `Init()` 才能使用其他函数
3. **资源清理**: 程序退出前应调用 `Cleanup()` 清理资源
4. **权限要求**: Windows 下修改路由需要管理员权限
5. **库文件位置**: 确保 DLL/dylib/so 文件与脚本在同一目录，或使用绝对路径

## 平台支持

- ✅ Windows (x64, x86)
- ✅ macOS (x64, arm64)
- ✅ Linux (x64, x86, arm64)

## 故障排除

### 找不到库文件

确保共享库文件与脚本在同一目录，或者修改 `getLibraryPath()` 函数使用绝对路径。

### 初始化失败

检查库文件是否正确编译，以及是否缺少依赖的动态库。

### 内存泄漏

确保所有返回的字符串都调用了 `FreeString()` 释放内存。

