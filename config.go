package main

type Rule struct {
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
