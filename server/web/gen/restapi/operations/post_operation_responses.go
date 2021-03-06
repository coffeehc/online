// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"online/server/web/gen/models"
)

// PostOperationOKCode is the HTTP code returned for type PostOperationOK
const PostOperationOKCode int = 200

/*PostOperationOK API 调用成功

swagger:response postOperationOK
*/
type PostOperationOK struct {

	/*
	  In: Body
	*/
	Payload *models.ActionSucceeded `json:"body,omitempty"`
}

// NewPostOperationOK creates PostOperationOK with default headers values
func NewPostOperationOK() *PostOperationOK {

	return &PostOperationOK{}
}

// WithPayload adds the payload to the post operation o k response
func (o *PostOperationOK) WithPayload(payload *models.ActionSucceeded) *PostOperationOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post operation o k response
func (o *PostOperationOK) SetPayload(payload *models.ActionSucceeded) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostOperationOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
