/* -*- Mode: C; tab-width: 8; indent-tabs-mode: nil; c-basic-offset: 2 -*- */
/* vim: set ts=2 et sw=2 tw=80: */
/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

/// @brief Sandboxed Lua execution @file
#include <stdlib.h>
#include <stdio.h>
#include <ctype.h>
#include <string.h>
#include <lua.h>
#include <lauxlib.h>
#include <lualib.h>
#include <time.h>
#include <lua_sandbox.h>
#include "_cgo_export.h"

////////////////////////////////////////////////////////////////////////////////
/// Calls to Lua
////////////////////////////////////////////////////////////////////////////////
int process_message(lua_sandbox* lsb)
{
    static const char* func_name = "process_message";
    lua_State* lua = lsb_get_lua(lsb);
    if (!lua) return 1;

    if (lsb_pcall_setup(lsb, func_name)) {
        char err[LSB_ERROR_SIZE];
        snprintf(err, LSB_ERROR_SIZE, "%s() function was not found", func_name);
        lsb_terminate(lsb, err);
        return 1;
    }

    if (lua_pcall(lua, 0, 1, 0) != 0) {
        char err[LSB_ERROR_SIZE];
        size_t len = snprintf(err, LSB_ERROR_SIZE, "%s() %s", func_name,
                              lua_tostring(lua, -1));
        if (len >= LSB_ERROR_SIZE) {
          err[LSB_ERROR_SIZE - 1] = 0;
        }
        lsb_terminate(lsb, err);
        return 1;
    }

    if (!lua_isnumber(lua, 1)) {
        char err[LSB_ERROR_SIZE];
        size_t len = snprintf(err, LSB_ERROR_SIZE,
                              "%s() must return a single numeric value", func_name);
        if (len >= LSB_ERROR_SIZE) {
          err[LSB_ERROR_SIZE - 1] = 0;
        }
        lsb_terminate(lsb, err);
        return 1;
    }

    int status = (int)lua_tointeger(lua, 1);
    lua_pop(lua, 1);

    lsb_pcall_teardown(lsb);

    return status;
}

////////////////////////////////////////////////////////////////////////////////
/// Calls from Lua
////////////////////////////////////////////////////////////////////////////////
int read_config(lua_State* lua)
{
    void* luserdata = lua_touserdata(lua, lua_upvalueindex(1));
    if (NULL == luserdata) {
        luaL_error(lua, "read_config() invalid lightuserdata");
    }
    lua_sandbox* lsb = (lua_sandbox*)luserdata;

    if (lua_gettop(lua) != 1) {
        luaL_error(lua, "read_config() must have a single argument");
    }
    const char* name = luaL_checkstring(lua, 1);

    struct go_lua_read_config_return gr;
    // Cast away constness of the Lua string, the value is not modified
    // and it will save a copy.
    gr = go_lua_read_config(lsb_get_parent(lsb), (char*)name);
    if (gr.r1 == NULL) {
        lua_pushnil(lua);
    } else {
        switch (gr.r0) {
        case 0:
            lua_pushlstring(lua, gr.r1, gr.r2);
            free(gr.r1);
            break;
        case 3:
            lua_pushnumber(lua, *((GoFloat64*)gr.r1));
            break;
        case 4:
            lua_pushboolean(lua, *((GoInt8*)gr.r1));
            break;
        default:
            lua_pushnil(lua);
            break;
        }
    }
    return 1;
}

////////////////////////////////////////////////////////////////////////////////
int write_value(lua_State* lua)
{
    void* luserdata = lua_touserdata(lua, lua_upvalueindex(1));
    if (NULL == luserdata) {
        luaL_error(lua, "write_value() invalid lightuserdata");
    }
    lua_sandbox* lsb = (lua_sandbox*)luserdata;

    int n = lua_gettop(lua);
    if (n != 1) {
        luaL_error(lua, "write_value() incorrect number of arguments");
    }

    int type = lua_type(lua, 1);
    int result;

    switch (type) {
    case LUA_TBOOLEAN: {
        int value = lua_toboolean(lua, 1);
        result = go_lua_write_bool(lsb_get_parent(lsb), value);
        break;
    }
    case LUA_TNUMBER: {
        lua_Number value = lua_tonumber(lua, 1);
            result = go_lua_write_double(lsb_get_parent(lsb), value);
        break;
    }
    case LUA_TSTRING: {
        size_t len;
        const char* value = lua_tostring(lua, 1);
        result = go_lua_write_string(lsb_get_parent(lsb), (char*)value);
        break;
    }
    default:
        luaL_error(lua, "write_message() only accepts numeric, string, or boolean field values");
    }

    if (result != 0) {
        luaL_error(lua, "write_message() failed");
    }
    return 0;
    return 0;
}

////////////////////////////////////////////////////////////////////////////////
int sandbox_init(lua_sandbox* lsb)
{
    if (!lsb) return 1;

    lsb_add_function(lsb, &read_config, "read_config");
    lsb_add_function(lsb, &write_value, "write_value");

    int result = lsb_init(lsb, (const char*) "");
    if (result) return result;

    return 0;
}
