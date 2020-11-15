package parser

import (
	"errors"
	"fmt"
)

type Attributes map[string]interface{}

func (a Attributes) FailureReason() error {
	code, codeExist := a["code"]
	description, descExist := a["description"]
	if !codeExist || !descExist {
		return errors.New("failure: unknown")
	}
	return fmt.Errorf("%s: %s", code, description)
}

func NewFailureAttributes(code, description interface{}) Attributes {
	a := make(Attributes)
	a["code"] = code
	a["description"] = description
	return a
}
