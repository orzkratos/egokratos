package erkgroup

import (
	"context"
	"sync/atomic"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/synckratos/internal/utils"
	"github.com/yyle88/must/mustnum"
)

type Task[A any, R any] struct {
	Arg A
	Res R
	Erk *errkratos.Erk
}

type TaskBatch[A any, R any] struct {
	Tasks []*Task[A, R]
	Index int64
	Glide bool // Glide 标志位，控制是否平滑继续，有的时候只要有一个子任务失败就算失败(set false)，而有时候它们是独立的(set true)
}

func NewTaskBatch[A any, R any](args []A) *TaskBatch[A, R] {
	tasks := make([]*Task[A, R], 0, len(args))
	for idx := 0; idx < len(args); idx++ {
		tasks = append(tasks, &Task[A, R]{
			Arg: args[idx],
			Res: utils.Zero[R](),
			Erk: nil,
		})
	}
	return &TaskBatch[A, R]{
		Tasks: tasks,
		Index: 0,
		Glide: false,
	}
}

func (t *TaskBatch[A, R]) GetRun(ctx context.Context, run func(ctx context.Context, arg A) (R, *errkratos.Erk)) func() *errkratos.Erk {
	newValue := atomic.AddInt64(&t.Index, 1)
	sliceIdx := int(newValue) - 1
	mustnum.Less(sliceIdx, len(t.Tasks)) //这里限制不要超过下标，这就需要外部控制调用次数，认为这是基本的，不应该用错
	task := t.Tasks[sliceIdx]
	return func() *errkratos.Erk {
		res, erk := run(ctx, task.Arg) //这里面你也不要panic，假如有panic需要调用者自己恢复
		if erk != nil {
			task.Erk = erk
			if t.Glide {
				return nil //当出错时，假如是设置“平滑继续”标志，就不返回错误，这样外层的 ctx 就不会被 cancel 掉，这也符合设计的目的
			}
			return erk
		}
		task.Res = res
		return nil
	}
}

func (t *TaskBatch[A, R]) SetGlide(glide bool) {
	t.Glide = glide
}
