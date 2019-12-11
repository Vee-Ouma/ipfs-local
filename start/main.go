package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/guillaumemichel/ipfs-local/config"
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
	nodes := flag.Int("nodes", DefaultIPFSInstances,
		"number of ipfs instances")
	clusters := flag.Int("clusters", DefaultIPFSInstances,
		"number of ipfs-cluster instances")

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

	instances := make([]config.ClusterInstance, clusterN)

	secret, p, err := config.SetupClusterLeader(leaderPath,
		ClusterInstanceName+"0", addrs[0], DefaultReplMin, DefaultReplMax)
	bootstrap := p.IP + strconv.Itoa(p.ClusterPort)
	if err != nil {
		fmt.Println("Error leader:", err)
	}
	fmt.Printf("\n***** IPFS cluster leader started at %s *****\n\n", bootstrap)
	instances[0] = *p

	wg := sync.WaitGroup{}

	currCluster := 1
	for i := 0; i < ipfsN; i++ {
		//count := currCluster
		clusterC := clusterN / ipfsN // +count
		if i < clusterN%ipfsN {
			clusterC++
		}
		if i == 0 {
			clusterC--
		}
		wg.Add(1)
		go func(n, cN int) {
			slavePath := path.Join(DefaultConfigPath, addrs[n])
			for j := cN; j < cN+clusterC; j++ {
				p, err := config.SetupClusterSlave(slavePath,
					ClusterInstanceName+strconv.Itoa(j), addrs[n],
					bootstrap, secret, DefaultReplMin, DefaultReplMax)
				if err != nil {
					fmt.Println("Error slave:", err)
				}
				instances[j] = *p
			}
			wg.Done()
		}(i, currCluster)
		currCluster += clusterC
	}
	wg.Wait()
	time.Sleep(14 * time.Second)
	fmt.Println("\n***** IPFS cluster instances started *****")

	config.SaveInstances(instances, "save0")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println()
	clean()

	os.Exit(0)
}
