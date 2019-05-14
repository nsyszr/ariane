package natsio

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/nsyszr/ariane/api/corev1"
	"github.com/nsyszr/ariane/pkg/api/handlers"
	"github.com/nsyszr/ariane/pkg/api/storage"

	"github.com/nats-io/go-nats"
)

type handler struct {
	nc         *nats.Conn
	pathPrefix string
	store      storage.Interface
}

type header map[string][]string

type request struct {
	Method string      `json:"method"`
	Header header      `json:"header,omitempty"`
	Body   interface{} `json:"body"`
}

type response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body,omitempty"`
	// Error      *errorObject `json:"error,omitempty"`
}

type errorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func (h *handler) Register(nc *nats.Conn) {
	subj := h.pathPrefix + ".*"
	nc.Subscribe(subj, h.dispatchRequest)
}

func (h *handler) dispatchRequest(msg *nats.Msg) {
	if strings.HasPrefix(msg.Subject, h.pathPrefix+".core.v1.namespaces") {
		h.dispatchCoreV1NamespacesRequest(msg)
	}
}

func (h *handler) dispatchCoreV1NamespacesRequest(msg *nats.Msg) {
	obj := &corev1.Namespace{}
	req := request{Body: obj}

	// Unmarshal the data to a namespace request, otherwise we have a bad payload
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		if err := h.replyError(msg.Reply, handlers.NewInvalidObjectError(err.Error())); err != nil {
			log.Print("failed to reply error: ", err)
		}
		return
	}

	// Has name in path?
	if strings.HasPrefix(msg.Subject, h.pathPrefix+".core.v1.namespaces.") {
		obj.Metadata.Name = strings.TrimPrefix(msg.Subject, h.pathPrefix+".core.v1.namespaces.")
	}

	if err := h.createObject(obj); err != nil {
		if err := h.replyError(msg.Reply, err); err != nil {
			log.Print("failed to reply error: ", err)
		}
	}
}

func (h *handler) createObject(obj interface{}) error {
	return nil
}

func (h *handler) replyError(subj string, e error) (err error) {
	var res []byte

	if he, ok := err.(*handlers.HandlerError); ok {
		res, err = json.Marshal(response{
			StatusCode: he.StatusCode(),
			Body:       he,
		})
	} else {
		res, err = json.Marshal(response{
			StatusCode: 500,
			Body: &errorObject{
				Code:    0,
				Message: e.Error(),
			},
		})
	}
	if err != nil {
		return
	}

	return h.nc.Publish(subj, res)
}
