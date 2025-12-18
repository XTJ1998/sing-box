@echo off
echo 正在构建 Windows DLL (32位和64位)...

set CGO_ENABLED=1
set GOOS=windows

:: ==========================================
:: 配置 MSYS2 编译器路径
:: ==========================================
set CC_64=D:\mysys64\mingw64\bin\gcc.exe
set CC_32=D:\mysys64\mingw32\bin\gcc.exe

echo.
echo ===== 构建 64位 DLL =====
set GOARCH=amd64
:: 强制指定 64位 编译器
set CC=%CC_64%
go build -tags "with_clash_api,with_gvisor,with_quic" -buildmode=c-shared -o signBox64.dll main_lib.go bridge.go
if %ERRORLEVEL% NEQ 0 (
    echo 64位构建失败！请检查 %CC_64% 是否存在。
    exit /b 1
)
move playfast.h signBox64.h >nul 2>&1
echo 64位构建成功！输出文件: signBox64.dll, signBox64.h

echo.
echo ===== 构建 32位 DLL =====
set GOARCH=386
:: 强制指定 32位 编译器
set CC=%CC_32%
go build -tags "with_clash_api,with_gvisor,with_quic" -buildmode=c-shared -o signBox32.dll main_lib.go bridge.go
if %ERRORLEVEL% NEQ 0 (
    echo.
        echo 警告: 32位构建失败！
        echo 请检查路径是否正确: %CC_32%
        echo 以及该路径下是否存在 gcc.exe
        echo.
) else (
    move playfast.h signBox32.h >nul 2>&1
    echo 32位构建成功！输出文件: signBox32.dll, signBox32.h
)

echo.
echo ===== 构建完成 =====
:: 清理环境变量
set CC=
if exist signBox64.dll (
    echo 生成的文件:
    echo   - signBox64.dll (64位)
    echo   - signBox64.h (64位头文件)
)
if exist signBox32.dll (
    echo   - signBox32.dll (32位)
    echo   - signBox32.h (32位头文件)
) else (
    echo 注意: 32位版本未生成，请检查是否安装了32位 C 编译器工具链1
)

