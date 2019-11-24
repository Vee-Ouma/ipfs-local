package cluster

import (
	"context"
	"fmt"
	"path/filepath"

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
func AddFile(c client.Client, path string) {
	ctx := context.Background()

	out := make(chan *api.AddedOutput)
	paths := []string{path}
	go c.Add(ctx, paths, api.DefaultAddParams(), out)
	ao := <-out
	fmt.Printf("\nAdded %s: %s\n", filepath.Base(path), ao.Name)
}
