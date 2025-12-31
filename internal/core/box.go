package core

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"net/netip"
	"os"
	"playfast/internal/node"
	"playfast/internal/path"
	"playfast/utils"
	"sync"
	"time"

	box "github.com/sagernet/sing-box"
	"github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/experimental/deprecated"
	"github.com/sagernet/sing-box/include"
	slog "github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	dns "github.com/sagernet/sing-dns"
	"github.com/sagernet/sing/common/json/badoption"
	"github.com/sagernet/sing/service"
)

////go:embed geoip-cn.srs
//var geoip []byte
//
////go:embed geosite-cn.srs
//var geosite []byte
//
//// black 包含需要拦截的黑名单域名列表
////
////go:embed black-list.json
//var black []byte
//
//// direct 包含需要直连的白名单域名列表
////
////go:embed direct-list.json
//var direct []byte

type Box struct {
	box              *box.Box
	ctx              context.Context
	router           bool
	appends          []string
	defaultInterface int
	sync.Mutex
}

func (b *Box) Start(region string, router bool, apps []string) error {
	b.Lock()
	defer b.Unlock()
	b.router = router

	if b.box != nil {
		_ = b.box.Close()
		b.box = nil
	}

	if b.box == nil {
		err := b.newBox(region, apps)
		if err != nil {
			return err
		}
	}
	err := b.box.Start()
	if err != nil {
		return err
	}
	// 如果启用了系统路由配置
	if router {
		// 获取系统默认网络接口
		var defaultInterface *net.Interface
		defaultInterface, err = utils.GetDefaultInterface()
		if err != nil {
			// 如果获取默认接口失败，停止 box 并返回错误
			_ = b.box.Close()
			b.box = nil
			return err
		}
		// 记录默认接口索引
		b.defaultInterface = defaultInterface.Index
		// 启用 IP 转发
		err = utils.SetIPForwarding(b.defaultInterface, true)
		if err != nil {
			// 如果设置 IP 转发失败，停止 box 并返回错误
			_ = b.box.Close()
			b.box = nil
			return err
		}
	}
	// 配置系统路由规则
	err = route(b.appends)
	if err != nil {
		// 如果路由设置失败，停止 box 并返回错误
		_ = b.box.Close()
		b.box = nil
		return err
	}
	return nil
}

// Stop 停止代理服务并清理路由规则
func (b *Box) Stop() error {
	// 加锁确保并发安全
	b.Lock()
	defer b.Unlock()
	if b.box == nil {
		return nil
	}

	// 关闭 box 实例
	err := b.box.Close()
	// 清空 box 指针
	b.box = nil

	// 如果启用了系统路由配置
	if b.router {
		// 禁用 IP 转发
		err = utils.SetIPForwarding(b.defaultInterface, false)
		if err != nil {
			return err
		}
	}

	// 清理系统路由规则
	deleteRoute(b.appends)
	return err
}

// New 创建一个新的 Box 实例
func New(ctx context.Context) *Box {
	// 为上下文添加标准错误管理器
	ctx = service.ContextWith(ctx, deprecated.NewStderrManager(slog.StdLogger()))
	// 为上下文添加 sing-box 注册表信息
	ctx = box.Context(ctx,
		include.InboundRegistry(),
		include.OutboundRegistry(),
		include.EndpointRegistry(),
		include.DNSTransportRegistry(),
		include.ServiceRegistry())

	// 将嵌入的配置文件写入到应用程序目录
	//_ = os.WriteFile(filepath.Join(path.Path(), "black-list.json"), black, 0644)
	//_ = os.WriteFile(filepath.Join(path.Path(), "direct-list.json"), direct, 0644)
	//_ = os.WriteFile(filepath.Join(path.Path(), "geoip-cn.srs"), geoip, 0644)
	//_ = os.WriteFile(filepath.Join(path.Path(), "geosite-cn.srs"), geosite, 0644)

	// 创建 Box 实例
	b := Box{
		ctx:     ctx,
		appends: []string{},
	}

	// 异步更新配置文件
	//go b.update()

	return &b
}

