package web

import (
  "testing"
  "net/url"
  "strings"
  "fmt"
  "time"
  "net"

  "github.com/reusee/goquery"
)

func TestHttpClient(t *testing.T) {
  client := NewClient()
  resp, err := client.PostFormWithHeaders("http://www.xiami.com/member/login", url.Values{
    "done": {"/"},
    "email": {"foo@foo.com"},
    "password": {"bar"},
    "submit": {"登 录"},
    "autologin": {"1"},
  }, &Header{
    "Origin": "http://www.xiami.com",
    "Referer": "http://www.xiami.com/",
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.31 (KHTML, like Gecko) Chrome/26.0.1410.63 Safari/537.31",
  })
  if err != nil {
    t.Fatalf("post form with headers", err)
  }
  doc, err := client.RespToDoc(resp)
  if err != nil {
    t.Fatalf("convert response to goquery doc", err)
  }
  divFound := false
  doc.Find("div").Each(func(i int, s *goquery.Selection) {
    divFound = true
  })
  if !divFound {
    t.Fatalf("div not found")
  }

  doc, err = client.GetDoc("http://www.xiami.com/song/showcollect/id/14916094")
  if err != nil {
    t.Fatalf("get doc", err)
  }
  songFound := false
  doc.Find("div#list_collect a").Each(func(i int, s *goquery.Selection) {
    href, exists := s.Attr("href")
    if !exists || !strings.Contains(href, "/song/") { return }
    fmt.Printf("%s\n", s.Text())
    songFound = true
  })
  if !songFound {
    t.Fatalf("song not found")
  }
}

func TestGetFind(t *testing.T) {
  client := NewClient()
  res, err := client.GetFind("http://www.qq.com", "a")
  if err != nil {
    t.Fatal(err)
  }
  res.Each(func(i int, s *goquery.Selection) {
    href, _ := s.Attr("href")
    fmt.Printf("%s\n", href)
  })
}

func TestTimeout(t *testing.T) {
  client := NewClient()
  SetDialTimeout(time.Second * 1)
  _, err := client.GetDoc("http://www.qq.com:8888")
  if !err.(*url.Error).Err.(*net.OpError).Timeout() {
    t.Fail()
  }
}

func TestEncoding(t *testing.T) {
  client := NewClient()
  client.Encoding = "gbk"
  doc, err := client.GetDoc("http://www.sina.com.cn")
  if err != nil { t.Fail() }
  fmt.Printf("%s\n", doc.Find(".main-nav").Text())
}
