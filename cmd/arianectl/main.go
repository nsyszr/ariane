package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/nsyszr/ariane/api/corev1"
	"github.com/nsyszr/ariane/api/metav1"
	"github.com/nsyszr/ariane/pkg/api"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	ns := corev1.Namespace{
		Kind:       "Namespace",
		APIVersion: "core/v1",
		Metadata: metav1.ObjectMeta{
			Name: "test-5678",
		},
	}

	req := api.Request{
		Method: "CREATE",
		Object: ns,
	}

	data, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}

	// Simulate a JSON Unmarshal error
	// data[0] = 0

	// Send the request
	msg, err := nc.Request("api.core.v1.namespace", data, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Use the response
	log.Printf("reply: %s", msg.Data)

	// Close the connection
	nc.Close()
}
