package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/guillaumemichel/local_clusters/config"
)

func clean() {
	exec.Command("killall", "-9", "ipfs").Run()
	exec.Command("killall", "-9", "ipfs-cluster-se").Run()
}

// StartIPFS start an ipfs instance
func startIPFS(ip string) {
	path := path.Join(DefaultConfigPath, ip, IPFSConfFolder)
	err := config.CreateEmptyDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// init ipfs in the desired folder
	exec.Command("ipfs", "-c"+path, "init").Run()

	// edit the ip in the config file
	config.EditIPFSConfig(path, ip)

	// start the ipfs daemon
	// we need to fork the process
	go func() {
		exec.Command("ipfs", "-c"+path, "daemon").Run()
		//fmt.Println("ipfs at ip", ip, "crashed")
	}()
}

// CreateNode start a node, containing an ipfs instance and cluster instances
func createNode(ip string) {
	path := path.Join(DefaultConfigPath, ip)

	// create the empty directory that will store the node configs
	err := config.CreateEmptyDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	//path := *flag.String("path", DefaultConfigPath, "path for the config files")
	nodes := flag.Int("nodes", DefaultIPFSInstances, "number of ipfs instances")
	clusters := flag.Int("clusters", DefaultIPFSInstances, "number of ipfs-cluster instances")

	flag.Parse()
	// help message
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	ipfsN := *nodes
	clusterN := *clusters

	if ipfsN < 2 {
		fmt.Println("we want more than 1 single node ...")
		os.Exit(1)
	}

	if clusterN < ipfsN {
		clusterN = ipfsN
	}

	clean()

	err := config.CreateEmptyDir(DefaultConfigPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	addrs := make([]string, ipfsN)
	for i := 0; i < ipfsN; i++ {
		addrs[i] = LocalAddr + strconv.Itoa(i+1)
	}

	fmt.Println("\n***** Starting IPFS nodes *****")
	fmt.Println(addrs)

	// start ipfs instances
	for i := 0; i < ipfsN; i++ {
		go func(i int) {
			createNode(addrs[i])
			startIPFS(addrs[i])
		}(i)
	}
	time.Sleep(IPFSStartupTime)

	fmt.Println("***** IPFS instances started successfully *****")
	fmt.Println()

	// setup the leader of the cluster
	leaderPath := path.Join(DefaultConfigPath, addrs[0])

	secret, bootstrap, err := config.SetupClusterLeader(leaderPath,
		ClusterInstanceName+"0", addrs[0], DefaultReplMin, DefaultReplMax)
	if err != nil {
		fmt.Println("Error leader:", err)
	}
	fmt.Println("***** IPFS cluster leader started at", bootstrap, "*****")

	currCluster := 0
	for i := 0; i < ipfsN; i++ {
		count := currCluster
		clusterC := clusterN/ipfsN + count
		if clusterC%ipfsN <= i {
			clusterC++
		}
		if i == 0 { // don't count master twice
			currCluster++
		}
		slavePath := path.Join(DefaultConfigPath, addrs[i])
		for j := currCluster; j < clusterC; j++ {
			go func(i, currCluster int) {
				err = config.SetupClusterSlave(slavePath,
					ClusterInstanceName+strconv.Itoa(currCluster), addrs[i],
					bootstrap, secret, DefaultReplMin, DefaultReplMax)
				if err != nil {
					fmt.Println("Error slave:", err)
				}
			}(i, currCluster)
			currCluster++
		}
	}
	time.Sleep(ClusterStartupTime)
	fmt.Println("***** IPFS cluster instances started *****")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println()
	clean()

	os.Exit(0)
}
