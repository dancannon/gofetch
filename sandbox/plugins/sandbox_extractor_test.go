package plugins

import (
	"fmt"
	. "github.com/dancannon/gofetch/sandbox"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

var value *interface{}
var result map[string]interface{}

func getTestMessage() Message {
	ptr := new(interface{})

	result = make(map[string]interface{})
	result["test"] = ptr

	return Message{
		Type:   "extractor",
		Name:   "test",
		Value:  ptr,
		Result: result,
		Config: map[string]interface{}{},
	}
}

func TestLuaSandbox(t *testing.T) {
	var err error

	// Setup message
	msg := getTestMessage()

	// Setup sandbox
	var sbc SandboxConfig
	sbc.ScriptType = "lua"
	sbc.ScriptFilename = "../lua/test_scripts/test.lua"

	// Setup extractor
	var extractor SandboxExtractor
	err = extractor.Init(&sbc)
	err = extractor.Extract(&msg)

	spew.Dump(msg.Result, err)
}

func TestJsSandbox(t *testing.T) {
	var err error

	// Setup message
	msg := getTestMessage()

	// Setup sandbox
	var sbc SandboxConfig
	sbc.ScriptType = "js"
	sbc.ScriptFilename = "../js/test_scripts/test.js"

	// Setup extractor
	var extractor SandboxExtractor
	err = extractor.Init(&sbc)
	spew.Dump(err)
	err = extractor.Extract(&msg)

	fmt.Print("\n")
	spew.Dump(msg.Result, err)
}
