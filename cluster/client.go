package cluster

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/ipfs/ipfs-cluster/api"
	"github.com/ipfs/ipfs-cluster/api/rest/client"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// ListPeers of a client
func ListPeers(c client.Client) {
	ctx := context.Background()
	peers, err := c.Peers(ctx)
	checkErr(err)

	fmt.Printf("\nPeers in the Cluster:\n")
	for _, p := range peers {
		fmt.Printf("%s: %s\n", p.Peername, p.Addresses[0])
	}
}

// AddFile to the cluster
func AddFile(c client.Client, path string) string {
	ctx := context.Background()

	cids := make(chan string, 10)
	out := make(chan *api.AddedOutput, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func(ch chan string) {
		defer wg.Done()
		for v := range out {
			if v == nil {
				return
			}
			ch <- v.Cid.String()
		}
	}(cids)

	paths := []string{path}
	start := time.Now()
	c.Add(ctx, paths, api.DefaultAddParams(), out)
	wg.Wait()
	fmt.Println(time.Now().Sub(start))
	name := <-cids
	fmt.Printf("\nAdded %s: %s\n", filepath.Base(path), name)
	return name
}

// CatFile cat an added text file
func CatFile(c client.Client, filename string) {
	ctx := context.Background()

	sh := c.IPFS(ctx)
	start := time.Now()
	rc, err := sh.Cat(filename)
	fmt.Println(time.Now().Sub(start))
	checkErr(err)

	buffer := make([]byte, 1024)
	n, err := rc.Read(buffer)
	checkErr(err)
	str := string(buffer[:n])
	fmt.Printf("\ncat "+filename+":\n%s", str)
}
