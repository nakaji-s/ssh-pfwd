package main

import (
	"fmt"
)

type Rule struct {
	Id       string
	Priority int
	IsLocal  bool
	Enable   bool
	SSHPortForward
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

func (c *Config) GetRule(id string) (Rule, error) {
	for _, rule := range c.Rules {
		if rule.Id == id {
			return rule, nil
		}
	}
	return Rule{}, fmt.Errorf("err id(%s) not found", id)
}

func (c *Config) UpdateRule(id string, newRule Rule) (Rule, error) {
	for i, _ := range c.Rules {
		if c.Rules[i].Id == id {
			c.Rules[i] = newRule
			return newRule, nil
		}
	}
	return Rule{}, fmt.Errorf("err id(%s) not found", id)
}

func (c *Config) GetRules() []Rule {
	return c.Rules
}
