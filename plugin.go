package gofetch

var (
	extractors map[string]Extractor = map[string]Extractor{}
)

func RegisterExtractor(extractor Extractor) {
	extractors[extractor.Id()] = extractor
}

func init() {
	// Register all plugins
	RegisterExtractor(new(OEmbedExtractor))
	RegisterExtractor(new(OpengraphExtractor))
	RegisterExtractor(new(SelectorExtractor))
	RegisterExtractor(new(SelectorTextExtractor))
	RegisterExtractor(new(TextExtractor))
	RegisterExtractor(new(TitleExtractor))
	RegisterExtractor(new(UrlMapperExtractor))
}
