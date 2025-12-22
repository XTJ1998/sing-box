package core

import (
	"bytes"
	"log"
	"net"
	"net/netip"
	"playfast/utils"

	"github.com/sagernet/sing-box/common/srs"
	"github.com/sagernet/sing/common/json/badoption"
)

func routeIps(appends []string) badoption.Listable[netip.Prefix] {
	prefixes := make([]netip.Prefix, 0)
	read, err := srs.Read(bytes.NewBuffer(geoip), false)
	if err != nil {
		return prefixes
	}
	for _, rule := range read.Options.Rules {
		if rule.DefaultOptions.IPSet != nil {
			for _, ipRange := range rule.DefaultOptions.IPSet.Ranges() {
				for _, prefix := range ipRange.Prefixes() {
					if !prefix.Addr().Is4() {
						continue
					}
					prefixes = append(prefixes, prefix)
				}
			}
		}
	}
	for _, s := range appends {
		prefixes = append(prefixes, netip.MustParsePrefix(s))
	}
	return prefixes
}

var defaultNetworkInfo *utils.NetworkInfo

func route(appends []string) error {
	var err error
	defaultNetworkInfo, err = utils.GetDefaultNetworkInfo()
	if err != nil {
		return err
	}
	err = utils.SetInterfaceMetric(defaultNetworkInfo.IfIndex, 1)
	if err != nil {
		return err
	}
	err = utils.UpdateDefaultMetric(defaultNetworkInfo.Gateway, defaultNetworkInfo.IfIndex, 10)
	if err != nil {
		return err
	}
	defaultNetworkInfo, err = utils.GetDefaultNetworkInfo()
	if err != nil {
		return err
	}
	for _, prefix := range routeIps(appends) {
		if prefix.Addr().Is4() {
			_, ipNet, _ := net.ParseCIDR(prefix.String())

			// 记录要添加的路由信息
			log.Printf("Adding route: %s mask %s via %s metric %d if %d",
				prefix.Addr(), net.IP(ipNet.Mask).String(), defaultNetworkInfo.Gateway,
				defaultNetworkInfo.Metric-2, defaultNetworkInfo.IfIndex)

			// 添加路由规则
			err = utils.AddRoute(
				prefix.Addr(),
				netip.MustParseAddr(net.IP(ipNet.Mask).String()),
				netip.MustParseAddr(defaultNetworkInfo.Gateway),
				defaultNetworkInfo.Metric-2,
				defaultNetworkInfo.IfIndex)
			if err != nil {
				log.Println("Failed to add route:", err, prefix)
			}
		}
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Failed to get network interfaces:", err)
		return err
	}
	for _, iface := range interfaces {
		if iface.Name == "utun25" {
			err = utils.AddRoute(
				netip.MustParseAddr("0.0.0.0"),
				netip.MustParseAddr("0.0.0.0"),
				netip.MustParseAddr("172.25.0.0"),
				defaultNetworkInfo.Metric-1,
				iface.Index)
			break
		}
	}

	return err
}
func deleteRoute(appends []string) {
	// 如果 defaultNetworkInfo 为 nil，说明 route() 函数没有成功执行，直接返回
	if defaultNetworkInfo == nil {
		log.Println("deleteRoute: defaultNetworkInfo is nil, skipping route cleanup")
		return
	}

	for _, prefix := range routeIps(appends) {
		if prefix.Addr().Is4() {
			_, ipNet, _ := net.ParseCIDR(prefix.String())
			err := utils.DeleteRoute(
				prefix.Addr(),
				netip.MustParseAddr(net.IP(ipNet.Mask).String()),
				netip.MustParseAddr(defaultNetworkInfo.Gateway),
				defaultNetworkInfo.Metric-2,
				defaultNetworkInfo.IfIndex)
			if err != nil {
				log.Println("Failed to delete route:", err, prefix)
			}
		}
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println("Failed to get network interfaces:", err)
		return
	}
	for _, iface := range interfaces {
		if iface.Name == "utun25" {
			err = utils.DeleteRoute(
				netip.MustParseAddr("0.0.0.0"),
				netip.MustParseAddr("0.0.0.0"),
				netip.MustParseAddr("172.25.0.0"),
				defaultNetworkInfo.Metric-2,
				iface.Index)
			break
		}
	}
	err = utils.SetInterfaceMetric(defaultNetworkInfo.IfIndex, 0)
	if err != nil {
		log.Println("deleteRoute: Failed to reset interface metric:", err)
	}
}