// update 从远程服务器更新配置文件
// 该方法会异步执行，尝试下载最新的配置文件，如果下载失败则使用嵌入的默认配置
//func (b *Box) update() {
//	var data []byte
//	var err error
//
//	// 更新黑名单配置文件
//	data, err = httpclient.GET(fmt.Sprintf("%s/black-list.json", api.GetApiDomain()))
//	if err != nil {
//		// 如果下载失败，使用嵌入的默认配置
//		data = black
//	}
//	b.Lock()
//	_ = os.WriteFile(filepath.Join(path.Path(), "black-list.json"), data, 0644)
//	b.Unlock()
//
//	// 更新白名单配置文件
//	data, err = httpclient.GET(fmt.Sprintf("%s/direct-list.json", api.GetApiDomain()))
//	if err != nil {
//		// 如果下载失败，使用嵌入的默认配置
//		data = direct
//	}
//	b.Lock()
//	_ = os.WriteFile(filepath.Join(path.Path(), "direct-list.json"), data, 0644)
//	b.Unlock()
//
//	// 更新中国IP地址数据库
//	data, err = httpclient.GET("https://raw.githubusercontent.com/SagerNet/sing-geoip/rule-set/geoip-cn.srs")
//	if err != nil {
//		// 如果下载失败，使用嵌入的默认配置
//		data = geoip
//	}
//	b.Lock()
//	_ = os.WriteFile(filepath.Join(path.Path(), "geoip-cn.srs"), data, 0644)
//	b.Unlock()
//
//	// 更新中国网站域名数据库
//	data, err = httpclient.GET("https://raw.githubusercontent.com/SagerNet/sing-geosite/rule-set/geosite-cn.srs")
//	if err != nil {
//		// 如果下载失败，使用嵌入的默认配置
//		data = geosite
//	}
//	b.Lock()
//	_ = os.WriteFile(filepath.Join(path.Path(), "geosite-cn.srs"), data, 0644)
//	b.Unlock()
//}

