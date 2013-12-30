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
	"io/ioutil"
	"log"
	"os"
)

func readConfig(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
	if jsb.msg == nil {
		result, _ := otto.ToValue(1)
		return result
	}

	name, _ := call.Argument(0).ToString()
	value := jsb.config[name]
	result, _ := otto.ToValue(value)

	return result
}

func writeValue(jsb *JsSandbox, call otto.FunctionCall) otto.Value {
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
	msg    *sandbox.Message
	script string
	output func(s string)
	config map[string]interface{}
	err    error
}

func CreateJsSandbox(conf *sandbox.SandboxConfig) (sandbox.Sandbox, error) {
	jsb := new(JsSandbox)
	jsb.or = otto.New()
	jsb.script = conf.ScriptFilename

	if jsb.or == nil {
		return nil, fmt.Errorf("Sandbox creation failed")
	}

	jsb.output = func(s string) { log.Println(s) }
	jsb.config = conf.Config
	return jsb, nil
}

func (this *JsSandbox) Init() error {
	// Load internal functions
	this.or.Set("readConfig", func(call otto.FunctionCall) otto.Value {
		return readConfig(this, call)
	})
	this.or.Set("writeValue", func(call otto.FunctionCall) otto.Value {
		return writeValue(this, call)
	})

	// Run script
	var script []byte
	var err error

	if this.script == "" || this.script == "-" {
		script, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("Can't read stdin: %v\n", err)
		}
	} else {
		script, err = ioutil.ReadFile(this.script)
		if err != nil {
			return fmt.Errorf("Can't open file \"%v\": %v\n", this.script, err)
		}
	}
	_, err = this.or.Run(string(script))
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

func (this *JsSandbox) ProcessMessage(msg *sandbox.Message) int {
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
