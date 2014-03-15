package js

import (
	"errors"
	"fmt"
	"github.com/dancannon/gofetch/sandbox"
	"github.com/robertkrimen/otto"
	// _ "github.com/robertkrimen/otto/underscore"
	"time"
)

var Halt = errors.New("Halt")

func getValue(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg != nil {
		value, _ := jsb.or.ToValue(jsb.msg.Value)
		return value
	}

	return otto.UndefinedValue()
}
func setValue(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg != nil {
		value, _ := call.Argument(0).Export()
		jsb.msg.Value = value
	}

	return otto.UndefinedValue()
}

func getPageType(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg != nil {
		value, _ := jsb.or.ToValue(jsb.msg.PageType)
		return value
	}

	return otto.UndefinedValue()
}

func setPageType(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg != nil {
		jsb.msg.PageType = call.Argument(0).String()
	}

	return otto.UndefinedValue()
}

type JsSandbox struct {
	or     *otto.Otto
	msg    *sandbox.SandboxMessage
	script string
}

func NewSandbox(conf sandbox.SandboxConfig) sandbox.Sandbox {
	jsb := new(JsSandbox)
	jsb.script = conf.Script
	return jsb
}

func (this *JsSandbox) Init() (err error) {
	this.or = otto.New()

	// Add timeout channel
	this.or.Interrupt = make(chan func())

	// Load internal functions
	this.or.Set("getValue", func(call otto.FunctionCall) otto.Value {
		return getValue(this, call)
	})
	this.or.Set("setValue", func(call otto.FunctionCall) otto.Value {
		return setValue(this, call)
	})
	this.or.Set("getPageType", func(call otto.FunctionCall) otto.Value {
		return getPageType(this, call)
	})
	this.or.Set("setPageType", func(call otto.FunctionCall) otto.Value {
		return setPageType(this, call)
	})

	return err
}

func (this *JsSandbox) Destroy() error {
	this.or.Interrupt = nil
	return nil
}

func (this *JsSandbox) ProcessMessage(msg *sandbox.SandboxMessage) (err error) {
	// Setup panic recovery
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == Halt {
				err = fmt.Errorf("The code took to long! Stopping after: %v\n", duration)
				return
			} else {
				panic(caught)
				err = fmt.Errorf("%v", caught)
				return
			}
		}
	}()

	// Add timeout handler
	this.or.Interrupt = make(chan func())
	go func() {
		time.Sleep(2 * time.Second) // Stop after two seconds
		this.or.Interrupt <- func() {
			panic(Halt)
		}
	}()

	// Setup sandbox with message data
	dv, _ := this.or.ToValue(msg.Document)
	this.or.Set("document", dv)
	this.msg = msg

	// Run script
	_, err = this.or.Run(this.script)

	this.msg = nil
	return err
}
