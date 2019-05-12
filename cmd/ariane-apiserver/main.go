package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nsyszr/ariane/api/corev1"
	"go.etcd.io/etcd/client"

	"github.com/nats-io/go-nats"
	"github.com/nsyszr/ariane/pkg/api"
)

func handleRequest(data []byte) ([]byte, error) {
	// We expect a request with contains the core/v1/Namespace object
	ns := &corev1.Namespace{}
	req := api.Request{Object: ns}

	// Unmarshal the data to a namespace request, otherwise we have a bad payload
	if err := json.Unmarshal(data, &req); err != nil {
		result, _ := json.Marshal(api.Response{
			StatusCode: 400,
			Error: &api.Error{
				Code:    100000,
				Message: "invalid data",
				Reason:  err.Error(),
			},
		})
		return result, err
	}

	// We're creating the object now
	id, err := uuid.NewUUID()
	if err != nil {
		result, _ := json.Marshal(api.Response{
			StatusCode: 500,
			Error: &api.Error{
				Code:    500000,
				Message: "server error",
				Reason:  err.Error(),
			},
		})
		return result, err
	}
	ns.Metadata.ID = id.String()
	createdAt := time.Now().Round(time.Second).UTC()
	ns.Metadata.CreatedAt = &createdAt
	ns.Status.Phase = "Active"

	exists, err := existsNamespace(ns)
	if err != nil {
		result, _ := json.Marshal(api.Response{
			StatusCode: 500,
			Error: &api.Error{
				Code:    500000,
				Message: "server error",
				Reason:  err.Error(),
			},
		})
		return result, err
	}

	if exists {
		result, _ := json.Marshal(api.Response{
			StatusCode: 409,
			Error: &api.Error{
				Code:    200000,
				Message: "namespace exists already",
			},
		})
		return result, err
	}

	if err := createNamespace(ns); err != nil {
		result, _ := json.Marshal(api.Response{
			StatusCode: 500,
			Error: &api.Error{
				Code:    500000,
				Message: "server error",
				Reason:  err.Error(),
			},
		})
		return result, err
	}

	// Reply with success and object
	result, _ := json.Marshal(api.Response{
		StatusCode: 201,
		Object:     ns,
	})

	return result, nil
}

func createNamespace(ns *corev1.Namespace) error {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: 10 * time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return err
	}
	kapi := client.NewKeysAPI(c)

	// set "/foo" key with "bar" value
	value, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	log.Printf("Setting key '%s' with value: %s", ns.Key(), string(value))
	resp, err := kapi.Set(context.Background(), ns.Key(), string(value), nil)
	if err != nil {
		return err
	}

	log.Printf("Set is done. Metadata is %q\n", resp)
	return nil
}

func existsNamespace(ns *corev1.Namespace) (bool, error) {
	cfg := client.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: 10 * time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return false, nil
	}
	kapi := client.NewKeysAPI(c)

	_, err = kapi.Get(context.Background(), ns.Key(), nil)
	if err != nil {
		if isKeyNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isKeyNotFound(err error) bool {
	if cErr, ok := err.(client.Error); ok {
		return cErr.Code == client.ErrorCodeKeyNotFound
	}
	return false
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Use a WaitGroup to wait for a message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Subscribe
	if _, err := nc.Subscribe("api.core.v1.namespace", func(msg *nats.Msg) {
		log.Printf("received: %s", msg.Data)

		data, err := handleRequest(msg.Data)
		if err != nil {
			log.Print("handle request error: ", err)
		}

		nc.Publish(msg.Reply, data)
		wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	// Wait for a message to come in
	wg.Wait()

	// Close the connection
	nc.Close()
}
