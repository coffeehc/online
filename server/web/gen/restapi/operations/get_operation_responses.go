// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"online/server/web/gen/models"
)

// GetOperationOKCode is the HTTP code returned for type GetOperationOK
const GetOperationOKCode int = 200

/*GetOperationOK 查询 Operation 记录

swagger:response getOperationOK
*/
type GetOperationOK struct {

	/*
	  In: Body
	*/
	Payload *models.OperationsResponse `json:"body,omitempty"`
}

// NewGetOperationOK creates GetOperationOK with default headers values
func NewGetOperationOK() *GetOperationOK {

	return &GetOperationOK{}
}

// WithPayload adds the payload to the get operation o k response
func (o *GetOperationOK) WithPayload(payload *models.OperationsResponse) *GetOperationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get operation o k response
func (o *GetOperationOK) SetPayload(payload *models.OperationsResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetOperationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
