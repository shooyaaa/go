package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//http mine type definition
const (
	MimeTypeEncode = "application/x-www-form-urlencoded"
	MimeTypeJson   = "application/json"
	MimeTypeHtml   = "text/html"
	MimeTypeAVIF   = "image/avif"
	MimeTypeWebp   = "image/webp"
	MimeTypeApng   = "image/png"
	MimeTypeSvg    = "image/svg+xml"
	MimeTypeForm   = "multipart/form-data"
)

//http headers
const (
	HeaderContentType = "Content-Type"
	HeaderRefer       = "Referer"
	HeaderAgent       = "User-Agent"
	HeaderAccept      = "Accept"
	HeaderAcceptLang  = "Accept-Language"
)

const (
	AgentStringChrome = "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"
)

type HttpReqMethod string

const HttpReqPost HttpReqMethod = "POST"
const HttpReqGet HttpReqMethod = "GET"
const HttpReqUpdate HttpReqMethod = "UPDATE"
const HttpReqDelete HttpReqMethod = "DELETE"

type Http struct {
}

func (http *Http) ParamsToForm(params map[string]interface{}) string {
	tokens := make([]string, 0)
	for key, raw := range params {
		value := ""
		switch raw.(type) {
		case int:
			value = strconv.Itoa(raw.(int))
		case float64:
			value = strconv.FormatFloat(raw.(float64), 'f', -1, 64)
		case string:
			value = raw.(string)
		}
		tokens = append(tokens, fmt.Sprintf("%v=%v", key, value))
	}
	return strings.Join(tokens, "&")
}

func (http *Http) HttpGet(url string, params map[string]interface{}, headers map[string]string) (interface{}, error) {
	return http.HttpRawRequest(HttpReqGet, params, url, headers)
}

func (http *Http) HttpPostForm(url string, params map[string]interface{}) (interface{}, error) {
	headers := make(map[string]string)
	headers[HeaderContentType] = MimeTypeForm
	return http.HttpRawRequest(HttpReqPost, params, url, headers)
}

func (http *Http) HttpPostJson(url string, params map[string]interface{}) (interface{}, error) {
	headers := make(map[string]string)
	headers[HeaderContentType] = MimeTypeJson
	return http.HttpRawRequest(HttpReqPost, params, url, headers)
}

func (hp *Http) HttpRawRequest(method HttpReqMethod, params map[string]interface{}, url string, headers map[string]string) (interface{}, error) {
	var reader io.Reader
	var write *multipart.Writer
	if method == HttpReqGet && params != nil {
		url += "?" + hp.ParamsToForm(params)
	} else if method == HttpReqPost {
		var buffer []byte
		if ct, ok := headers[HeaderContentType]; ok && ct == MimeTypeJson {
			buffer, _ = json.Marshal(params)
		} else if ct == MimeTypeForm {
			var body bytes.Buffer
			write = multipart.NewWriter(&body)
			for key, value := range params {
				sValue, ok := value.(string)
				if ok {
					part, _ := write.CreateFormFile(key, key)
					if sValue[0:1] == "@" {
						file, err := os.Open(sValue[1:])
						if err != nil {
							fmt.Println(err)
						}
						contents, _ := ioutil.ReadAll(file)
						part.Write(contents)
					} else if sValue[0:1] == "#" {
						part.Write([]byte(sValue[1:]))
					}

				} else {
					write.WriteField(key, fmt.Sprintf("%v", value))
				}
			}
			write.Close()
			buffer = body.Bytes()
		}
		if buffer != nil {
			reader = bytes.NewReader(buffer)
		}
	}
	req, err := http.NewRequest(string(method), url, reader)
	if err != nil {
		return req, err
	}
	for header, value := range headers {
		req.Header.Set(header, value)
	}
	req.Header.Set(HeaderAgent, AgentStringChrome)
	if write != nil {
		req.Header.Set(HeaderContentType, write.FormDataContentType())
	}
	client := &http.Client{Timeout: time.Second * 3}
	return hp.HandleResponse(client.Do(req))
}

func (http *Http) HandleResponse(resp *http.Response, err error) (interface{}, error) {
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()

	var ret interface{}
	ct := resp.Header.Get(HeaderContentType)
	switch {
	case strings.Contains(ct, MimeTypeJson):
		ret = make(map[string]interface{})
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal(body, &ret)
		return ret, err
	case strings.Contains(ct, MimeTypeHtml):
		ret, err = goquery.NewDocumentFromReader(resp.Body)
	case strings.Contains(ct, MimeTypeApng):
		bts, _ := io.ReadAll(resp.Body)
		ret = string(bts)
	}
	return ret, err
}
