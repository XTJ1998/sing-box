package main

import (
	"log"
	"os"
	"playfast/internal/path"
)

// 这个文件用于构建共享库（DLL/dylib）
// 实际的导出函数在 bridge.go 中定义
// 注意：即使构建共享库，Go 编译器仍然需要一个 main 函数作为入口点
func init() {
	// 设置日志
	appLog := path.Path() + "/app.log"
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	open, err := os.Create(appLog)
	if err == nil {
		log.SetOutput(open)
	}
}

// main 函数是必需的，即使构建共享库时不会被调用
// 实际的导出函数在 bridge.go 中通过 //export 指令定义
func main() {
	// 共享库模式下，这个函数不会被调用
	// 所有功能通过 CGO 导出的函数提供
}
