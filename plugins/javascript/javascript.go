package oembed

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/plugins"
	. "github.com/dancannon/gofetch/sandbox"
	"github.com/dancannon/gofetch/sandbox/js"

	"fmt"
)

type JavaScriptExtractor struct {
	sb  Sandbox
	sbc *SandboxConfig
}

func (e *JavaScriptExtractor) Setup(config interface{}) error {
	var params map[string]interface{}
	if p, ok := config.(map[string]interface{}); !ok {
		params = make(map[string]interface{})
	} else {
		params = p
	}

	e.sb = nil
	e.sbc = new(SandboxConfig)
	if script, ok := params["script"]; !ok {
		return fmt.Errorf("The javascript sandbox extractor must be passed an script")
	} else {
		e.sbc.Script = script.(string)
		e.sbc.ScriptType = "js"
	}

	switch e.sbc.ScriptType {
	case "js":
		var err error
		e.sb, err = js.CreateJsSandbox(e.sbc)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported script type: %s", e.sbc.ScriptType)
	}

	return e.sb.Init()
}

func (e *JavaScriptExtractor) Extract(doc document.Document) (interface{}, error) {
	if e.sb == nil {
		return nil, fmt.Errorf("Sandbox not setup")
	}

	// Create message
	msg := &SandboxMessage{
		Document: doc,
	}

	retval := e.sb.ProcessMessage(msg)
	if retval > 0 {
		return nil, fmt.Errorf("FATAL: %s", e.sb.LastError())
	} else if retval < 0 {
		return nil, fmt.Errorf("Failed extracting value")
	}

	return msg.Value, nil
}

func (e *JavaScriptExtractor) ExtractValues(doc document.Document) (interface{}, string, error) {
	if e.sb == nil {
		return nil, "", fmt.Errorf("Sandbox not setup")
	}

	// Create message
	msg := &SandboxMessage{
		Document: doc,
	}

	retval := e.sb.ProcessMessage(msg)
	if retval > 0 {
		return nil, "", fmt.Errorf("FATAL: %s", e.sb.LastError())
	} else if retval < 0 {
		return nil, "", fmt.Errorf("Failed extracting value")
	}

	return msg.Value, msg.PageType, nil
}

func init() {
	RegisterPlugin("javascript", new(JavaScriptExtractor))
}
