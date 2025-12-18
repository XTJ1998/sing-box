# 权限说明

## Windows 权限要求

在 Windows 系统上运行 PlayFast 共享库时，**需要管理员权限**才能正常工作。

### 为什么需要管理员权限？

1. **创建 TUN 接口**：需要管理员权限来创建虚拟网络接口
2. **修改路由表**：需要管理员权限来添加/删除系统路由
3. **修改网络接口设置**：需要管理员权限来修改网络接口的指标和 IP 转发设置

### 如何以管理员权限运行？

#### 方法 1: 以管理员身份运行 Node.js

1. 找到 Node.js 可执行文件（通常在 `C:\Program Files\nodejs\node.exe`）
2. 右键点击，选择"以管理员身份运行"
3. 在管理员权限的命令行窗口中运行你的脚本：

```bash
node node_playfast_koffi.js
```

#### 方法 2: 以管理员身份运行命令行/终端

1. 在开始菜单搜索 "cmd" 或 "PowerShell"
2. 右键点击，选择"以管理员身份运行"
3. 在管理员权限的命令行中切换到项目目录：

```bash
cd D:\go-object\self\playfast
node node_playfast_koffi.js
```

#### 方法 3: 通过 npm 脚本（需要管理员权限的命令行）

```bash
# 在管理员权限的命令行中
npm start
```

### 常见错误

#### 错误 1: `Access is denied`

```
✗ 启动失败: start inbound/tun[tun-in]: configure tun interface: Access is denied.
```

**解决方案**：以管理员权限运行 Node.js 或命令行。

#### 错误 2: 路由操作失败

如果看到路由相关的错误，通常也是权限问题。请确保：

1. 以管理员身份运行
2. 关闭可能冲突的 VPN 软件
3. 确保没有其他程序占用网络接口

### Electron 应用中的处理

如果你的 Electron 应用需要自动请求管理员权限，可以：

1. **检测权限**：使用 `is-elevated` 等 npm 包检测是否有管理员权限
2. **提示用户**：如果没有权限，提示用户以管理员身份运行
3. **自动提升权限**：在 Windows 上可以尝试使用 `sudo-prompt` 或 `node-windows` 包

示例代码：

```javascript
const isElevated = require('is-elevated');

async function checkPermissions() {
  const elevated = await isElevated();
  if (!elevated) {
    console.error('需要管理员权限！请以管理员身份运行。');
    process.exit(1);
  }
}

// 在调用 Init() 之前检查权限
await checkPermissions();
playfast.Init();
```

### macOS/Linux 权限

- **macOS**: 通常需要用户密码来授权网络设置
- **Linux**: 通常需要 root 权限或使用 `sudo` 运行

