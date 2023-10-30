package rulemap

import (
	"sync"

	i "conectivity-checker-wizard/rulemanager/interfaces"
)

type RuleMap struct {
	m map[string]i.Rule
}

var instance *RuleMap

// sync.Once type to ensure that the creation of the Singleton instance is thread-safe.
var once sync.Once

func GetInstance() *RuleMap {
	// The once.Do function will guarantee that the initialization code is executed only once.
	once.Do(func() {
		instance = new(RuleMap)
		instance.m = make(map[string]i.Rule)
	})
	return instance
}

func (rm *RuleMap) GetRuleByName(ruleName string) (i.Rule, bool) {
	v, ok := rm.m[ruleName]
	return v, ok
}

func (rm *RuleMap) AddRule(ruleName string, rule i.Rule) {
	rm.m[ruleName] = rule
}
