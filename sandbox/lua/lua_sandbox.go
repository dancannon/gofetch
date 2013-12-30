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
package lua

/*
#cgo CFLAGS: -std=gnu99
#cgo CFLAGS: -I/home/daniel/go/src/github.com/dancannon/gofetch/sandbox/lua
#cgo LDFLAGS: -lluasandbox -llua -llpeg -lcjson -lm
#include <stdlib.h>
#include <lua_sandbox.h>
#include "lua_sandbox_interface.h"
*/
import "C"

import (
	"fmt"
	"github.com/dancannon/gofetch/sandbox"
	"log"
	"unsafe"
)

type ValueType int32

const (
	Type_STRING  ValueType = 0
	Type_BYTES   ValueType = 1
	Type_INTEGER ValueType = 2
	Type_DOUBLE  ValueType = 3
	Type_BOOL    ValueType = 4
)

//export go_lua_write_string
func go_lua_write_string(ptr unsafe.Pointer, value *C.char) int {

	var lsb *LuaSandbox = (*LuaSandbox)(ptr)
	if lsb.msg == nil {
		// pipeline.Globals().LogMessage("go_lua_write_string", "No sandbox pack.")
		return 1
	}

	// vC, ok := value.(*C.char)
	// if !ok {
	// 	// return fmt.Errorf("type error, '%s' is a string field", lsb.msg.Name)
	// 	return 1
	// }

	v := C.GoString(value)
	lsb.msg.Value = v
	return 0
}

//export go_lua_write_double
func go_lua_write_double(ptr unsafe.Pointer, value C.double) int {

	var lsb *LuaSandbox = (*LuaSandbox)(ptr)
	if lsb.msg == nil {
		// pipeline.Globals().LogMessage("go_lua_write_string", "No sandbox pack.")
		return 1
	}

	// v, ok := value.(float64)
	// if !ok {
	// 	// return fmt.Errorf("type error, '%s' is a double field", lsb.msg.Name)
	// 	return 1
	// }

	lsb.msg.Value = float64(value)
	return 0
}

//export go_lua_write_bool
func go_lua_write_bool(ptr unsafe.Pointer, value bool) int {

	var lsb *LuaSandbox = (*LuaSandbox)(ptr)
	if lsb.msg == nil {
		// pipeline.Globals().LogMessage("go_lua_write_string", "No sandbox pack.")
		return 1
	}

	// v, ok := value
	// if !ok {
	// 	// return fmt.Errorf("type error, '%s' is an boolean field", lsb.msg.Name)
	// 	return 1
	// }

	lsb.msg.Value = bool(value)
	return 0
}

//export go_lua_read_config
func go_lua_read_config(ptr unsafe.Pointer, c *C.char) (int, unsafe.Pointer, int) {
	name := C.GoString(c)
	var lsb *LuaSandbox = (*LuaSandbox)(ptr)
	if lsb.config == nil {
		return 0, unsafe.Pointer(nil), 0
	}

	v := lsb.config[name]
	switch v.(type) {
	case string:
		s := v.(string)
		cs := C.CString(s) // freed by the caller
		return int(Type_STRING), unsafe.Pointer(cs), len(s)
	case bool:
		b := v.(bool)
		return int(Type_BOOL), unsafe.Pointer(&b), 0
	case int64:
		d := float64(v.(int64))
		return int(Type_INTEGER), unsafe.Pointer(&d), 0
	case float64:
		d := v.(float64)
		return int(Type_DOUBLE), unsafe.Pointer(&d), 0
	}
	return 0, unsafe.Pointer(nil), 0
}

type LuaSandbox struct {
	lsb    *C.lua_sandbox
	msg    *sandbox.Message
	output func(s string)
	config map[string]interface{}
	field  int
}

func CreateLuaSandbox(conf *sandbox.SandboxConfig) (sandbox.Sandbox, error) {
	lsb := new(LuaSandbox)
	cs := C.CString(conf.ScriptFilename)
	defer C.free(unsafe.Pointer(cs))
	md := C.CString(conf.ModuleDirectory)
	defer C.free(unsafe.Pointer(md))
	lsb.lsb = C.lsb_create(unsafe.Pointer(lsb),
		cs,
		md,
		C.uint(conf.MemoryLimit),
		C.uint(conf.InstructionLimit),
		C.uint(conf.OutputLimit))
	if lsb.lsb == nil {
		return nil, fmt.Errorf("Sandbox creation failed")
	}
	lsb.output = func(s string) { log.Println(s) }
	lsb.config = conf.Config
	return lsb, nil
}

func (this *LuaSandbox) Init() error {
	r := int(C.sandbox_init(this.lsb))
	if r != 0 {
		return fmt.Errorf("Init() %s", this.LastError())
	}
	return nil
}

func (this *LuaSandbox) Destroy() error {
	c := C.lsb_destroy(this.lsb, C.CString(""))
	if c != nil {
		err := C.GoString(c)
		C.free(unsafe.Pointer(c))
		return fmt.Errorf("Destroy() %s", err)
	}
	return nil
}
func (this *LuaSandbox) Status() int {
	return int(C.lsb_get_state(this.lsb))
}

func (this *LuaSandbox) Output() string {
	var l C.size_t
	return C.GoString(C.lsb_get_output(this.lsb, &l))
}

func (this *LuaSandbox) LastError() string {
	return C.GoString(C.lsb_get_error(this.lsb))
}

func (this *LuaSandbox) ProcessMessage(msg *sandbox.Message) int {
	this.field = 0
	this.msg = msg
	retval := int(C.process_message(this.lsb))
	this.msg = nil
	return retval
}
