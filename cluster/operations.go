package cluster

import (
	"fmt"
	"strconv"

	"github.com/guillaumemichel/ipfs-local/config"

	"github.com/ipfs/ipfs-cluster/api/rest/client"
	ma "github.com/multiformats/go-multiaddr"
)

// DoThings on the cluster
func DoThings(instances []config.ClusterInstance) {
	proxyAddr, err := ma.NewMultiaddr(fmt.Sprintf(instances[0].IP +
		strconv.Itoa(instances[0].IPFSProxyPort)))
	checkErr(err)
	apiAddr, err := ma.NewMultiaddr(fmt.Sprintf(instances[0].IP +
		strconv.Itoa(instances[0].RestAPIPort)))
	checkErr(err)

	conf := client.Config{
		APIAddr:   apiAddr,
		ProxyAddr: proxyAddr,
	}

	c, err := client.NewDefaultClient(&conf)
	checkErr(err)

	ListPeers(c)
	cid := AddFile(c, "../data/hello.txt")
	AddFile(c, "../data/rfc1918.txt")
	AddFile(c, "../data/kiddo.gif")

	CatFile(c, cid)

	/*
		id, err := c.ID(ctx)
		checkErr(err)
		fmt.Println("ID:", id)
		peers, err := c.Peers(ctx)
		checkErr(err)

		fmt.Println()
		for i, p := range peers {
			fmt.Printf("Peer %d: %s, %s\n", i, p.Peername, p.ID)
		}
		out := make(chan *api.AddedOutput)
		paths := []string{"data/hello.txt"}
		go c.Add(ctx, paths, api.DefaultAddParams(), out)
		checkErr(err)
		fmt.Println("Added hello.txt, waiting for CID")
		ao := <-out
		fmt.Println(ao)
		sh := c.IPFS(ctx)
		fmt.Println("Got shell", sh)
		rc, err := sh.Cat(ao.Name)
		checkErr(err)

		fmt.Println("Still alive", rc)
		buffer := make([]byte, 1024)
		n, err := rc.Read(buffer)
		checkErr(err)
		fmt.Println("I read the shit")
		str := string(buffer[:n])
		fmt.Println(str)
	*/
}

func Test0(instances []config.ClusterInstance) {
	proxyAddr, err := ma.NewMultiaddr(fmt.Sprintf(instances[0].IP +
		strconv.Itoa(instances[0].IPFSProxyPort)))
	checkErr(err)
	apiAddr, err := ma.NewMultiaddr(fmt.Sprintf(instances[0].IP +
		strconv.Itoa(instances[0].RestAPIPort)))
	checkErr(err)

	conf := client.Config{
		APIAddr:   apiAddr,
		ProxyAddr: proxyAddr,
	}

	_, err = client.NewDefaultClient(&conf)
	checkErr(err)

}
