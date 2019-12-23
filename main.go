package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	redis "github.com/go-redis/redis/v7"

	"Polymer/assets"
	_ "Polymer/assets/cnode"
	_ "Polymer/assets/juejin"
	_ "Polymer/assets/v2ex"
	_ "Polymer/assets/weibo"
	_ "Polymer/assets/zhihu"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func main() {
	fmt.Println("the server is listening on port 8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html, err := ioutil.ReadFile("web/index.html")
		if err != nil {
			fmt.Fprint(w, err)
		}
		w.Write(html)
	})

	http.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		asset := r.Form.Get("asset")
		tab := r.Form.Get("tab")
		thirdparty, exist := assets.Assets[asset]
		if !exist {
			fmt.Fprintf(w, "查询的 thirdparty 不存在")
			return
		}
		topics, err := thirdparty.GetTopics(tab)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		topics = insertToRedis(topics)
		json, err := json.Marshal(topics)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		io.WriteString(w, string(json))
	})

	http.ListenAndServe(":8080", nil)
}

func insertToRedis(topics []assets.Topic) []assets.Topic {
	const (
		KEY    = "topics"
		EXPIRE = time.Hour * 8
	)

	var (
		result = []assets.Topic{}
		shouldAddToRedis []string
	)

	v, _ := client.Exists(KEY).Result()
	if v == 0 {
		// 为了防止 topics 没设上缓存，然后存入数据后，将永远不会过期
		if client.SAdd(KEY, nil).Err() != nil || client.Expire(KEY, EXPIRE).Err() != nil {
			// BUG client.SAdd 成功, client.Expire 失败, client.Exists 的结果就是 1 了，导致永远不能设置上 Expire
			client.Del(KEY) // 先暂时这样，但是不能保证 Del 是成功的
			return topics
		}
	}

	for _, topic := range topics {
		exist, _ := client.SIsMember(KEY, topic.Link).Result()
		if exist {
			continue
		} else {
			result = append(result, topic)
			shouldAddToRedis = append(shouldAddToRedis, topic.Link)
		}
	}
	client.SAdd(KEY, shouldAddToRedis)
	return result
}
