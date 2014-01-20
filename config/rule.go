package config

type Rule struct {
	Id          string        `json:"id"`
	Type        string        `json:"type"`
	Host        string        `json:"host"`
	PathPattern string        `json:"path_pattern"`
	UrlPattern  string        `json:"url_pattern"` // This overrides the path and host
	Values      []interface{} `json:"values"`
	Priority    int           `json:"priority,omitempty"`
}

type ProviderConfig struct {
	Id         string            `json:"provider"`
	Parameters map[string]string `json:"params"`
}

type RuleSlice []Rule

func (s RuleSlice) Len() int           { return len(s) }
func (s RuleSlice) Less(i, j int) bool { return s[i].Priority < s[j].Priority }
func (s RuleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
