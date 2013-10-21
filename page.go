package gofetch

type PageType string

const (
	Unknown   PageType = "unknown"
	PlainText          = "plaintext"
	Article            = "article"
	Audio              = "audio"
	Image              = "image"
	Video              = "video"
	Gallery            = "gallery"
	Flash              = "flash"
)
