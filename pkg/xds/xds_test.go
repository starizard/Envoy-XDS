package xds

import (
	"testing"
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
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

func TestAddCluster(t *testing.T) {
	timeout := 2 * time.Second
	tests := []struct {
		name    string
		cluster *Cluster
		want    *v2.Cluster
	}{{
		"",
		&Cluster{
			Name:           "app-backend",
			ConnectTimeout: 2 * time.Second,
			SNI:            "www.google.com",
			Hosts: []Host{{
				Name: "www.google.com",
				Port: 443,
			},
			},
		},
		&v2.Cluster{
			Name:                 "app-backend",
			ConnectTimeout:       &timeout,
			ClusterDiscoveryType: &v2.Cluster_Type{Type: v2.Cluster_LOGICAL_DNS},
			DnsLookupFamily:      v2.Cluster_V4_ONLY,
			LbPolicy:             v2.Cluster_ROUND_ROBIN,
			Hosts: []*core.Address{&core.Address{Address: &core.Address_SocketAddress{
				SocketAddress: &core.SocketAddress{
					Address:  "www.google.com",
					Protocol: core.TCP,
					PortSpecifier: &core.SocketAddress_PortValue{
						PortValue: uint32(443),
					},
				},
			}}},
			TlsContext: &auth.UpstreamTlsContext{
				Sni: "www.google.com",
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AddCluster(tt.cluster)
			assert.Equal(t, got, tt.want)
		})
	}

}
