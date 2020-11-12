package parser

import (
	"github.com/clbanning/mxj/v2/x2j-wrapper"
)

func parseXML(content []byte) (Attributes, bool) {
	result := make(map[string]interface{})

	if err := x2j.Unmarshal(content, &result); err != nil {
		return nil, false
	}

	response, ok := result["serviceResponse"].(map[string]interface{})["authenticationSuccess"]
	if !ok {
		return nil, false
	}

	return response.(map[string]interface{}), true
}
