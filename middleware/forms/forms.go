package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// create a custom form field
type Form struct {
	*url.Values
	Errors errors
}

// initialize a form struct
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		getField := f.Get(field)
		if strings.TrimSpace(getField) == "" {
			f.Errors.Add(getField, "This field cannot be blank")
		}
	}
}
func New(data *url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}
func (f *Form) Has(field string, r *http.Request) bool {
	formField := r.Form.Get(field)
	if formField == "" {
		f.Errors.Add(field, "this field should not blank")
		return false
	}
	return true
}
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	x := r.Form.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("%s  should be %d", field, length))
		return false
	}
	return true

}
func (f *Form) IsValidEmail(field string) bool {
	email := f.Get(field)
	if !govalidator.IsEmail(email) {
		f.Errors.Add(field, "Enter valid email")
		return false
	}
	return true

}
