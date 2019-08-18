package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"google.golang.org/grpc"

	xdsUtils "github.com/starizard/envoy-xds/pkg/xds"
)

func main() {
	var configPath string
	var nodeID string
	flag.StringVar(&configPath, "config", "config.json", "path of the config file")
	flag.StringVar(&nodeID, "nodeID", "envoy-1", "Node ID")

	flag.Parse()
	file, _ := ioutil.ReadFile(configPath)
	config := xdsUtils.Config{}
	err := json.Unmarshal([]byte(file), &config)
	if err != nil {
		panic(err)
	}
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
