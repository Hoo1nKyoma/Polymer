package assets

var Assets = map[string]ThirdParty{}

type ThirdParty interface {
	GetTopics(tab string) ([]Topic, error)
}

type Topic struct {
	Title string
	Link  string
}

func Register(name string, thirdParty ThirdParty) {
	Assets[name] = thirdParty
}

func Contain(array []string, item string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
