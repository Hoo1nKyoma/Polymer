package weibo

import (
  "Polymer/assets"
  "github.com/tidwall/gjson"
  "io/ioutil"
  "net/http"
)

type weibo struct {
  url string
  headers map[string]string
}

func init () {
  assets.Register("weibo", newWeibo())
}

func newWeibo () weibo {
  return weibo{
    url: "https://m.weibo.cn/api/container/getIndex?containerid=106003type%3D25%26t%3D3%26disable_hot%3D1%26filter_type%3Drealtimehot&title=%E5%BE%AE%E5%8D%9A%E7%83%AD%E6%90%9C&extparam=pos%3D0_0%26mi_cid%3D100103%26cate%3D10103%26filter_type%3Drealtimehot%26c_type%3D30%26display_time%3D1575949167&luicode=10000011&lfid=231583",
    headers: map[string]string{
      "cache-control": "no-cache",
    },
  }
}

func (w weibo) GetTopics (tab string) ([]assets.Topic, error) {
  client := &http.Client{}
  request, _ := http.NewRequest("GET", w.url, nil)
  for k, v := range w.headers {
    request.Header.Add(k, v)
  }
  response, err := client.Do(request)
  if err != nil {
    return nil, err
  }
  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return nil, err
  }
  result := gjson.Get(string(body), "data.cards.1.card_group")
  topics := []assets.Topic{}
  if result.IsArray() {
    for _, value := range result.Array() {
      title := value.Get("desc").String()
      link := value.Get("scheme").String()
      topics = append(topics, assets.Topic{ Title: title, Link: link })
    }
  }
  return topics, nil
}