package rules

import "sync"

type RuleMap struct {
	Map map[string]Rule
}

var instance *RuleMap

// sync.Once type to ensure that the creation of the Singleton instance is thread-safe.
var once sync.Once

func GetInstance() *RuleMap {
	// The once.Do function will guarantee that the initialization code is executed only once.
	once.Do(func() {
		instance = new(RuleMap)
		instance.Map = make(map[string]Rule)
	})
	return instance
}

func (rm *RuleMap) GetRuleByName(ruleName string) Rule {
	for name, r := range rm.Map {
		if name == ruleName {
			return r
		}
	}
	return nil
}
