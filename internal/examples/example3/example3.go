package example3

import (
	"github.com/orzkratos/egokratos"
)

type Step1Param struct {
	NumA int
}

type Step1Result struct {
	ResA         string
	Step2Outputs egokratos.TaskOutputList[*Step2Param, *Step2Result]
}

type Step2Param struct {
	NumB int
}

type Step2Result struct {
	ResB         string
	Step3Outputs egokratos.TaskOutputList[*Step3Param, *Step3Result]
}

type Step3Param struct {
	NumC int
}

type Step3Result struct {
	ResC string
}

func NewStep1Params(paramCount int) []*Step1Param {
	var params = make([]*Step1Param, 0, paramCount)
	for idx := 0; idx < paramCount; idx++ {
		params = append(params, &Step1Param{NumA: idx})
	}
	return params
}

func NewStep2Params(paramCount int) []*Step2Param {
	var params = make([]*Step2Param, 0, paramCount)
	for idx := 0; idx < paramCount; idx++ {
		params = append(params, &Step2Param{NumB: idx})
	}
	return params
}

func NewStep3Params(paramCount int) []*Step3Param {
	var params = make([]*Step3Param, 0, paramCount)
	for idx := 0; idx < paramCount; idx++ {
		params = append(params, &Step3Param{NumC: idx})
	}
	return params
}
