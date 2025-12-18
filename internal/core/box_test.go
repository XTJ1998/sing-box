package core

import (
	"fmt"
	"sort"
	"testing"

	ps "github.com/mitchellh/go-ps"
)

func TestBox(t *testing.T) {

	GetProcessList()

	//box := New(context.Background())
	//err := box.Start("香港(hysteria2)", false, nil)
	//if err != nil {
	//	t.Error(err)
	//}
	//defer func() { _ = box.Stop() }()
	//select {}
}

func GetProcessList() {
	processes, err := ps.Processes()
	if err != nil {
		fmt.Printf("获取进程列表失败: %v", err)
	}

	// 使用 map 去重
	nameMap := make(map[string]bool)
	for _, p := range processes {
		name := p.Executable()
		// 简单的过滤逻辑：过滤掉空名或显然是系统进程的（可选）
		if name == "" {
			continue
		}
		// Windows 下统一转为小写比较好控制，但显示时保留原样
		// 这里为了简单，直接存储原名
		nameMap[name] = true
	}

	// 转为 slice
	names := make([]string, 0, len(nameMap))
	for name := range nameMap {
		names = append(names, name)
	}

	// 排序，方便前端展示
	sort.Strings(names)

	for _, name := range names {
		fmt.Println(name)
	}

	//fmt.Printf("%v\n", names)

	//
	//// 序列化为 JSON
	//jsonData, err := json.Marshal(names)
	//if err != nil {
	//	return C.CString("[]")
	//}
	//
	//return C.CString(string(jsonData))
}
