package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

type Form struct {
	url.Values
	Errors errors
}

func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Has(field string) bool {
	formField := f.Get(field)
	if formField == "" {
		return false
	}
	return true
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if len(strings.TrimSpace(value)) == 0 {
			f.Errors.Add(field, "This field cannot be empty.")
		}
	}
}

func (f *Form) MinLength(field string, length int) bool {
	actualLength := f.Get(field)
	if len(strings.TrimSpace(actualLength)) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must have at least %d characters.", length))
		return false
	}
	return true
}

func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Please enter a valid email address.")
	}
}
