package config

type Rule struct {
	Extractor string   `xml:"extractor,attr"`
	Urls      []string `xml:"url"`
	Values    []Value  `xml:"values>value"`
	Priority  int      `xml:"priority,attr"`
}

type Value struct {
	Name       string      `xml:"name,attr"`
	Value      string      `xml:"value,attr"`
	Parameters []Parameter `xml:"parameter"`
	Children   []Value     `xml:"value"`
}

type ruleProviderConfig struct {
	Id         string      `xml:"id"`
	Parameters []Parameter `xml:"parameter"`
}

type RuleSlice []Rule

func (s RuleSlice) Len() int           { return len(s) }
func (s RuleSlice) Less(i, j int) bool { return s[i].Priority < s[j].Priority }
func (s RuleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
