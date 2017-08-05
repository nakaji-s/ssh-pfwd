package main

import (
	"fmt"

	"github.com/labstack/echo"
)

type Rule struct {
	Id       string
	Priority int
	IsLocal  bool
	Enable   bool
	SSHPortForward
}

type Config interface {
	AddRule(newRule Rule)
	DeleteRule(id string) error
	GetRule(id string) (Rule, error)
	UpdateRule(id string, c echo.Context) (Rule, error)
	GetRules() []Rule
}

type InMemoryConfig struct {
	Rules []Rule
}

func (c *InMemoryConfig) AddRule(newRule Rule) {
	for i, rule := range c.Rules {
		if rule.Priority > newRule.Priority {
			c.Rules = append(c.Rules[:i], append([]Rule{newRule}, c.Rules[i:]...)...)
			return
		}
	}
	c.Rules = append(c.Rules, newRule)
}

func (c *InMemoryConfig) DeleteRule(id string) error {
	for i, rule := range c.Rules {
		if rule.Id == id {
			c.Rules = append(c.Rules[:i], c.Rules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("err id(%s) not found", id)
}

func (c *InMemoryConfig) GetRule(id string) (Rule, error) {
	for _, rule := range c.Rules {
		if rule.Id == id {
			return rule, nil
		}
	}
	return Rule{}, fmt.Errorf("err id(%s) not found", id)
}

func (c *InMemoryConfig) UpdateRule(id string, ctx echo.Context) (Rule, error) {
	for i, _ := range c.Rules {
		if c.Rules[i].Id == id {
			if err := ctx.Bind(&c.Rules[i]); err != nil {
				return Rule{}, err
			}
			return c.Rules[i], nil
		}
	}
	return Rule{}, fmt.Errorf("err id(%s) not found", id)
}

func (c *InMemoryConfig) GetRules() []Rule {
	return c.Rules
}
