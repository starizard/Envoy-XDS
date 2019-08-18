package xds

import "time"

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
