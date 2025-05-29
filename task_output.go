package egokratos

import (
	"github.com/orzkratos/egokratos/internal/utils"
	"github.com/orzkratos/errkratos"
)

type TaskOutput[ARG any, RES any] struct {
	Arg ARG
	Res RES
	Erk *errkratos.Erk
}

func NewOkTaskOutput[ARG any, RES any](arg ARG, res RES) *TaskOutput[ARG, RES] {
	return &TaskOutput[ARG, RES]{
		Arg: arg,
		Res: res,
		Erk: nil,
	}
}

func NewWaTaskOutput[ARG any, RES any](arg ARG, erk *errkratos.Erk) *TaskOutput[ARG, RES] {
	return &TaskOutput[ARG, RES]{
		Arg: arg,
		Res: utils.Zero[RES](),
		Erk: erk,
	}
}

type TaskOutputList[ARG any, RES any] []*TaskOutput[ARG, RES]

func (rs TaskOutputList[ARG, RES]) OkList() TaskOutputList[ARG, RES] {
	var results TaskOutputList[ARG, RES]
	for _, one := range rs {
		if one.Erk == nil {
			results = append(results, one)
		}
	}
	return results
}

func (rs TaskOutputList[ARG, RES]) WaList() TaskOutputList[ARG, RES] {
	var results TaskOutputList[ARG, RES]
	for _, one := range rs {
		if one.Erk != nil {
			results = append(results, one)
		}
	}
	return results
}

func (rs TaskOutputList[ARG, RES]) OkCount() int {
	var cnt int
	for _, one := range rs {
		if one.Erk == nil {
			cnt++
		}
	}
	return cnt
}

func (rs TaskOutputList[ARG, RES]) WaCount() int {
	var cnt int
	for _, one := range rs {
		if one.Erk != nil {
			cnt++
		}
	}
	return cnt
}

func (rs TaskOutputList[ARG, RES]) OkResults() []RES {
	var results []RES
	for _, one := range rs {
		if one.Erk == nil {
			results = append(results, one.Res)
		}
	}
	return results
}

func (rs TaskOutputList[ARG, RES]) WaReasons() []*errkratos.Erk {
	var reasons []*errkratos.Erk
	for _, one := range rs {
		if one.Erk != nil {
			reasons = append(reasons, one.Erk)
		}
	}
	return reasons
}
