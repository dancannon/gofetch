package sandbox

import (
	"github.com/dancannon/gofetch/document"
)

const (
	STAT_LIMIT   = 0
	STAT_CURRENT = 1
	STAT_MAXIMUM = 2

	TYPE_MEMORY       = 0
	TYPE_INSTRUCTIONS = 1
	TYPE_OUTPUT       = 2
)

type SandboxMessage struct {
	PageType string
	Value    interface{}

	Document document.Document
}

type Sandbox interface {
	// Sandbox control
	Init() error
	Destroy() error

	// Plugin functions
	ProcessMessage(msg *SandboxMessage) error
}

type SandboxConfig struct {
	Script     string `json:"script"`
	ScriptType string `json:"script_type"`
}
