package sandbox

import (
	. "github.com/dancannon/gofetch/message"
)

const (
	STAT_LIMIT   = 0
	STAT_CURRENT = 1
	STAT_MAXIMUM = 2

	TYPE_MEMORY       = 0
	TYPE_INSTRUCTIONS = 1
	TYPE_OUTPUT       = 2
)

type Sandbox interface {
	// Sandbox control
	Init() error
	Destroy() error

	// Sandbox state
	Status() int
	Output() string
	LastError() string

	// Plugin functions
	ProcessMessage(msg *ExtractMessage) int
}

type SandboxConfig struct {
	Id               string `json:id"`
	ScriptType       string `json:"script_type"`
	ScriptFilename   string `json:"filename"`
	ModuleDirectory  string `json:"module_directory"`
	PreserveData     bool   `json:"preserve_data"`
	MemoryLimit      uint   `json:"memory_limit"`
	InstructionLimit uint   `json:"instruction_limit"`
	OutputLimit      uint   `json:"output_limit"`
	Config           map[string]interface{}
}
