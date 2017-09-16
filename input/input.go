package input

import (
	"encoding/json"
	"net/http"
)

type InputInterface interface {
	Get(key string) string
	Has(key string) bool
	All() map[string]string
}

type input struct {
	values map[string]string
}

func Parse(r *http.Request) *input {

	contentType := r.Header.Get("Content-Type")
	inputs := make(map[string]string)

	if contentType == "application/json" {
		json.NewDecoder(r.Body).Decode(&inputs)
	} else {
		r.ParseForm()
		for fieldName := range r.Form {
			inputs[fieldName] = r.Form.Get(fieldName)
		}
	}

	return &input{inputs}
}

func (i *input) Get(key string) string {
	return i.values[key]
}

func (i *input) Has(key string) bool {
	return i.values[key] != ""
}

func (i *input) All() map[string]string {
	return i.values
}
