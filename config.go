package main

import (
	"fmt"
)

type Rule struct {
	Id         string
	Priority   int
	SshAddr    string
	LocalAddr  string
	RemoteAddr string
	IsLocal    bool
}

type Config struct {
	Rules []Rule
}

func (c *Config) AddRule(newRule Rule) {
	for i, rule := range c.Rules {
		if rule.Priority > newRule.Priority {
			c.Rules = append(c.Rules[:i], append([]Rule{newRule}, c.Rules[i:]...)...)
			return
		}
	}
	c.Rules = append(c.Rules, newRule)
}

func (c *Config) DeleteRule(id string) error {
	for i, rule := range c.Rules {
		if rule.Id == id {
			c.Rules = append(c.Rules[:i], c.Rules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("err id(%s) not found", id)
}
