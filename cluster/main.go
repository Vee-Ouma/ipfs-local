package cluster

import (
	"fmt"

	"github.com/guillaumemichel/ipfs-local/config"
)

func main() {
	instances := config.LoadInstances("save0")
	fmt.Println(instances)
}
