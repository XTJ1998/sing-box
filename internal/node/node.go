package node

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"playfast/internal/api"
	"playfast/internal/echo"
	"playfast/internal/http-client"
	"time"

	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/include"
	slog "github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/json/badoption"
)

type Proxy struct {
	Name     string `json:"name"`
	Method   string `json:"method"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
}

func Get() []Proxy {
	data := make([]Proxy, 0)
	all, err := http_client.GET(fmt.Sprintf("%s/proxy.json", api.GetApiDomain()))
	if err != nil {
		return data
	}
	_ = json.Unmarshal(all, &data)
	return data
}

func GetOutbound(proxy string) (*option.Outbound, string, error) {
	data := Get()
	var bestOut *option.Outbound
	var bestHost string
	// 初始化最小延迟为一个较大的数，或者通过逻辑判断
	var minLatency int64 = 99999999

	for i, p := range data {
		var out option.Outbound
		switch p.Protocol {
		case "shadowsocks":
			out = option.Outbound{
				Type: constant.TypeShadowsocks,
				Tag:  "proxy",
				Options: &option.ShadowsocksOutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     p.Host,
						ServerPort: p.Port,
					},
					Method:   p.Method,
					Password: p.Password,
					UDPOverTCP: &option.UDPOverTCPOptions{
						Enabled: true,
						Version: 2,
					},
				},
			}
		case "vless":
			out = option.Outbound{
				Type: constant.TypeVLESS,
				Tag:  "proxy",
				Options: &option.VLESSOutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     p.Host,
						ServerPort: p.Port,
					},
					UUID:                        p.Password,
					OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{},
					Multiplex: &option.OutboundMultiplexOptions{
						Enabled:        true,
						Protocol:       "h2mux",
						MaxConnections: 8,
						MinStreams:     16,
						Padding:        false,
					},
				},
			}
		case "socks":
			out = option.Outbound{
				Type: constant.TypeSOCKS,
				Tag:  "proxy",
				Options: &option.SOCKSOutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     p.Host,
						ServerPort: p.Port,
					},
					Version:  "5",
					Username: "playfast",
					Password: p.Password,
					UDPOverTCP: &option.UDPOverTCPOptions{
						Enabled: true,
						Version: 2,
					},
				},
			}
		case "hysteria2":
			out = option.Outbound{
				Type: constant.TypeHysteria2,
				Tag:  "proxy",
				Options: &option.Hysteria2OutboundOptions{
					ServerOptions: option.ServerOptions{
						Server:     p.Host,
						ServerPort: p.Port,
					},
					Password: p.Password,
					OutboundTLSOptionsContainer: option.OutboundTLSOptionsContainer{
						TLS: &option.OutboundTLSOptions{
							Enabled:    true,
							ServerName: "gpp",
							Insecure:   true,
							ALPN:       badoption.Listable[string]{"h3"},
						},
					},
				},
			}
		default:
			continue
		}
		registryOut := include.OutboundRegistry()
		createOutbound, err := registryOut.CreateOutbound(context.Background(), nil, slog.StdLogger(), out.Type, out.Type, out.Options)
		if err != nil {
			continue
		}
		client := echo.NewClient("1.1.1.1:80", echo.WithTimeout(3*time.Second), echo.WithDialer(createOutbound.DialContext))
		err = client.Connect(context.Background())
		if err != nil {
			log.Println(fmt.Sprintf("节点连接失败:ID:%d 节点:%s，UUID:%s 错误:%v", i, p.Name, p.Password, err))
			continue
		}
		var ms int64
		result := client.Test(context.Background(), []byte("GET / HTTP/1.1\r\nHost: 1.1.1.1\r\nAccept: *\r\n\r\n\r\n"))
		ms = result.Latency.Milliseconds()
		log.Println(fmt.Sprintf("节点检测:ID:%d 节点:%s 延迟=%dms", i, p.Name, ms))

		if ms > 0 && ms < minLatency {
			minLatency = ms
			// 必须复制一份 out，否则循环变量可能会被覆盖（但在 switch case 重新赋值结构体一般是安全的，这里为了保险起见直接保存需要的值）
			currentOut := out
			bestOut = &currentOut
			bestHost = p.Host
		}
	}

	if bestOut != nil {
		log.Println(fmt.Sprintf("最终选择节点: %s, 延迟: %dms", bestHost, minLatency))
		return bestOut, bestHost, nil
	}

	return nil, "", errors.New("not found valid Outbound")
}
