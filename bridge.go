package main

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"playfast/internal/api"
	"playfast/internal/core"
	http_client "playfast/internal/http-client"
	"playfast/internal/node"
	"playfast/utils"
	"sync"
	"time"
	"unsafe"
)

var (
	box      *core.Box
	ctx      context.Context
	Version  = "1.0.0"
	initChan = make(chan struct{}) // 新增：用于通知初始化完成的通道
	initOnce sync.Once             // 新增：确保初始化只执行一次
)

//export Init
func Init() {
	// 使用 sync.Once 保证线程安全，防止多次调用 Init 导致重复关闭 channel 而 panic
	initOnce.Do(func() {
		ctx = context.Background()
		box = core.New(ctx)
		// 关闭通道，这会广播信号给所有阻塞在 <-initChan 的 goroutine
		close(initChan)
		log.Println("sing-box 核心模块初始化完成")
	})
}

// Switch 启动或停止加速
// status: true=启动, false=停止
// proxy: 代理节点名称
// route: 是否启用路由模式
// 返回: 错误信息，如果成功则返回空字符串
//
//export Switch
func Switch(status C.int, proxy *C.char, route C.int, appList *C.char) *C.char {
	// 优化后的等待逻辑
	select {
	case <-initChan:
		// 通道已关闭（即 Init 已完成），立即继续执行，无延迟
	case <-time.After(30 * time.Second): // 保持原有的 30s 超时逻辑
		// 超时处理：原代码超时后仍继续执行可能会导致空指针崩溃，这里建议直接返回错误
		return C.CString("错误：核心模块初始化超时或未调用 Init()")
	}

	proxyStr := C.GoString(proxy)
	statusBool := status != 0
	routeBool := route != 0

	var appNames []string
	appListStr := C.GoString(appList)

	if len(appListStr) > 0 {
		if err := json.Unmarshal([]byte(appListStr), &appNames); err != nil {
			log.Printf("解析应用列表失败: %v", err)
			// 如果解析失败，您可以选择返回错误，或者默认为空列表（即不加速任何程序）
			appNames = []string{}
		}
	}

	var err error
	if statusBool {
		if routeBool {
			ip, gateway, mask, err2 := utils.RandIP()
			if err2 != nil {
				errMsg := fmt.Sprintf("获取本地IP失败: %s", err2.Error())
				return C.CString(errMsg)
			}
			// 注意：这里不再显示对话框，Electron 端需要自己处理
			_ = ip
			_ = gateway
			_ = mask
		}

		if err = box.Start(proxyStr, routeBool, appNames); err != nil {
			return C.CString(err.Error())
		}
	} else {
		if err = box.Stop(); err != nil {
			return C.CString(err.Error())
		}
	}
	return C.CString("")
}

// GetProxyList 获取代理列表
// 返回: JSON 格式的代理列表字符串
//
//export GetProxyList
func GetProxyList() *C.char {
	proxies := node.Get()
	names := make([]string, 0, len(proxies))
	for _, proxy := range proxies {
		names = append(names, proxy.Name)
	}
	jsonData, err := json.Marshal(names)
	if err != nil {
		return C.CString("[]")
	}
	return C.CString(string(jsonData))
}

// GetAnnouncement 获取公告
// 返回: 公告内容字符串
//
//export GetAnnouncement
func GetAnnouncement() *C.char {
	all, err := http_client.GET(fmt.Sprintf("%s/announcement", api.GetApiDomain()))
	if err != nil {
		log.Printf("获取公告失败: %v", err)
		return C.CString("")
	}
	return C.CString(string(all))
}

// GetVersion 获取版本号
// 返回: 版本号字符串
//
//export GetVersion
func GetVersion() *C.char {
	return C.CString(Version)
}

// GetRouteInfo 获取路由配置信息（仅在启用路由模式时使用）
// 返回: JSON 格式的路由信息 {"ip": "...", "gateway": "...", "mask": "..."}
//
//export GetRouteInfo
func GetRouteInfo() *C.char {
	ip, gateway, mask, err := utils.RandIP()
	if err != nil {
		return C.CString("{}")
	}
	info := map[string]string{
		"ip":      ip,
		"gateway": gateway,
		"mask":    mask,
	}
	jsonData, err := json.Marshal(info)
	if err != nil {
		return C.CString("{}")
	}
	return C.CString(string(jsonData))
}

// FreeString 释放 C 字符串内存
//
//export FreeString
func FreeString(str *C.char) {
	if str != nil {
		C.free(unsafe.Pointer(str))
	}
}

// Stop 停止加速并清理资源
//
//export Stop
func Stop() {
	if box != nil {
		_ = box.Stop()
	}
}

// Cleanup 清理所有资源
//
//export Cleanup
func Cleanup() {
	Stop()
	box = nil
	ctx = nil
}
