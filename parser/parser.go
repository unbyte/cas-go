package parser

type Parser func([]byte) (Attributes Attributes, success bool)
