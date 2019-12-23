package zhihu

import (
  "fmt"
  "io/ioutil"
  "net/http"

  "github.com/tidwall/gjson"

  "Polymer/assets"
)

type zhihu struct {
  validityTabs []string
  url string
}

func init () {
  assets.Register("zhihu", newZhihu())
}

func newZhihu () zhihu {
  return zhihu{
    validityTabs: []string{"total", "digital"},
    url: "https://www.zhihu.com/api/v3/feed/topstory/hot-lists/%s?limit=50&desktop=true",
  }
}

func (z zhihu) GetTopics (tab string) ([]assets.Topic, error) {
  exist := assets.Contain(z.validityTabs, tab)
  if !exist {
    return nil, fmt.Errorf("查询的 tab 不存在")
  }
  response, err := http.Get(fmt.Sprintf(z.url, tab))
  if err != nil {
    return nil, err
  }
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return nil, err
  }
  result := gjson.Get(string(body), "data")
  topics := []assets.Topic{}
  if result.IsArray() {
    for _, value := range result.Array() {
      title := value.Get("target.title").String()
      url := value.Get("target.url").String()
      topics = append(topics, assets.Topic{ Title: title, Link: url })
    }
  }
  return topics, nil
}
