package client

import (
	"github.com/nsyszr/ariane/pkg/api/client"
	"github.com/nsyszr/ariane/pkg/api/client/natsio"
)

func NewConfigForNATS(url string) *natsio.Config {
	return &natsio.Config{
		URL: url,
	}
}

func NewClientSetForNATS(cfg *natsio.Config) (client.ClientSet, error) {
	return natsio.NewClientSet(cfg)
}
