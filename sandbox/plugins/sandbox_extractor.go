package plugins

import (
	"github.com/dancannon/gofetch/document"
	. "github.com/dancannon/gofetch/sandbox"
	"github.com/dancannon/gofetch/sandbox/js"

	"errors"
	"fmt"
)

type SandboxExtractor struct {
	sb  Sandbox
	sbc *SandboxConfig
	err error
}

func (s *SandboxExtractor) Init(config interface{}) (err error) {
	if s.sb != nil {
		return // no-op already initialized
	}
	s.sbc = config.(*SandboxConfig)

	switch s.sbc.ScriptType {
	case "js":
		s.sb, err = js.CreateJsSandbox(s.sbc)
		if err != nil {
			return
		}
	default:
		return fmt.Errorf("unsupported script type: %s", s.sbc.ScriptType)
	}

	err = s.sb.Init()
	return
}

func (e *SandboxExtractor) Setup(config interface{}) error {
	return nil
}

func (s *SandboxExtractor) Shutdown() {
	if s.sb != nil {
		s.sb.Destroy()
		s.sb = nil
	}
}

func (s *SandboxExtractor) Extract(doc document.Document) (interface{}, error) {
	if s.sb == nil {
		return nil, s.err
	}

	// Create message
	msg := &SandboxMessage{
		Document: doc,
	}

	retval := s.sb.ProcessMessage(msg)
	if retval > 0 {
		s.err = errors.New("FATAL: " + s.sb.LastError())
	} else if retval < 0 {
		s.err = fmt.Errorf("Failed extracting value")
	}

	return msg.Value, s.err
}

func (s *SandboxExtractor) ExtractValues(doc document.Document) (interface{}, string, error) {
	if s.sb == nil {
		return nil, "", s.err
	}

	// Create message
	msg := &SandboxMessage{
		Document: doc,
	}

	retval := s.sb.ProcessMessage(msg)
	if retval > 0 {
		s.err = errors.New("FATAL: " + s.sb.LastError())
	} else if retval < 0 {
		s.err = fmt.Errorf("Failed extracting value")
	}

	return msg.Value, msg.PageType, s.err
}
