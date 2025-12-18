#!/bin/bash

echo "========================================"
echo "构建所有平台的共享库 (32位和64位)"
echo "========================================"
echo ""

# 检测当前平台
PLATFORM=$(uname -s)

if [ "$PLATFORM" = "Darwin" ]; then
    echo "[1/3] 构建 macOS dylib..."
    chmod +x build_dylib.sh
    ./build_dylib.sh
    if [ $? -ne 0 ]; then
        echo "macOS dylib 构建失败！"
        exit 1
    fi
    echo ""
    
    echo "注意: Windows 和 Linux 版本需要在对应平台上构建"
    echo "  - Windows: build_dll.bat"
    echo "  - Linux: ./build_so.sh"
    
elif [ "$PLATFORM" = "Linux" ]; then
    echo "[1/3] 构建 Linux .so..."
    chmod +x build_so.sh
    ./build_so.sh
    if [ $? -ne 0 ]; then
        echo "Linux .so 构建失败！"
        exit 1
    fi
    echo ""
    
    echo "注意: Windows 和 macOS 版本需要在对应平台上构建"
    echo "  - Windows: build_dll.bat"
    echo "  - macOS: ./build_dylib.sh"
    
else
    echo "未知平台: $PLATFORM"
    echo "请手动运行对应的构建脚本:"
    echo "  - Windows: build_dll.bat"
    echo "  - macOS: ./build_dylib.sh"
    echo "  - Linux: ./build_so.sh"
    exit 1
fi

echo ""
echo "========================================"
echo "构建完成！"
echo "========================================"

