package parser

import (
	"fmt"
)

type Result struct {
	SuccessResponse map[string]interface{}
	FailureResponse *resultFailureResponse
}

type resultFailureResponse struct {
	Code, Description string
}

func (r *Result) FailureReason() error {
	if r.FailureResponse == nil {
		return nil
	}
	return fmt.Errorf("%s: %s", r.FailureResponse.Code, r.FailureResponse.Description)
}

func (r *Result) GetData(indexes ...interface{}) interface{} {
	if r.SuccessResponse == nil {
		return nil
	}
	var temp interface{} = r.SuccessResponse
	var ok bool
	for _, index := range indexes {
		switch t := index.(type) {
		case int:
			if temp, ok = temp.([]interface{}); !ok {
				return nil
			}
			temp = temp.([]interface{})[t]
		case string:
			if temp, ok = temp.(map[string]interface{}); !ok {
				return nil
			}
			temp = temp.(map[string]interface{})[t]
		default:
			return nil
		}
	}
	return temp
}

func NewFailureResult(code, description interface{}) *Result {
	var c, d string
	var ok bool
	if c, ok = code.(string); !ok {
		c = "failure"
	}
	if d, ok = description.(string); !ok {
		d = "unknown"
	}
	return &Result{
		FailureResponse: &resultFailureResponse{
			Code:        c,
			Description: d,
		},
	}
}

func NewSuccessResult(data map[string]interface{}) *Result {
	return &Result{
		SuccessResponse: data,
	}
}
