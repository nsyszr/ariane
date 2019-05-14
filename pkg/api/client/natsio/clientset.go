package natsio

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/nsyszr/ariane/api/corev1"
	"github.com/nsyszr/ariane/pkg/api"
	"github.com/nsyszr/ariane/pkg/api/client"
	"github.com/nsyszr/ariane/pkg/api/runtime"
)

type Config struct {
	URL string
}

type clientSet struct {
	cfg    *Config
	nc     *nats.Conn
	coreV1 client.CoreV1ClientSet
}

func NewClientSet(cfg *Config) (client.ClientSet, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, err
	}

	return &clientSet{
		cfg:    cfg,
		nc:     nc,
		coreV1: newCoreV1ClientSet(cfg, nc),
	}, nil
}

func (cs *clientSet) Close() {
	cs.nc.Close()
}

func (cs *clientSet) CoreV1() client.CoreV1ClientSet {
	return cs.coreV1
}

type coreV1ClientSet struct {
	cfg        *Config
	nc         *nats.Conn
	namespaces client.NamespacesClient
}

func newCoreV1ClientSet(cfg *Config, nc *nats.Conn) client.CoreV1ClientSet {
	return &coreV1ClientSet{
		cfg:        cfg,
		nc:         nc,
		namespaces: newNamespacesClient(cfg, nc),
	}
}

func (cs *coreV1ClientSet) Namespaces() client.NamespacesClient {
	return cs.namespaces
}

type namespacesClient struct {
	cfg *Config
	nc  *nats.Conn
}

func newNamespacesClient(cfg *Config, nc *nats.Conn) client.NamespacesClient {
	return &namespacesClient{
		cfg: cfg,
		nc:  nc,
	}
}

func (c *namespacesClient) Create(obj *corev1.Namespace) (*corev1.Namespace, error) {
	req := api.Request{
		Method: "CREATE",
		Object: obj,
	}

	data, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	// Send the request
	msg, err := c.nc.Request(runtime.NATSRoutes().CoreV1Namespace(), data, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// We expect a response with a namespace object
	out := &corev1.Namespace{}
	res := api.Response{Object: out}

	if err := json.Unmarshal(msg.Data, &res); err != nil {
		return nil, err
	}

	// Check if the object was created successful, otherwise return an error
	if res.StatusCode != 201 {
		return nil, res.Error
	}

	return out, nil
}
