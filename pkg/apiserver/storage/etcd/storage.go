package etcd

import (
	"path"

	"github.com/nsyszr/ariane/pkg/apiserver/storage"
	"go.etcd.io/etcd/clientv3"
)

type store struct {
	client     *clientv3.Client
	pathPrefix string
}

// New returns an etcd3 implementation of storage.Interface.
func New(c *clientv3.Client, prefix string) storage.Interface {
	return newStore(c, prefix)
}

func newStore(c *clientv3.Client, prefix string) *store {
	result := &store{
		client:     c,
		pathPrefix: path.Join("/", prefix),
	}
	return result
}
