package simple

import (
	"fmt"
	"strings"
	"sync"

	"github.com/objenious/mproxy/logger"
	"github.com/objenious/mproxy/pkg/session"
)

var _ session.Handler = (*Handler)(nil)

// Handler implements mqtt.Handler interface
type Handler struct {
	logger logger.Logger
	DLMap  sync.Map
}

// New creates new Event entity
func New(logger logger.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// AuthConnect is called on device connection,
// prior forwarding to the MQTT broker
func (h *Handler) AuthConnect(c *session.Client) error {
	h.logger.Info(fmt.Sprintf("AuthConnect() - clientID: %s, username: %s, password: %s, client_CN: %s", c.ID, c.Username, string(c.Password), c.Cert.Subject.CommonName))
	return nil
}

// AuthPublish is called on device publish,
// prior forwarding to the MQTT broker
func (h *Handler) AuthPublish(c *session.Client, topic *string, payload *[]byte) error {
	h.logger.Info(fmt.Sprintf("AuthPublish() - clientID: %s, topic: %s, payload: %s", c.ID, *topic, string(*payload)))
	return nil
}

// AuthSubscribe is called on device publish,
// prior forwarding to the MQTT broker
func (h *Handler) AuthSubscribe(c *session.Client, topics *[]string) error {
	h.logger.Info(fmt.Sprintf("AuthSubscribe() - clientID: %s, topics: %s", c.ID, strings.Join(*topics, ",")))
	for i := range *topics {
		h.DLMap.Store("/obj_dl/1/test/"+(*topics)[i], (*topics)[i])
		(*topics)[i] = "/obj_dl/1/test/"+(*topics)[i]
		h.logger.Debug("sub to " + (*topics)[i])
	}
	return nil
}

// Connect - after client successfully connected
func (h *Handler) Connect(c *session.Client) {
	h.logger.Info(fmt.Sprintf("Connect() - username: %s, clientID: %s", c.Username, c.ID))
}

// Publish - after client successfully published
func (h *Handler) Publish(c *session.Client, topic *string, payload *[]byte) {
	h.logger.Info(fmt.Sprintf("Publish() - username: %s, clientID: %s, topic: %s, payload: %s", c.Username, c.ID, *topic, string(*payload)))
}

// Subscribe - after client successfully subscribed
func (h *Handler) Subscribe(c *session.Client, topics *[]string) {
	h.logger.Info(fmt.Sprintf("Subscribe() - username: %s, clientID: %s, topics: %s", c.Username, c.ID, strings.Join(*topics, ",")))
}

// OnSendToSubscriber - after client successfully subscribed
func (h *Handler) OnSendToSubscriber(c *session.Client, topic *string, payload *[]byte) {
	h.logger.Info(fmt.Sprintf("OnSendToSubscriber() - clientID: %s, topic: %s, payload: %s", c.ID, *topic, string(*payload)))

	v, ok := h.DLMap.Load(*topic)
	if !ok {
		// use the topic as received
		return
	}

	*topic = v.(string)
	h.logger.Info(fmt.Sprintf("end OnSendToSubscriber() - clientID: %s, topic: %s, payload: %s", c.ID, *topic, string(*payload)))

}

// Unsubscribe - after client unsubscribed
func (h *Handler) Unsubscribe(c *session.Client, topics *[]string) {
	h.logger.Info(fmt.Sprintf("Unsubscribe() - username: %s, clientID: %s, topics: %s", c.Username, c.ID, strings.Join(*topics, ",")))
}

// Disconnect on conection lost
func (h *Handler) Disconnect(c *session.Client) {
	h.logger.Info(fmt.Sprintf("Disconnect() - client with username: %s and ID: %s disconenected", c.Username, c.ID))
}
