package erkgroup

import "github.com/orzkratos/errkratos"

type Task[A any, R any] struct {
	Arg A
	Res R
	Erk *errkratos.Erk
}

type Tasks[A any, R any] []*Task[A, R]

func (tasks Tasks[A, R]) OkTasks() Tasks[A, R] {
	var okTasks Tasks[A, R]
	for _, task := range tasks {
		if task.Erk == nil {
			okTasks = append(okTasks, task)
		}
	}
	return okTasks
}

func (tasks Tasks[A, R]) WaTasks() Tasks[A, R] {
	var waTasks Tasks[A, R]
	for _, task := range tasks {
		if task.Erk != nil {
			waTasks = append(waTasks, task)
		}
	}
	return waTasks
}

func (tasks Tasks[A, R]) Flatten(newWaFunc func(arg A, erk *errkratos.Erk) R) []R {
	var results = make([]R, 0, len(tasks))
	for _, task := range tasks {
		if task.Erk != nil {
			results = append(results, newWaFunc(task.Arg, task.Erk))
		} else {
			results = append(results, task.Res)
		}
	}
	return results
}
