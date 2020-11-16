package main

import (
	"fmt"
	"github.com/unbyte/cas-go"
	"github.com/unbyte/cas-go/api"
	"github.com/unbyte/cas-go/parser"
)

func main() {
	data, err := cas.New(cas.Option{
		APIInstance:    api.NewAPIv2("https://pass.neu.edu.cn/tpass/", "http://202.118.31.197/"),
		Store:          cas.DefaultStore(0),
		SessionManager: cas.DefaultSessionManager(""),
		ResultHandler: func(result *parser.Result) interface{} {
			return result.GetData("attributes", 0, "ID_NUMBER")
		},
	}).ValidateTicket("ST-1513787-tBFjv1XJJjYF0Lk2eECc-tpass")
	fmt.Println(err)
	fmt.Printf("%+v\n", data)
}
