package plugins

import (
	"github.com/dancannon/gofetch/message"
	. "github.com/dancannon/gofetch/sandbox"
	"github.com/dancannon/gofetch/sandbox/js"
	"github.com/dancannon/gofetch/sandbox/lua"

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
	case "lua":
		s.sb, err = lua.CreateLuaSandbox(s.sbc)
		if err != nil {
			return
		}
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

func (s *SandboxExtractor) Extract(msg *message.ExtractMessage) (err error) {
	if s.sb == nil {
		err = s.err
		return
	}

	retval := s.sb.ProcessMessage(msg)
	if retval > 0 {
		s.err = errors.New("FATAL: " + s.sb.LastError())
	} else if retval < 0 {
		s.err = fmt.Errorf("Failed extracting value")
	}
	err = s.err
	return
}
