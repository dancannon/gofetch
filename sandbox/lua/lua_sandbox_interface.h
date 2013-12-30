/* -*- Mode: C; tab-width: 8; indent-tabs-mode: nil; c-basic-offset: 2 -*- */
/* vim: set ts=2 et sw=2 tw=80: */
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/// Heka Go interfaces for the Lua sandbox @file
#ifndef lua_sandbox_interface_
#define lua_sandbox_interface_

#include "lua.h"
#include "lua_sandbox.h"

// LMW_ERR_*: Lua Message Write errors
extern const int LMW_ERR_NO_SANDBOX_PACK;
extern const int LMW_ERR_WRONG_TYPE;
extern const int LMW_ERR_NEWFIELD_FAILED;
extern const int LMW_ERR_BAD_FIELD_INDEX;
extern const int LMW_ERR_BAD_ARRAY_INDEX;
extern const int LMW_ERR_INVALID_FIELD_NAME;

/**
* Passes a Heka message down to the sandbox for processing. The instruction
* count limits are active during this call.
*
* @param lsb Pointer to the sandbox
*
* @return int Zero on success, non-zero on failure.
*/
int process_message(lua_sandbox* lsb);

/**
* Reads a configuration variable provided in the Heka toml and returns the
* value.
*
* @param lua Pointer to the Lua state.
*
* @return int Returns one value on the stack.
*/
int read_config(lua_State* lua);
int write_value(lua_State* lua);

/**
 * Initializes the sandbox and sets up the above callbacks.
 *
 * @param lsb Pointer to the sandbox.
 * @param data_file File used for the data restoration (empty or NULL for no
 *                  restoration)
 *
 * @return int 0 on success
 */
int sandbox_init(lua_sandbox* lsb);

#endif

