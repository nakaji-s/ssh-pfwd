package main

import (
	"errors"

	"github.com/jinzhu/gorm"
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
	AddRule(newRule Rule) error
	DeleteRule(id string) error
	GetRule(id string) (Rule, error)
	UpdateRule(id string, c echo.Context) (Rule, error)
	GetRules() ([]Rule, error)
}

// InMemory
type InMemoryConfig struct {
	Rules []Rule
}

var ErrIdNotFound = errors.New("id not found")

func (c *InMemoryConfig) AddRule(newRule Rule) error {
	for i, rule := range c.Rules {
		if rule.Priority > newRule.Priority {
			c.Rules = append(c.Rules[:i], append([]Rule{newRule}, c.Rules[i:]...)...)
			return nil
		}
	}
	c.Rules = append(c.Rules, newRule)
	return nil
}

func (c *InMemoryConfig) DeleteRule(id string) error {
	for i, rule := range c.Rules {
		if rule.Id == id {
			c.Rules = append(c.Rules[:i], c.Rules[i+1:]...)
			return nil
		}
	}
	return ErrIdNotFound
}

func (c *InMemoryConfig) GetRule(id string) (Rule, error) {
	for _, rule := range c.Rules {
		if rule.Id == id {
			return rule, nil
		}
	}
	return Rule{}, ErrIdNotFound
}

func (c *InMemoryConfig) UpdateRule(id string, ctx echo.Context) (Rule, error) {
	for i, _ := range c.Rules {
		if c.Rules[i].Id == id {
			if err := ctx.Bind(&c.Rules[i]); err != nil {
				return Rule{}, err
			}
			// Reconnect
			c.Rules[i].SSHPortForward.Stop()
			if c.Rules[i].Enable == true {
				c.Rules[i].SSHPortForward.Start()
			}
			return c.Rules[i], nil
		}
	}
	return Rule{}, ErrIdNotFound
}

func (c *InMemoryConfig) GetRules() ([]Rule, error) {
	return c.Rules, nil
}

// sqlite
type SqliteConfig struct {
	Db *gorm.DB
}

func (c *SqliteConfig) AddRule(newRule Rule) error {
	return c.Db.Create(newRule).Error
}

func (c *SqliteConfig) DeleteRule(id string) error {
	return c.Db.Where("id = ?", id).Delete(&Rule{}).Error
}

func (c *SqliteConfig) GetRule(id string) (Rule, error) {
	var rule Rule
	err := c.Db.Where("id = ?", id).First(&rule).Error

	return rule, err
}

func (c *SqliteConfig) UpdateRule(id string, ctx echo.Context) (Rule, error) {
	var rule Rule
	err := c.Db.Where("id = ?", id).First(&rule).Error
	if err != nil {
		return Rule{}, err
	}

	if err = ctx.Bind(&rule); err != nil {
		return Rule{}, err
	}
	// Reconnect
	rule.SSHPortForward.Stop()
	if rule.Enable == true {
		rule.SSHPortForward.Start()
	}

	err = c.Db.Model(&rule).Where("id = ?", id).Update(&rule).Error

	return rule, err
}

func (c *SqliteConfig) GetRules() ([]Rule, error) {
	var rules []Rule
	err := c.Db.Find(&rules).Error
	return rules, err
}
