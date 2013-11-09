package gofetch

type Rule struct {
	Extractor  string            `xml:"extractor,attr"`
	Urls       []string          `xml:"url"`
	Parameters []configParameter `xml:"parameter"`
	Priority   int               `xml:"priority"`
}

type RuleSlice []Rule

func (s RuleSlice) Len() int           { return len(s) }
func (s RuleSlice) Less(i, j int) bool { return s[i].Priority < s[j].Priority }
func (s RuleSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
