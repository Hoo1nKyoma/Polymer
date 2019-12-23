package v2ex

import (
  "fmt"
  "net/http"

  "github.com/PuerkitoBio/goquery"

  "Polymer/assets"
)

type v2ex struct {
  validityTabs []string
  url string
}

func init () {
  assets.Register("v2ex", newV2ex())
}

func newV2ex () v2ex {
  return v2ex{
    validityTabs: []string{"tech", "creative", "play", "apple", "jobs", "deals", "city", "qna", "hot"},
    url: "https://www.v2ex.com/?tab=%s",
  }
}

func (v v2ex) GetTopics (tab string) ([]assets.Topic, error) {
  exist := assets.Contain(v.validityTabs, tab)
  if !exist {
    return nil, fmt.Errorf("查询的 tab 不存在")
  }
  response, err := http.Get(fmt.Sprintf(v.url, tab))
  if err != nil {
    return nil, err
  }
  defer response.Body.Close()
  doc, err := goquery.NewDocumentFromReader(response.Body)
  if err != nil {
    return nil, err
  }
  topics := []assets.Topic{}
  doc.Find("#Main .cell.item .item_title .topic-link").Each(func (i int, s *goquery.Selection) {
    title := s.Text()
    link := fmt.Sprintf("https://v2ex.com%s", s.AttrOr("href", ""))
    topics = append(topics, assets.Topic{ Title: title, Link: link })
  })
  return topics, nil
}
