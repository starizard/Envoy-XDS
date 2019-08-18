package xds

import (
	"testing"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/stretchr/testify/assert"
)

func TestAddRouteConfig(t *testing.T) {
	tests := []struct {
		name        string
		routeConfig *RouteConfig
		want        *v2.RouteConfiguration
	}{{
		"",
		&RouteConfig{
			Name:    "route-config-test-1",
			Domains: []string{"*"},
			Routes: []Route{
				{
					Regex: "*",
					Action: RouteAction{
						ClusterName:   "test-cluster",
						PrefixRewrite: "/rewrite",
					},
				},
			},
		},
		&v2.RouteConfiguration{
			Name: "route-config-test-1",
			VirtualHosts: []*route.VirtualHost{&route.VirtualHost{
				Name:    "route-config-test-1",
				Domains: []string{"*"},

				Routes: []*route.Route{{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Regex{
							Regex: "*",
						},
					},
					Action: &route.Route_Route{
						Route: &route.RouteAction{
							ClusterSpecifier: &route.RouteAction_Cluster{
								Cluster: "test-cluster",
							},
							PrefixRewrite: "/rewrite",
						},
					},
				}}}},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddRouteConfig(tt.routeConfig)
			assert.Equal(t, got, tt.want)
		})
	}

}
