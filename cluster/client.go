package cluster

import (
	"context"
	"fmt"

	"github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/api/rest/client"

	ma "github.com/multiformats/go-multiaddr"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// DoThings on the cluster
func DoThings(ctx context.Context, ip, port string) {
	proxyAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", 14001))
	checkErr(err)

	conf := client.Config{
		Host:      ip,
		Port:      port,
		ProxyAddr: proxyAddr,
	}

	c, err := client.NewDefaultClient(&conf)
	checkErr(err)
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

}
