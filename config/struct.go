package config

// ClusterInstance details of a cluster
type ClusterInstance struct {
	IP            string
	IPFSAPIPort   int
	RestAPIPort   int
	IPFSProxyPort int
	ClusterPort   int
}
