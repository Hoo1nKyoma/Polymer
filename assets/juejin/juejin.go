package juejin

import (
  "io/ioutil"
  "net/http"
  "strings"

  "github.com/tidwall/gjson"

  "Polymer/assets"
)

type juejin struct {
  url string
  query string
  headers map[string]string
}

func init () {
  assets.Register("juejin", newJuejin())
}

func newJuejin () juejin {
  return juejin{
    url: "https://web-api.juejin.im/query",
    query: `{
      "operationName": "",
      "query": "",
      "variables": {
        "tags": [],
        "category": "5562b415e4b00c57d9b94ac8",
        "first": 20,
        "after": "",
        "order": "POPULAR"
      },
      "extensions": {
        "query": {
          "id": "653b587c5c7c8a00ddf67fc66f989d42"
        }
      }
    }`,
    headers: map[string]string{
      "Content-Type": "application/json",
      "X-Agent": "Juejin/Web",
    },
  }
}

func (j juejin) GetTopics (tab string) ([]assets.Topic, error) {
  client := &http.Client{}
  request, err := http.NewRequest("POST", j.url, strings.NewReader(j.query))
  if err != nil {
    return nil, err
  }
  for k, v := range j.headers {
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
  result := gjson.Get(string(body), "data.articleFeed.items.edges")
  topics := []assets.Topic{}
  if result.IsArray() {
    for _, value := range result.Array() {
      title := value.Get("node.title").String()
      link := value.Get("node.originalUrl").String()
      topics = append(topics, assets.Topic{ Title: title, Link: link })
    }
  }
  return topics, nil
}