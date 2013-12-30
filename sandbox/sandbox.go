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
	ProcessMessage(msg *Message) int
}

type SandboxConfig struct {
	ScriptType       string `toml:script_type"`
	ScriptFilename   string `toml:"filename"`
	ModuleDirectory  string `toml:"module_directory"`
	PreserveData     bool   `toml:"preserve_data"`
	MemoryLimit      uint   `toml:"memory_limit"`
	InstructionLimit uint   `toml:"instruction_limit"`
	OutputLimit      uint   `toml:"output_limit"`
	Config           map[string]interface{}
}
