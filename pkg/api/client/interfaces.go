package client

import (
	"github.com/nsyszr/ariane/api/corev1"
)

type ClientSet interface {
	CoreV1() CoreV1ClientSet
	Close()
}

type CoreV1ClientSet interface {
	Namespaces() NamespacesClient
}

type NamespacesClient interface {
	Create(*corev1.Namespace) (*corev1.Namespace, error)
}
