package parser

type Parser func([]byte) (result *Result, success bool)

type ResultHandler func(result *Result) interface{}
