# 构建说明

## 架构支持

所有构建脚本都会尝试构建 32 位和 64 位版本：

- **64位 (amd64)**: `playfast_amd64.dll/dylib/so` + `playfast_amd64.h`
- **32位 (386)**: `playfast_386.dll/dylib/so` + `playfast_386.h`

## 32位构建注意事项

### Windows
- 需要安装 32 位 MinGW 工具链
- 如果只有 64 位工具链，32 位构建会失败，但会给出警告
- 64 位版本仍然可以正常使用

### macOS
- **macOS 10.15 (Catalina) 及更高版本不再支持 32 位应用**
- 32 位构建通常会失败，这是正常的
- 建议只使用 64 位版本

### Linux
- 需要安装 32 位 gcc 工具链
- Ubuntu/Debian: `sudo apt-get install gcc-multilib`
- CentOS/RHEL: `sudo yum install glibc-devel.i686`
- 如果只有 64 位工具链，32 位构建会失败，但会给出警告

## 构建命令

### Windows
```bash
build_dll.bat
```

### macOS
```bash
chmod +x build_dylib.sh
./build_dylib.sh
```

### Linux
```bash
chmod +x build_so.sh
./build_so.sh
```

## Electron 集成

在 Electron 应用中，根据运行时的架构自动选择对应的库文件：

```javascript
const os = require('os');
const arch = os.arch(); // 'x64' 或 'ia32'

// Windows
const libName = arch === 'x64' ? 'playfast_amd64.dll' : 'playfast_386.dll';

// macOS
const libName = arch === 'x64' ? 'playfast_amd64.dylib' : 'playfast_386.dylib';

// Linux
const libName = arch === 'x64' ? 'playfast_amd64.so' : 'playfast_386.so';
```

参考 `electron_example.js` 查看完整示例。

