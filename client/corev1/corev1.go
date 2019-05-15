package corev1

import "github.com/nsyszr/ariane/api/corev1"

type GroupSet interface {
	Namespaces() NamespacesClient
}

type NamespacesClient interface {
	Create(*corev1.Namespace) (*corev1.Namespace, error)
}
