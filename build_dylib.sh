#!/bin/bash

echo "正在构建 macOS dylib (32位和64位)..."

export CGO_ENABLED=1
export GOOS=darwin

echo ""
echo "===== 构建 64位 dylib ====="
export GOARCH=amd64
go build -tags with_clash_api -buildmode=c-shared -o playfast_amd64.dylib main_lib.go bridge.go
if [ $? -ne 0 ]; then
    echo "64位构建失败！"
    exit 1
fi
mv playfast.h playfast_amd64.h 2>/dev/null
echo "64位构建成功！输出文件: playfast_amd64.dylib, playfast_amd64.h"

echo ""
echo "===== 构建 32位 dylib ====="
export GOARCH=386
go build -tags with_clash_api -buildmode=c-shared -o playfast_386.dylib main_lib.go bridge.go
if [ $? -ne 0 ]; then
    echo ""
    echo "警告: 32位构建失败！"
    echo "可能原因: macOS 10.15+ 已不再支持32位应用"
    echo "如果只需要64位版本，可以忽略此错误"
    echo ""
else
    mv playfast.h playfast_386.h 2>/dev/null
    echo "32位构建成功！输出文件: playfast_386.dylib, playfast_386.h"
fi

echo ""
echo "===== 构建完成 ====="
if [ -f playfast_amd64.dylib ]; then
    echo "生成的文件:"
    echo "  - playfast_amd64.dylib (64位)"
    echo "  - playfast_amd64.h (64位头文件)"
fi
if [ -f playfast_386.dylib ]; then
    echo "  - playfast_386.dylib (32位)"
    echo "  - playfast_386.h (32位头文件)"
else
    echo "注意: 32位版本未生成（macOS 10.15+ 不支持32位应用）"
fi

