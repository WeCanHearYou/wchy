package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var (
	//ErrContentTypeNotAllowed is used when POSTing a body that is not json
	ErrContentTypeNotAllowed = errors.New("Only Content-Type application/json is allowed")
	intType                  = reflect.TypeOf(0)
	stringType               = reflect.TypeOf("")
)

//DefaultBinder is the default HTTP binder
type DefaultBinder struct {
}

//NewDefaultBinder creates a new default binder
func NewDefaultBinder() *DefaultBinder {
	return &DefaultBinder{}
}

func methodHasBody(method string) bool {
	return method == http.MethodPost ||
		method == http.MethodDelete
}

//Bind request data to object i
func (b *DefaultBinder) Bind(target interface{}, c *Context) error {
	if methodHasBody(c.Request.Method) && c.Request.ContentLength > 0 {
		contentType := strings.Split(c.Request.Header.Get("Content-Type"), ";")
		if len(contentType) == 0 || contentType[0] != JSONContentType {
			return ErrContentTypeNotAllowed
		}

		if err := json.NewDecoder(c.Request.Body).Decode(target); err != nil {
			return err
		}
	}

	targetValue := reflect.ValueOf(target).Elem()
	targetType := targetValue.Type()
	for i := 0; i < targetValue.NumField(); i++ {
		b.bindRoute(i, targetValue, targetType, c.params)
		b.format(i, targetValue, targetType)
	}
	return nil
}

func (b *DefaultBinder) bindRoute(idx int, target reflect.Value, targetType reflect.Type, params StringMap) error {
	name := targetType.Field(idx).Tag.Get("route")
	if name != "" {
		field := target.Field(idx)
		fieldType := field.Type()
		if fieldType == intType {
			value, _ := strconv.ParseInt(params[name], 10, 64)
			field.SetInt(value)
		} else if fieldType == stringType {
			field.SetString(params[name])
		}
	}

	return nil
}

func (b *DefaultBinder) format(idx int, target reflect.Value, targetType reflect.Type) {
	field := target.Field(idx)

	if field.Type() != stringType {
		return
	}

	format := targetType.Field(idx).Tag.Get("format")
	str := field.Interface().(string)
	str = strings.TrimSpace(str)
	if format == "lower" {
		str = strings.ToLower(str)
	} else if format == "upper" {
		str = strings.ToUpper(str)
	}
	field.SetString(str)
}
