@echo off
echo 正在构建 Windows DLL (32位和64位)...

set CGO_ENABLED=1
set GOOS=windows

echo.
echo ===== 构建 64位 DLL =====
set GOARCH=amd64
go build -tags "with_clash_api,with_gvisor,with_quic" -buildmode=c-shared -o playfast_amd64.dll main_lib.go bridge.go
if %ERRORLEVEL% NEQ 0 (
    echo 64位构建失败！
    exit /b 1
)
move playfast.h playfast_amd64.h >nul 2>&1
echo 64位构建成功！输出文件: playfast_amd64.dll, playfast_amd64.h

echo.
echo ===== 构建 32位 DLL =====
set GOARCH=386
go build -tags "with_clash_api,with_gvisor,with_quic" -buildmode=c-shared -o playfast_386.dll main_lib.go bridge.go
if %ERRORLEVEL% NEQ 0 (
    echo.
    echo 警告: 32位构建失败！
    echo 可能原因: 当前环境缺少32位 C 编译器工具链
    echo 如果只需要64位版本，可以忽略此错误
    echo 要构建32位版本，需要安装32位 MinGW 工具链
    echo.
) else (
    move playfast.h playfast_386.h >nul 2>&1
    echo 32位构建成功！输出文件: playfast_386.dll, playfast_386.h
)

echo.
echo ===== 构建完成 =====
if exist playfast_amd64.dll (
    echo 生成的文件:
    echo   - playfast_amd64.dll (64位)
    echo   - playfast_amd64.h (64位头文件)
)
if exist playfast_386.dll (
    echo   - playfast_386.dll (32位)
    echo   - playfast_386.h (32位头文件)
) else (
    echo 注意: 32位版本未生成，请检查是否安装了32位 C 编译器工具链
)

