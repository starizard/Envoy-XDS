package xds

import (
	"fmt"
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hconnManager "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
)

type Config struct {
	Version      string
	Listeners    []Listener
	Clusters     []Cluster
	RouteConfigs []RouteConfig
}

type Listener struct {
	Name    string
	Port    int
	Address string
}

type Cluster struct {
	Name           string
	ConnectTimeout time.Duration
	SNI            string
	Hosts          []Host
}

type Host struct {
	Name string
	Port int
}

type RouteConfig struct {
	Name    string
	Domains []string
	Routes  []Route
}

type Route struct {
	Regex  string
	Prefix string
	Action RouteAction
}

type RouteAction struct {
	ClusterName   string
	PrefixRewrite string
}

func AddRouteConfig(r *RouteConfig) *v2.RouteConfiguration {
	var routes []*route.Route
	for _, route := range r.Routes {
		routes = append(routes, AddRoute(route))
	}
	return &v2.RouteConfiguration{
		Name: r.Name,
		VirtualHosts: []*route.VirtualHost{&route.VirtualHost{
			Name:    r.Name,
			Domains: r.Domains,
			Routes:  routes,
		}},
	}
}

func AddRoute(r Route) *route.Route {
	return &route.Route{
		Match: &route.RouteMatch{
			PathSpecifier: &route.RouteMatch_Regex{
				Regex: r.Regex,
			},
		},
		Action: &route.Route_Route{
			Route: &route.RouteAction{
				ClusterSpecifier: &route.RouteAction_Cluster{
					Cluster: r.Action.ClusterName,
				},
				PrefixRewrite: r.Action.PrefixRewrite,
			},
		},
	}
}

func AddCluster(c *Cluster) *v2.Cluster {
	var hosts []*core.Address
	for _, host := range c.Hosts {
		hosts = append(hosts, AddHost(host))
	}
	cluster := &v2.Cluster{
		Name:                 c.Name,
		ConnectTimeout:       &c.ConnectTimeout,
		ClusterDiscoveryType: &v2.Cluster_Type{Type: v2.Cluster_LOGICAL_DNS},
		DnsLookupFamily:      v2.Cluster_V4_ONLY,
		LbPolicy:             v2.Cluster_ROUND_ROBIN,
		Hosts:                hosts,
		TlsContext: &auth.UpstreamTlsContext{
			Sni: c.SNI,
		},
	}
	fmt.Printf("\n%v\n", cluster)
	return cluster
}

func AddHost(h Host) *core.Address {
	return &core.Address{Address: &core.Address_SocketAddress{
		SocketAddress: &core.SocketAddress{
			Address:  h.Name,
			Protocol: core.TCP,
			PortSpecifier: &core.SocketAddress_PortValue{
				PortValue: uint32(h.Port),
			},
		},
	}}
}

func AddHTTPConnectionManager(r *RouteConfig) *types.Struct {
	routeConfig := &hconnManager.HttpConnectionManager_RouteConfig{
		RouteConfig: AddRouteConfig(r),
	}
	httpConnManager := &hconnManager.HttpConnectionManager{
		CodecType:      hconnManager.AUTO,
		StatPrefix:     "ingress_http",
		RouteSpecifier: routeConfig,
		HttpFilters: []*hconnManager.HttpFilter{{
			Name: util.Router,
		}},
	}
	configStruct, _ := util.MessageToStruct(httpConnManager)
	return configStruct
}

func AddListener(l *Listener, configStruct *types.Struct) *v2.Listener {
	return &v2.Listener{
		Name: l.Name,
		Address: &core.Address{
			Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Protocol: core.TCP,
					Address:  l.Address,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(l.Port),
					},
				},
			},
		},
		FilterChains: []*listener.FilterChain{{
			Filters: []*listener.Filter{{
				Name: util.HTTPConnectionManager,
				ConfigType: &listener.Filter_Config{
					Config: configStruct,
				},
			}},
		}},
	}
}
