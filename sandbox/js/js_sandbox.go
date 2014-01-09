/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2012
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Mike Trinkala (trink@mozilla.com)
#   Rob Miller (rmiller@mozilla.com)
#
# ***** END LICENSE BLOCK *****/
package js

import (
	"fmt"
	"github.com/dancannon/gofetch/sandbox"
	"github.com/robertkrimen/otto"
	_ "github.com/robertkrimen/otto/underscore"
	"log"
)

func getValue(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg == nil {
		result, _ := otto.ToValue(1)
		return result
	}

	name, _ := call.Argument(0).ToString()
	var result otto.Value
	switch name {
	case "PageType":
		result, _ = otto.ToValue(jsb.msg.PageType)
	case "Value":
		result, _ = otto.ToValue(jsb.msg.Value)
	case "Document.Url":
		result, _ = otto.ToValue(jsb.msg.Document.Url)
	case "Document.Title":
		result, _ = otto.ToValue(jsb.msg.Document.Title)
	case "Document.Meta":
		result, _ = otto.ToValue(jsb.msg.Document.Meta)
	case "Document.Doc":
		result, _ = jsb.or.ToValue(jsb.msg.Document.Doc)
	case "Document.Body":
		result, _ = jsb.or.ToValue(jsb.msg.Document.Body)
	default:
		result = otto.UndefinedValue()
	}

	return result
}

func setValue(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg == nil {
		result, _ := otto.ToValue(1)
		return result
	}

	value, _ := call.Argument(0).Export()
	jsb.msg.Value = value

	return otto.UndefinedValue()
}

func setPageType(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg == nil {
		result, _ := otto.ToValue(1)
		return result
	}

	value, _ := call.Argument(0).Export()
	jsb.msg.Value = value

	return otto.UndefinedValue()
}

type JsSandbox struct {
	or     *otto.Otto
	msg    *sandbox.SandboxMessage
	script string
	output func(s string)
	config map[string]interface{}
	err    error
}

func CreateJsSandbox(conf *sandbox.SandboxConfig) (sandbox.Sandbox, error) {
	jsb := new(JsSandbox)
	jsb.or = otto.New()
	jsb.script = conf.Script

	if jsb.or == nil {
		return nil, fmt.Errorf("Sandbox creation failed")
	}

	jsb.output = func(s string) { log.Println(s) }
	jsb.config = conf.Config
	return jsb, nil
}

func (this *JsSandbox) Init() error {
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
	var err error
	_, err = this.or.Run(this.script)
	this.err = err

	return err
}

func (this *JsSandbox) Destroy() error {
	return nil
}

func (this *JsSandbox) Status() int {
	if this.err != nil {
		return int(1)
	} else {
		return int(0)
	}
}

func (this *JsSandbox) Output() string {
	return ""
}

func (this *JsSandbox) LastError() string {
	return this.err.Error()
}

func (this *JsSandbox) ProcessMessage(msg *sandbox.SandboxMessage) int {
	this.msg = msg
	ret, err := this.or.Call("processMessage", nil)
	this.msg = nil

	if err != nil {
		this.err = err
		return 1
	}

	reti, _ := ret.ToInteger()
	return int(reti)
}
