#!/bin/bash

echo "正在构建 Linux .so (32位和64位)..."

export CGO_ENABLED=1
export GOOS=linux

echo ""
echo "===== 构建 64位 .so ====="
export GOARCH=amd64
go build -tags with_clash_api -buildmode=c-shared -o playfast_amd64.so main_lib.go bridge.go
if [ $? -ne 0 ]; then
    echo "64位构建失败！"
    exit 1
fi
mv playfast.h playfast_amd64.h 2>/dev/null
echo "64位构建成功！输出文件: playfast_amd64.so, playfast_amd64.h"

echo ""
echo "===== 构建 32位 .so ====="
export GOARCH=386
go build -tags with_clash_api -buildmode=c-shared -o playfast_386.so main_lib.go bridge.go
if [ $? -ne 0 ]; then
    echo ""
    echo "警告: 32位构建失败！"
    echo "可能原因: 当前环境缺少32位 C 编译器工具链"
    echo "如果只需要64位版本，可以忽略此错误"
    echo "要构建32位版本，需要安装32位 gcc 工具链（如 gcc-multilib）"
    echo ""
else
    mv playfast.h playfast_386.h 2>/dev/null
    echo "32位构建成功！输出文件: playfast_386.so, playfast_386.h"
fi

echo ""
echo "===== 构建完成 ====="
if [ -f playfast_amd64.so ]; then
    echo "生成的文件:"
    echo "  - playfast_amd64.so (64位)"
    echo "  - playfast_amd64.h (64位头文件)"
fi
if [ -f playfast_386.so ]; then
    echo "  - playfast_386.so (32位)"
    echo "  - playfast_386.h (32位头文件)"
else
    echo "注意: 32位版本未生成，请检查是否安装了32位 C 编译器工具链"
fi

