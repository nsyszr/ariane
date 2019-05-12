package corev1

import (
	"fmt"
	"strings"

	"github.com/nsyszr/ariane/api/metav1"
)

type Namespace struct {
	Kind       string            `json:"kind"`
	APIVersion string            `json:"apiVersion"`
	Metadata   metav1.ObjectMeta `json:"metadata"`
	Spec       NamespaceSpec     `json:"spec,omitempty"`
	Status     NamespaceStatus   `json:"status,omitempty"`
	// namespaced bool
}

func (obj *Namespace) Key() string {
	return fmt.Sprintf("/api/%s/%s/%s",
		strings.ToLower(obj.APIVersion),
		strings.ToLower(obj.Kind),
		strings.ToLower(obj.Metadata.Name))
}

type NamespaceSpec struct {
}

type NamespaceStatus struct {
	Phase string `json:"phase,omitempty"`
}
