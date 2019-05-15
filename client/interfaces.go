package client

import (
	corev1client "github.com/nsyszr/ariane/client/corev1"
)

type ClientSet interface {
	CoreV1() corev1client.GroupSet
	Close()
}
