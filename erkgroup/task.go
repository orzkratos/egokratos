package erkgroup

import (
	"context"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/synckratos/internal/utils"
	"github.com/yyle88/must"
	"github.com/yyle88/must/mustnum"
)

type Task[A any, R any] struct {
	Arg A
	Res R
	Erk *errkratos.Erk
}

type TaskBatch[A any, R any] struct {
	Tasks     []*Task[A, R]
	Glide     bool // Glide 标志位，控制是否平滑继续，有的时候只要有一个子任务失败就算失败(set false)，而有时候它们是独立的(set true)
	newCtxErk func(erx error) *errkratos.Erk
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
		Glide: false,
	}
}

// GetRun 通过传递 idx 依次获取 errgroup.Go 需要的运行函数，依次使用 errgroup.Go 调用即可得到结果
func (t *TaskBatch[A, R]) GetRun(idx int, run func(ctx context.Context, arg A) (R, *errkratos.Erk)) func(ctx context.Context) *errkratos.Erk {
	mustnum.Less(idx, len(t.Tasks)) //这里限制不要超过下标，这就需要外部控制调用次数，认为这是基本的，不应该用错
	task := t.Tasks[idx]
	return func(ctx context.Context) *errkratos.Erk {
		if t.newCtxErk != nil && ctx.Err() != nil {
			erk := t.newCtxErk(ctx.Err())
			must.Full(erk) //这里避免被外部诓骗，你的错误函数不能返回假的
			task.Erk = erk
			if t.Glide {
				return nil //这个标志位是"平滑继续"的作用。即使 ctx 出错时后续也无法执行，但依然需要他们都走到这里，把错误设置到结果里
			}
			return erk
		}
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

// EgoRun 既演示如何使用 GetRun，而且当任务逻辑较重而调用层逻辑较轻时，也可以直接反过来以调用层为参数，相当于传过来个调度器，调度执行函数
func (t *TaskBatch[A, R]) EgoRun(ego *Group, run func(ctx context.Context, arg A) (R, *errkratos.Erk)) {
	for idx := 0; idx < len(t.Tasks); idx++ {
		ego.Go(t.GetRun(idx, run))
	}
}

func (t *TaskBatch[A, R]) SetGlide(glide bool) {
	t.Glide = glide
}

func (t *TaskBatch[A, R]) SetNewCtxErk(newCtxErk func(erx error) *errkratos.Erk) {
	t.newCtxErk = newCtxErk
}
