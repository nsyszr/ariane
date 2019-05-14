package main

import (
	"log"

	"github.com/nsyszr/ariane/api/corev1"
	"github.com/nsyszr/ariane/api/metav1"
	"github.com/nsyszr/ariane/client"
)

func main() {
	cs, err := client.NewClientSetForNATS(client.NewConfigForNATS("localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer cs.Close()

	obj := &corev1.Namespace{
		Kind:       "Namespace",
		APIVersion: "core/v1",
		Metadata: metav1.ObjectMeta{
			Name: "test-1212",
		},
	}

	out, err := cs.CoreV1().Namespaces().Create(obj)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v", out)
}

/*func main() {
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
}*/
