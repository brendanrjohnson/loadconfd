package backends

import (
	"errors"
	"strings"

	"github.com/brendanrjohnson/loadconfd/backends/consul"
	"github.com/brendanrjohnson/loadconfd/backends/env"
	"github.com/brendanrjohnson/loadconfd/backends/etcd"
	"github.com/brendanrjohnson/loadconfd/backends/zookeeper"
	"github.com/kelseyhightower/confd/log"
)

// The StoreClient interface is implemented by objects that can retrieve
// key/value pairs from a backend store.
type StoreClient interface {
	GetValues(key []string) (map[string]string, error)
	WatchPrefix(prefix string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

//New is used to create a storage client based on our configuration
func New(config Config) (StoreClient, error) {
	if config.Backend == "" {
		config.Backend = "etcd"
	}
	backendNodes := config.BackendNodes
	log.Notice("Backend nodes set to " + strings.Join(backendNodes, ", "))
	switch config.Backend {
	case "consul":
		return consul.NewConsulClient(backendNodes)
	case "etcd":
		// Create the etcd client upfront and use it for the lif of the process.
		// The etcdClient is an http.Client and designed to be reused
		return etcd.NewEtcdClient(backendNodes, config.ClientCert, config.ClientKey, config.ClientCaKeys)
	case "zookeeper":
		return zookeeper.NewZookeeperClient(backendNodes)
	case "env":
		return env.NewEnvClient()
	}
}
