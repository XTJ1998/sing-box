#!/bin/bash

echo "正在构建 macOS dylib (32位和64位)..."

export CGO_ENABLED=1
export GOOS=darwin

:: ==========================================
:: 配置 MSYS2 编译器路径
:: ==========================================
set CC_64=D:\mysys64\mingw64\bin\gcc.exe
set CC_32=D:\mysys64\mingw32\bin\gcc.exe


echo ""
echo "===== 构建 64位 dylib ====="
export GOARCH=amd64
:: 强制指定 64位 编译器
set CC=%CC_64%
go build -tags with_clash_api -buildmode=c-shared -o signBox64.dylib main_lib.go bridge.go
if [ $? -ne 0 ]; then
    echo "64位构建失败！"
    exit 1
fi
mv playfast.h signBox64.h 2>/dev/null
echo "64位构建成功！输出文件: signBox64.dylib, signBox64.h"

echo ""
echo "===== 构建 32位 dylib ====="
export GOARCH=386
set GOARCH=386
go build -tags with_clash_api -buildmode=c-shared -o signBox32.dylib main_lib.go bridge.go
if [ $? -ne 0 ]; then
    echo ""
    echo "警告: 32位构建失败！"
    echo "可能原因: macOS 10.15+ 已不再支持32位应用"
    echo "如果只需要64位版本，可以忽略此错误"
    echo ""
else
    mv playfast.h signBox32.h 2>/dev/null
    echo "32位构建成功！输出文件: signBox32.dylib, signBox32.h"
fi

echo ""
echo "===== 构建完成 ====="
:: 清理环境变量
set CC=
if [ -f signBox64.dylib ]; then
    echo "生成的文件:"
    echo "  - signBox64.dylib (64位)"
    echo "  - signBox64.h (64位头文件)"
fi
if [ -f signBox32.dylib ]; then
    echo "  - signBox32.dylib (32位)"
    echo "  - signBox32.h (32位头文件)"
else
    echo "注意: 32位版本未生成（macOS 10.15+ 不支持32位应用）"
fi

