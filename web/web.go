package web

import (
  "net/http"
  "net/http/cookiejar"
  "net/url"
  "strings"
  "net"
  "time"
  "io"
  "bytes"

  "github.com/reusee/goquery"
  "code.google.com/p/go.net/html"
  "github.com/reusee/go-packages/rune_conv"
)

type Header map[string]string

type Client struct {
  *http.Client
  Encoding string
}

var dialTimeout = time.Second * 10
var headerTimeout = time.Second * 10

func SetDialTimeout(t time.Duration) {
  dialTimeout = t
}

func SetHeaderTimeout(t time.Duration) {
  headerTimeout = t
}

func dialWithTimeout(network, addr string) (net.Conn, error) {
  return net.DialTimeout(network, addr, dialTimeout)
}

var transport = &http.Transport{
  Dial: dialWithTimeout,
  Proxy: http.ProxyFromEnvironment,
  ResponseHeaderTimeout: headerTimeout,
}

func NewClient() *Client {
  jar, _ := cookiejar.New(nil)
  client := &http.Client{
    Transport: transport,
  }
  client.Jar = jar
  return &Client{
    client,
    "utf-8",
  }
}

func (self *Client) PostFormWithHeaders(url string, values url.Values, headers *Header) (resp *http.Response, err error) {
  req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
  if err != nil {
    return nil, err
  }
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  if headers != nil {
    for key, value := range *headers {
      req.Header.Add(key, value)
    }
  }
  return self.Do(req)
}

func (self *Client) RespToDoc(resp *http.Response) (*goquery.Document, error) {
  var reader io.Reader = resp.Body
  var err error
  if self.Encoding != "utf-8" {
    buf := new(bytes.Buffer)
    io.Copy(buf, resp.Body)
    runes, err := rune_conv.From(self.Encoding, buf.Bytes())
    if err != nil {
      return nil, err
    }
    reader = bytes.NewReader([]byte(string(runes)))
  }
  node, err := html.Parse(reader)
  if err != nil {
    print("here\n")
    return nil, err
  }
  resp.Body.Close()
  return goquery.NewDocumentFromNode(node), nil
}

func (self *Client) GetDoc(url string) (*goquery.Document, error) {
  resp, err := self.Get(url)
  if err != nil {
    return nil, err
  }
  return self.RespToDoc(resp)
}

func (self *Client) GetFind(url string, selector string) (*goquery.Selection, error) {
  doc, err := self.GetDoc(url)
  if err != nil {
    return nil, err
  }
  return doc.Find(selector), nil
}
