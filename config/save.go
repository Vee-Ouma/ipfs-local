package config

import (
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// SaveInstances save ipfs and ipfs-cluster instances ip and ports on disk
func SaveInstances(instances []ClusterInstance, filename string) {
	str := ""
	// store instances one by one in the format
	/*
		 	Node_0:
			  IP: /ip4/127.0.0.1/tcp/
			  IPFS API Port: 5001
			  IPFS cluster port: 14002
			  Cluster proxy port: 14001
			  Cluster API port: 14000

			Node_1:
			  IP: /ip4/127.0.0.2/tcp/
			  IPFS API Port: 5001
			  IPFS cluster port: 14005
			  Cluster proxy port: 14004
			  Cluster API port: 14003

			...
	*/
	for i, instance := range instances {
		str += nodeIdentifier + strconv.Itoa(i) + ":\n"
		str += ipIdentifier + instance.IP + "\n"
		str += ipfsAPIIdentifier + strconv.Itoa(instance.IPFSAPIPort) + "\n"
		str += clusterPortIdentifier + strconv.Itoa(instance.ClusterPort) + "\n"
		str += clusterProxyPortIdentifier +
			strconv.Itoa(instance.IPFSProxyPort) + "\n"
		str += clusterAPIPortIdentifier +
			strconv.Itoa(instance.RestAPIPort) + "\n\n"
	}
	// write file to system
	ioutil.WriteFile(filepath.Join(statepath, filename), []byte(str), 0)
}

// LoadInstances load saved ipfs and ipfs-cluster instances from disk
func LoadInstances(filename string) []ClusterInstance {
	bytes, err := ioutil.ReadFile(filepath.Join(statepath, filename))
	if err != nil {
		panic(err)
	}
	str := string(bytes)
	lines := strings.Split(str, "\n")

	instances := make([]ClusterInstance, len(lines)/linesPerInstance)
	for i := 0; i < len(instances); i++ {
		instances[i].IP = strings.Split(lines[i*linesPerInstance+1], ipIdentifier)[1]
		instances[i].IPFSAPIPort, err = strconv.Atoi(strings.Split(lines[i*linesPerInstance+2],
			ipfsAPIIdentifier)[1])
		checkErr(err)
		instances[i].ClusterPort, err = strconv.Atoi(strings.Split(lines[i*linesPerInstance+3],
			clusterPortIdentifier)[1])
		checkErr(err)
		instances[i].IPFSProxyPort, err = strconv.Atoi(strings.Split(lines[i*linesPerInstance+4],
			clusterProxyPortIdentifier)[1])
		checkErr(err)
		instances[i].RestAPIPort, err = strconv.Atoi(strings.Split(lines[i*linesPerInstance+5],
			clusterAPIPortIdentifier)[1])
		checkErr(err)
	}
	return instances
}
