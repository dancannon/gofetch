package js

import (
	"errors"
	"fmt"
	"github.com/dancannon/gofetch/sandbox"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"time"
)

var Halt = errors.New("Halt")

func getValue(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg == nil {
		result, _ := jsb.or.ToValue(1)
		return result
	}

	name, _ := call.Argument(0).ToString()
	var result otto.Value
	switch name {
	case "PageType":
		result, _ = jsb.or.ToValue(jsb.msg.PageType)
	case "Value":
		result, _ = jsb.or.ToValue(jsb.msg.Value)
	case "Document.URL":
		result, _ = jsb.or.ToValue(jsb.msg.Document.URL.String())
	case "Document.Title":
		result, _ = jsb.or.ToValue(jsb.msg.Document.Title)
	case "Document.Meta":
		result, _ = jsb.or.ToValue(jsb.msg.Document.Meta)
	case "Document.Doc":
		result, _ = jsb.or.ToValue(*jsb.msg.Document.Doc)
	case "Document.Body":
		result, _ = jsb.or.ToValue(*jsb.msg.Document.Body)
	default:
		result = otto.UndefinedValue()
	}

	return result
}

func setValue(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg == nil {
		result, _ := jsb.or.ToValue(1)
		return result
	}

	value, _ := call.Argument(0).Export()
	jsb.msg.Value = value

	return otto.UndefinedValue()
}

func setPageType(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg == nil {
		result, _ := jsb.or.ToValue(1)
		return result
	}

	jsb.msg.PageType = call.Argument(0).String()

	return otto.UndefinedValue()
}

type JsSandbox struct {
	or     *otto.Otto
	msg    *sandbox.SandboxMessage
	script string
}

func NewSandbox(conf sandbox.SandboxConfig) (sandbox.Sandbox, error) {
	jsb := new(JsSandbox)
	jsb.or = otto.New()
	jsb.script = conf.Script

	if jsb.or == nil {
		return nil, fmt.Errorf("Sandbox creation failed")
	}

	return jsb, nil
}

func (this *JsSandbox) Init() (err error) {
	// Setup panic recovery
	start := time.Now()
	defer func() {
		duration := time.Since(start)
		if caught := recover(); caught != nil {
			if caught == Halt {
				err = fmt.Errorf("The code took to long! Stopping after: %v\n", duration)
				return
			} else {
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

	// Load internal functions
	this.or.Set("getValue", func(call otto.FunctionCall) otto.Value {
		return getValue(this, call)
	})
	this.or.Set("setValue", func(call otto.FunctionCall) otto.Value {
		return setValue(this, call)
	})
	this.or.Set("setPageType", func(call otto.FunctionCall) otto.Value {
		return setPageType(this, call)
	})

	// Run script
	_, err = this.or.Run(this.script)

	this.or.Interrupt = nil

	return err
}

func (this *JsSandbox) Destroy() error {
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

	this.msg = msg

	_, err = this.or.Call("processMessage", nil)

	this.msg = nil
	this.or.Interrupt = nil

	return err
}
