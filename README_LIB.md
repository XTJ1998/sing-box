# PlayFast 核心库

这是 PlayFast 的核心加速功能库，已移除 Wails 前端代码，可以编译为共享库（DLL/dylib/so）供 Electron 应用使用。

## 构建

所有构建脚本都会同时生成 32 位和 64 位版本。

### Windows DLL
```bash
build_dll.bat
```
将生成：
- `playfast_amd64.dll` 和 `playfast_amd64.h` (64位)
- `playfast_386.dll` 和 `playfast_386.h` (32位)

### macOS dylib
```bash
chmod +x build_dylib.sh
./build_dylib.sh
```
将生成：
- `playfast_amd64.dylib` 和 `playfast_amd64.h` (64位)
- `playfast_386.dylib` 和 `playfast_386.h` (32位)

**注意**: 现代 macOS 可能不支持 32 位应用，32 位构建可能会失败。

### Linux .so
```bash
chmod +x build_so.sh
./build_so.sh
```
将生成：
- `playfast_amd64.so` 和 `playfast_amd64.h` (64位)
- `playfast_386.so` 和 `playfast_386.h` (32位)

### 一键构建（当前平台）
- Windows: `build_all.bat`
- macOS/Linux: `chmod +x build_all.sh && ./build_all.sh`

## 导出的 C 函数

### Init()
初始化核心模块，必须在调用其他函数之前调用。

### Switch(status, proxy, route) -> error
启动或停止加速。
- `status`: 1=启动, 0=停止
- `proxy`: 代理节点名称（C 字符串）
- `route`: 1=启用路由模式, 0=禁用
- 返回: 错误信息（C 字符串），成功返回空字符串

### GetProxyList() -> json
获取代理列表，返回 JSON 格式的字符串数组。

### GetAnnouncement() -> string
获取公告内容。

### GetVersion() -> string
获取版本号。

### GetRouteInfo() -> json
获取路由配置信息（仅在启用路由模式时使用），返回 JSON 格式：`{"ip": "...", "gateway": "...", "mask": "..."}`

### Stop()
停止加速并清理资源。

### Cleanup()
清理所有资源。

### FreeString(str)
释放 C 字符串内存，用于释放上述函数返回的字符串。

## Electron 集成示例

### 使用 ffi-napi

```javascript
const ffi = require('ffi-napi');
const ref = require('ref-napi');
const os = require('os');
const path = require('path');

// 根据平台和架构选择库文件
function getLibraryPath() {
  const platform = os.platform();
  const arch = os.arch();
  
  let libName;
  if (platform === 'win32') {
    libName = arch === 'x64' ? 'playfast_amd64.dll' : 'playfast_386.dll';
  } else if (platform === 'darwin') {
    libName = arch === 'x64' ? 'playfast_amd64.dylib' : 'playfast_386.dylib';
  } else {
    libName = arch === 'x64' ? 'playfast_amd64.so' : 'playfast_386.so';
  }
  
  return path.join(__dirname, libName);
}

// 定义函数签名
const playfast = ffi.Library(getLibraryPath(), {
  'Init': ['void', []],
  'Switch': ['string', ['int', 'string', 'int']],
  'GetProxyList': ['string', []],
  'GetAnnouncement': ['string', []],
  'GetVersion': ['string', []],
  'GetRouteInfo': ['string', []],
  'Stop': ['void', []],
  'Cleanup': ['void', []],
  'FreeString': ['void', ['string']]
});

// 初始化
playfast.Init();

// 获取代理列表
const proxyListJson = playfast.GetProxyList();
const proxyList = JSON.parse(proxyListJson);
playfast.FreeString(proxyListJson);

// 启动加速
const error = playfast.Switch(1, '节点名称', 0);
if (error) {
  console.error('启动失败:', error);
  playfast.FreeString(error);
}

// 停止加速
playfast.Stop();

// 清理
playfast.Cleanup();
```

### 使用 node-ffi-napi (推荐)

```bash
npm install ffi-napi ref-napi
```

## 注意事项

1. 在调用任何函数之前，必须先调用 `Init()`
2. 所有返回的字符串都需要使用 `FreeString()` 释放内存
3. 程序退出前应调用 `Cleanup()` 清理资源
4. Windows 下需要管理员权限才能修改路由
5. 路由模式需要将设备连接到同一网络

