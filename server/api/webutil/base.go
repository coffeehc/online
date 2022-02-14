package webutil

import (
	"fmt"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/gorm"
	"math"
	"net/http"
	"online/server/web/gen/models"
)

type error500Respond struct {
	Payload *models.ActionFailed
}

// WriteResponse to the client
func (o *error500Respond) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(500)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

type ok200Rrespond struct {
	Payload *models.ActionSucceeded
}

// WriteResponse to the client
func (o *ok200Rrespond) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

var (
	pTrue  *bool
	pFalse *bool
)

func init() {
	t := true
	pTrue = &t

	f := false
	pFalse = &f
}

func Int64P(i int64) *int64 {
	return &i
}

func BoolP(b bool) *bool {
	return &b
}

func StrP(i string) *string {
	return &i
}

func Float64P(i float64) *float64 {
	return &i
}
func Float64PWithFixed(i float64) *float64 {
	i = math.Round(i*100) / 100
	return &i
}

func NewActionErrorResponder(from, reason string, items ...interface{}) middleware.Responder {
	var r *string
	if len(items) > 0 {
		r = StrP(fmt.Sprintf(reason, items...))
	} else {
		r = &reason
	}
	return &error500Respond{
		Payload: &models.ActionFailed{
			Ok:     pFalse,
			From:   &from,
			Reason: r,
		},
	}
}

func NewActionSucceedResponder(from string) middleware.Responder {
	return &ok200Rrespond{
		Payload: &models.ActionSucceeded{
			From: &from,
			Ok:   pTrue,
		},
	}
}

func GetInt64ValueOr(raw *int64, value int64) int64 {
	if raw == nil {
		return value
	}
	return *raw
}

func GetStrValueOr(raw *string, value string) string {
	if raw == nil {
		return value
	}
	return *raw
}

func StrEmptyOr(raw *string, value string) string {
	if raw == nil {
		return value
	}

	if *raw == "" {
		return value
	}

	return *raw
}

func GetFloat64ValueOr(raw *float64, value float64) float64 {
	if raw == nil {
		return value
	}
	return *raw
}

func Str(raw *string) string {
	return GetStrValueOr(raw, "")
}

func Int64(raw *int64) int64 {
	return GetInt64ValueOr(raw, 0)
}

func Int(raw *int64) int {
	return int(Int64(raw))
}

func Int64P2Int(raw *int64) int {
	return int(Int64(raw))
}

func Float64(raw *float64) float64 {
	return GetFloat64ValueOr(raw, 0)
}

func IntToInt64P(i int) *int64 {
	return Int64P(int64(i))
}

func Bool(raw *bool) bool {
	if raw == nil {
		return false
	}
	return *raw
}

func ConvertGormModelToSwagger(g gorm.Model) models.GormBaseModel {
	return models.GormBaseModel{
		CreatedAt: Int64P(g.CreatedAt.Unix()),
		ID:        Int64P(int64(g.ID)),
		UpdatedAt: Int64P(g.UpdatedAt.Unix()),
	}
}
