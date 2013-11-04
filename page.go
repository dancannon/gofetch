package gofetch

type PageType string

const (
	Unknown   PageType = "unknown"
	PlainText          = "plaintext"
	Text               = "text"
	Audio              = "audio"
	Image              = "image"
	Video              = "video"
	Gallery            = "gallery"
	Flash              = "flash"
)
