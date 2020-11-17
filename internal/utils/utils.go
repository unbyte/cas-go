package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetCallbackURLFromRequest(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	params := r.URL.Query()
	params.Del("ticket")
	return fmt.Sprintf(`%s://%s/%s?%s`, scheme, r.Host, r.URL.Path, params.Encode())
}

func GetRequest(c *http.Client, url string) (body []byte, header http.Header, err error) {
	resp, err := c.Get(url)
	if err != nil {
		return
	}
	result, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return
	}
	return result, resp.Header, nil
}

func GenerateSessionID() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return url.QueryEscape(base64.URLEncoding.EncodeToString(b))
}
