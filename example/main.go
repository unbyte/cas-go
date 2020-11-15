package main

import (
	"fmt"
	"github.com/unbyte/cas-go"
	"github.com/unbyte/cas-go/api"
)

func main() {
	a, err := cas.New(cas.Option{
		APIInstance:     api.NewAPIv2("https://pass.neu.edu.cn/tpass/", "http://202.118.31.197/"),
		AttributesStore: cas.DefaultAttributesStore(),
		SessionManager:  cas.DefaultSessionManager("nb"),
	}).ValidateTicket("ST-1493992-jnckLnXTqdd1cVuKyWWv-tpass")
	fmt.Println(err)
	fmt.Printf("%+v\n", a)
}
