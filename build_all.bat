@echo off
echo ========================================
echo 构建所有平台的共享库 (32位和64位)
echo ========================================
echo.

echo [1/3] 构建 Windows DLL...
call build_dll.bat
if %ERRORLEVEL% NEQ 0 (
    echo Windows DLL 构建失败！
    exit /b 1
)
echo.

echo ========================================
echo 构建完成！
echo ========================================
echo.
echo 注意: macOS 和 Linux 版本需要在对应平台上构建
echo 使用以下命令:
echo   - macOS: chmod +x build_dylib.sh ^&^& ./build_dylib.sh
echo   - Linux: chmod +x build_so.sh ^&^& ./build_so.sh
echo.

