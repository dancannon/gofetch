package config

import ()

type RuleProvider interface {
	Setup(map[string]string)
	Provide() []Rule
}

// RethinkDB Provider
