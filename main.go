package main

import (
	"fmt"
	"net"
	"time"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"google.golang.org/grpc"

	xdsUtils "github.com/starizard/envoy-xds/pkg/xds"
)

var nodeID = "envoy-1"

// TODO: get this from file
var config xdsUtils.Config = xdsUtils.Config{
	Version: "1.0.1",
	Listeners: []xdsUtils.Listener{{
		Name:    "listener-1",
		Port:    8081,
		Address: "127.0.0.1",
	}},
	Clusters: []xdsUtils.Cluster{
		{
			Name:           "app-backend",
			ConnectTimeout: 2 * time.Second,
			SNI:            "example.com",
			Hosts: []xdsUtils.Host{{
				Name: "example.com",
				Port: 443,
			},
			},
		},
	},
	RouteConfig: xdsUtils.RouteConfig{
		Name:    "route-config-test-1",
		Domains: []string{"*"},
		Routes: []xdsUtils.Route{
			{
				Regex: "/*",
				Action: xdsUtils.RouteAction{
					ClusterName:   "app-backend",
					PrefixRewrite: "/rewrite",
				},
			},
		},
	},
}

func main() {

	snapshotCache := cache.NewSnapshotCache(false, xdsUtils.Hasher{}, nil)
	server := xds.NewServer(snapshotCache, nil)
	grpcServer := grpc.NewServer()
	listeners, clusters := config.Make()
	fmt.Printf("\n%v\n%v\n", listeners, clusters)
	snap := cache.NewSnapshot(config.Version, nil, clusters, nil, listeners)
	snapshotCache.SetSnapshot(nodeID, snap)

	lis, _ := net.Listen("tcp", ":18000")
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	api.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	api.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	api.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	api.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	func() {
		fmt.Printf("\nXDS server started\n")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()
}
