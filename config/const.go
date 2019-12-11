package config

const (
	//DefaultIPFSAPIPort DefaultIPFSAPIPort
	DefaultIPFSAPIPort = 5001
	// DefaultIPFSGatewayPort DefaultIPFSGatewayPort
	DefaultIPFSGatewayPort = 8080
	// DefaultIPFSSwarmPort DefaultIPFSSwarmPort
	DefaultIPFSSwarmPort = 4001

	// IPVersion default ip version
	IPVersion string = "/ip4/"
	// TransportProtocol default transport protocol
	TransportProtocol string = "/tcp/"

	statepath                  = "../config/saves"
	nodeIdentifier             = "Node_"
	ipIdentifier               = "  IP: "
	ipfsAPIIdentifier          = "  IPFS API Port: "
	clusterPortIdentifier      = "  IPFS cluster port: "
	clusterProxyPortIdentifier = "  Cluster proxy port: "
	clusterAPIPortIdentifier   = "  Cluster API port: "
	linesPerInstance           = 7
)
