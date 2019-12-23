package cnode

import (
  "fmt"
  "net/http"

  "github.com/PuerkitoBio/goquery"

  "Polymer/assets"
)

/**
 * "全部": "https://cnodejs.org/?tab=all"
 * "招聘": "https://cnodejs.org/?tab=job"
 */

type cnode struct {
  validityTabs []string
  url string
}

func init () {
  assets.Register("cnode", newCnode())
}

func newCnode () cnode {
  return cnode{
   validityTabs: []string{"all", "job"},
   url: "https://cnodejs.org/?tab=%s",
  }
}

func (c cnode) GetTopics (tab string) ([]assets.Topic, error) {
  exist := assets.Contain(c.validityTabs, tab)
  if !exist {
    return nil, fmt.Errorf("查询的 tab 不存在")
  }
  response, err := http.Get(fmt.Sprintf(c.url, tab))
  if err != nil {
    return nil, err
  }
  defer response.Body.Close()
  doc, err := goquery.NewDocumentFromReader(response.Body)
  if err != nil {
    return nil, err
  }
  topics := []assets.Topic{}
  doc.Find("#main #content .inner #topic_list > .cell .topic_title_wrapper .topic_title").Each(func (i int, s *goquery.Selection) {
    title := s.Text()
    link := fmt.Sprintf("https://cnodejs.org%s", s.AttrOr("href", ""))
    topics = append(topics, assets.Topic{ Title: title, Link: link })
  })
  return topics, nil
}
