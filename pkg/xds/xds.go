package xds

type Config struct {
}

type Listener struct {
	Name    string
	Port    int
	Address string
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
