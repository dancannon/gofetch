package gofetch

var (
	extractors map[string]Extractor = map[string]Extractor{}
)

func RegisterExtractor(extractor Extractor) {
	extractors[extractor.Id()] = extractor
}
