package etcd

import (
	"context"
	"path"

	"github.com/nsyszr/ariane/pkg/api/runtime"
	"github.com/nsyszr/ariane/pkg/api/storage"
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

func (s *store) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {
	// Encode the object to a byte slice
	data, err := runtime.Encode(s.codec, obj)
	if err != nil {
		return err
	}
	key = path.Join(s.pathPrefix, key)

	tx, err := s.client.KV.Txn(ctx).If(
		notFound(key),
	).Then(
		clientv3.OpPut(key, string(data), nil),
	).Commit()

	if err != nil {
		return err
	}
	if !tx.Succeeded {
		return storage.NewKeyExistsError(key)
	}

	/*if out != nil {
		putResp := tx.Responses[0].GetResponsePut()
		return decode(s.codec, s.versioner, data, out, putResp.Header.Revision)
	}*/

	return nil
}

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}
