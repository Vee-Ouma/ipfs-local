package main

import (
	"github.com/guillaumemichel/ipfs-local/cluster"
	"github.com/guillaumemichel/ipfs-local/config"
)

func main() {
	instances := config.LoadInstances("save0")
	cluster.DoThings(instances)
}
