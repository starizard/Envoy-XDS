package xds

import (
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
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
	Host           string
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
