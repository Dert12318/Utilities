package context

import (
	"github.com/Dert12318/Utilities/common/constant/header"
)

type (
	Source interface {
		Apply(c *Context)
	}
	httpSource       struct{}
	MandatoryRequest struct {
		requestID string `validate:"required"`
		token     string
	}
	messagingSource struct {
		header map[string]string
	}
)

func (m messagingSource) Apply(c *Context) {
	requestID := m.header[header.MessagingRequestID]
	token := m.header[header.MessagingAuthorization]

	c.mandatory = MandatoryRequest{
		requestID: requestID,
		token:     token,
	}
}

func (m MandatoryRequest) RequestID() string {
	return m.requestID
}

func (m MandatoryRequest) Token() string {
	return m.token
}

func (c *Context) SetMandatory(source Source) {
	source.Apply(c)
}

func (m MandatoryRequest) String() string {
	var stringBuilder = "Mandatory : {RequestID:" + m.requestID

	return stringBuilder + "}"
}

func (m MandatoryRequest) Source() map[string]interface{} {
	var response = make(map[string]interface{})

	if m.RequestID() != "" {
		response["requestId"] = m.RequestID()
	}
	return response
}