// newBox 创建并配置一个新的 sing-box 实例
func (b *Box) newBox(proxy string, apps []string) error {
	// 获取代理出站配置
	proxyOutbound, proxyOutboundHost, err := node.GetOutbound(proxy)
	if err != nil {
		return err
	}

	// 初始化路由白名单
	b.appends = make([]string, 0)

	// 获取代理服务器的IP地址
	proxyOutboundIp, err := utils.GetIPsFromString(proxyOutboundHost)
	if err != nil {
		return err
	}

	// 将代理服务器IP添加到路由白名单
	b.appends = append(b.appends, fmt.Sprintf("%s/32", proxyOutboundIp))
	// 创建 sing-box 配置选项
	options := box.Options{
		Options: option.Options{
			// 日志配置
			Log: &option.LogOptions{
				// 禁用日志记录
				Disabled: true,
			},
			// DNS配置
			DNS: &option.DNSOptions{
				RawDNSOptions: option.RawDNSOptions{
					// DNS服务器列表
					Servers: []option.DNSServerOptions{
						// 代理DNS服务器：使用Cloudflare DNS，通过代理访问
						{
							Type: "https",
							Tag:  "proxyDns",
							Options: &option.RemoteHTTPSDNSServerOptions{
								RemoteTLSDNSServerOptions: option.RemoteTLSDNSServerOptions{
									RemoteDNSServerOptions: option.RemoteDNSServerOptions{
										LocalDNSServerOptions: option.LocalDNSServerOptions{
											DialerOptions: option.DialerOptions{
												// 通过代理访问
												Detour: "proxy",
											},
										},
										DNSServerAddressOptions: option.DNSServerAddressOptions{
											// Cloudflare DNS服务器
											Server:     "cloudflare-dns.com",
											ServerPort: 443,
										},
									},
								},
							},
						},
						// 本地DNS服务器：使用阿里云DNS，直接访问
						{
							Type: "https",
							Tag:  "localDns",
							Options: &option.RemoteHTTPSDNSServerOptions{
								RemoteTLSDNSServerOptions: option.RemoteTLSDNSServerOptions{
									RemoteDNSServerOptions: option.RemoteDNSServerOptions{
										DNSServerAddressOptions: option.DNSServerAddressOptions{
											// 阿里云DNS服务器
											Server:     "223.5.5.5",
											ServerPort: 443,
										},
									},
								},
							},
						},
					},
					// DNS路由规则
					Rules: []option.DNSRule{
						// 中国网站使用本地DNS
						//{
						//	Type: constant.RuleTypeDefault,
						//	DefaultOptions: option.DefaultDNSRule{
						//		RawDefaultDNSRule: option.RawDefaultDNSRule{
						//			RuleSet: []string{
						//				"geosite-cn",
						//			},
						//		},
						//		DNSRuleAction: option.DNSRuleAction{
						//			Action: constant.RuleActionTypeRoute,
						//			RouteOptions: option.DNSRouteActionOptions{
						//				Server: "localDns",
						//			},
						//		},
						//	},
						//},
						// 代理服务器使用本地DNS
						{
							Type: constant.RuleTypeDefault,
							DefaultOptions: option.DefaultDNSRule{
								RawDefaultDNSRule: option.RawDefaultDNSRule{
									Domain: []string{proxyOutboundHost},
								},
								DNSRuleAction: option.DNSRuleAction{
									Action: constant.RuleActionTypeRoute,
									RouteOptions: option.DNSRouteActionOptions{
										Server: "localDns",
									},
								},
							},
						},
					},
					// 默认DNS服务器
					Final: "proxyDns",
					// DNS客户端选项
					DNSClientOptions: option.DNSClientOptions{
						// 域名解析策略：优先使用IPv4
						Strategy: option.DomainStrategy(dns.DomainStrategyUseIPv4),
						// DNS缓存容量
						CacheCapacity: 2048,
					},
				},
			},
			// 入站配置
			Inbounds: []option.Inbound{
				// TUN接口配置
				{
					Type: constant.TypeTun,
					Tag:  "tun-in",
					Options: &option.TunInboundOptions{
						// TUN接口名称
						InterfaceName: "utun25",
						// 最大传输单元
						MTU: 1500,
						// TUN接口IP地址
						Address: badoption.Listable[netip.Prefix]{
							netip.MustParsePrefix("172.25.0.0/30"),
						},
						// 注释掉的路由地址配置
						//RouteAddress: in(),
						// 注释掉的自动路由配置
						//AutoRoute:    true,
						// 注释掉的严格路由配置
						//StrictRoute:  true,
						// UDP超时时间
						UDPTimeout: option.UDPTimeoutCompat(time.Second * 300),
						// 使用gvisor作为TUN栈
						Stack: "gvisor",
					},
				},
			},
			// 路由配置
			Route: &option.RouteOptions{
				// 路由规则集合
				RuleSet: nil,
				//RuleSet: []option.RuleSet{
				//	// 中国网站域名集合
				//	{
				//		Type:         constant.RuleSetTypeLocal,
				//		Tag:          "geosite-cn",
				//		Format:       constant.RuleSetFormatBinary,
				//		LocalOptions: option.LocalRuleSet{Path: filepath.Join(path.Path(), "geosite-cn.srs")},
				//	},
				//	// 中国IP地址集合
				//	{
				//		Type:         constant.RuleSetTypeLocal,
				//		Tag:          "geoip-cn",
				//		Format:       constant.RuleSetFormatBinary,
				//		LocalOptions: option.LocalRuleSet{Path: filepath.Join(path.Path(), "geoip-cn.srs")},
				//	},
				//	// 黑名单域名集合
				//	{
				//		Type:         constant.RuleSetTypeLocal,
				//		Tag:          "black-list",
				//		Format:       constant.RuleSetFormatSource,
				//		LocalOptions: option.LocalRuleSet{Path: filepath.Join(path.Path(), "black-list.json")},
				//	},
				//	// 白名单域名集合
				//	{
				//		Type:         constant.RuleSetTypeLocal,
				//		Tag:          "direct-list",
				//		Format:       constant.RuleSetFormatSource,
				//		LocalOptions: option.LocalRuleSet{Path: filepath.Join(path.Path(), "direct-list.json")},
				//	},
				//},
				// 自动检测网络接口
				AutoDetectInterface: true,
				// 初始化空的路由规则列表
				Rules: []option.Rule{},
			}, // 出站配置：定义流量的出口方式
			Outbounds: []option.Outbound{
				// 代理出站：使用从node.GetOutbound获取的代理配置
				*proxyOutbound,
				// 直连出站：直接连接到目标服务器
				{Type: constant.TypeDirect, Tag: "direct"},
			},
			// 实验性功能配置
			Experimental: &option.ExperimentalOptions{
				// Clash API配置：提供与Clash兼容的API接口
				ClashAPI: &option.ClashAPIOptions{
					// 外部控制器地址：用于Clash客户端连接和管理
					ExternalController: "127.0.0.1:54713",
				},
			},
		},
		Context: b.ctx,
	}
	options.Options.Route.Rules = append(options.Options.Route.Rules, []option.Rule{
		{
			Type: constant.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RawDefaultRule: option.RawDefaultRule{
					Network: []string{
						"udp",
					},
					Port: []uint16{
						443,
					},
				},
				RuleAction: option.RuleAction{
					Action: constant.RuleActionTypeReject,
					RejectOptions: option.RejectActionOptions{
						Method: "default",
						NoDrop: false,
					},
				},
			},
		}, //禁止http3
		{
			Type: constant.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RawDefaultRule: option.RawDefaultRule{
					Invert: true,
				},
				RuleAction: option.RuleAction{
					Action: constant.RuleActionTypeSniff,
					SniffOptions: option.RouteActionSniff{
						Sniffer: []string{
							"dns", "http", "tls", "quic",
						},
					},
				},
			},
		}, //解析协议域名
		//{
		//	Type: constant.RuleTypeDefault,
		//	DefaultOptions: option.DefaultRule{
		//		RawDefaultRule: option.RawDefaultRule{
		//			RuleSet: []string{"black-list"},
		//		},
		//		RuleAction: option.RuleAction{
		//			Action: constant.RuleActionTypeReject,
		//			RejectOptions: option.RejectActionOptions{
		//				Method: constant.RuleActionRejectMethodDefault,
		//				NoDrop: false,
		//			},
		//		},
		//	},
		//}, //过滤黑名单
		//{
		//	Type: constant.RuleTypeDefault,
		//	DefaultOptions: option.DefaultRule{
		//		RawDefaultRule: option.RawDefaultRule{
		//			RuleSet: []string{"direct-list"},
		//		},
		//		RuleAction: option.RuleAction{
		//			Action: constant.RuleActionTypeRoute,
		//			RouteOptions: option.RouteActionOptions{
		//				Outbound: "direct",
		//			},
		//		},
		//	},
		//}, //直连白名单
		//{
		//	Type: constant.RuleTypeDefault,
		//	DefaultOptions: option.DefaultRule{
		//		RawDefaultRule: option.RawDefaultRule{
		//			Protocol: []string{"dns"},
		//		},
		//		RuleAction: option.RuleAction{
		//			Action: constant.RuleActionTypeHijackDNS,
		//		},
		//	},
		//}, //dns劫持
		//{
		//	Type: constant.RuleTypeDefault,
		//	DefaultOptions: option.DefaultRule{
		//		RawDefaultRule: option.RawDefaultRule{
		//			RuleSet: []string{
		//				"geosite-cn", "geoip-cn",
		//			},
		//		},
		//		RuleAction: option.RuleAction{
		//			Action: constant.RuleActionTypeRoute,
		//			RouteOptions: option.RouteActionOptions{
		//				Outbound: "direct",
		//			},
		//		},
		//	},
		//},
	}...)
	// 2. [新增] 动态应用程序规则
	if len(apps) > 0 {
		appRule := option.Rule{
			Type: constant.RuleTypeDefault,
			DefaultOptions: option.DefaultRule{
				RawDefaultRule: option.RawDefaultRule{
					ProcessName: apps, // 使用传入的程序列表
				},
				RuleAction: option.RuleAction{
					Action: constant.RuleActionTypeRoute,
					RouteOptions: option.RouteActionOptions{
						Outbound: "proxy", // 指定程序 -> 走代理
					},
				},
			},
		}
		options.Options.Route.Rules = append(options.Options.Route.Rules, appRule)
	}

	// 最终路由规则：匹配所有剩余流量
	// 这是默认路由规则，所有未被其他规则匹配的流量都将使用此规则
	finalRule := option.Rule{
		Type: constant.RuleTypeDefault,
		DefaultOptions: option.DefaultRule{
			RawDefaultRule: option.RawDefaultRule{
				Invert: true, // 反转匹配，即匹配所有未被其他规则匹配的流量
			},
			RuleAction: option.RuleAction{
				Action: constant.RuleActionTypeRoute,
				RouteOptions: option.RouteActionOptions{
					Outbound: "direct", // [重点] 默认为直连，确保未明确指定的流量不会自动走代理
				},
			},
		},
	}
	// 将最终路由规则添加到规则列表
	options.Options.Route.Rules = append(options.Options.Route.Rules, finalRule)

	// 更新日志配置：启用日志记录到文件
	// 首先删除旧的日志文件
	_ = os.Remove(path.Path() + "/run.log")
	// 配置新的日志选项
	options.Log = &option.LogOptions{
		Disabled:     false,                            // 启用日志记录
		Level:        slog.FormatLevel(slog.LevelInfo), // 日志级别设置为Info
		Output:       path.Path() + "/run.log",         // 日志输出到指定文件
		Timestamp:    true,                             // 记录日志时间戳
		DisableColor: true,                             // 禁用彩色日志输出
	}
	b.box, err = box.New(options)
	return err
}
