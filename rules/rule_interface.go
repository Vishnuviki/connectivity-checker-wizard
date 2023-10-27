package rules

import "conectivity-checker-wizard/models"

type Rule interface {
	SetNextRule(nextRule Rule)
	SetName(ruleName string)
	Execute(inputData models.InputData) models.ResponseData
}
