package runtime

type Routes interface {
	CoreV1Namespace() string
}

type natsRoutes struct {
}

func NATSRoutes() Routes {
	return &natsRoutes{}
}

func (r *natsRoutes) CoreV1Namespace() string {
	return "api.core.v1.namespace"
}
