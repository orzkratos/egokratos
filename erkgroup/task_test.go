package erkgroup

import (
	"context"
	"strconv"
	"testing"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/erkrequire"
	"github.com/orzkratos/synckratos/internal/errors_example"
	"github.com/stretchr/testify/require"
)

func TestTaskBatch_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	var taskBatch = NewTaskBatch[uint64, string](args)
	ctx := context.Background()
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(ctx, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run()
		t.Log(erk)
		erkrequire.NoError(t, erk)
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)

		//因为数据是 0 1 2 3 4 5 因此测试时恰好可以和下标比较，看看两种方式得到的结果是否相同
		require.Equal(t, idx, int(task.Arg))
		require.Equal(t, strconv.Itoa(idx), task.Res)
		erkrequire.NoError(t, task.Erk)
	}
}

func TestTaskBatch_SetGlide_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	var taskBatch = NewTaskBatch[uint64, string](args)
	taskBatch.SetGlide(true)
	ctx := context.Background()
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(ctx, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
			if arg%2 == 0 {
				return "", errors_example.ErrorServerDbError("wrong db")
			}
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run()
		t.Log(erk)
		erkrequire.NoError(t, erk) //当设置 "平滑继续" 时这里不返回错误
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)

		//因为数据是 0 1 2 3 4 5 因此测试时恰好可以和下标比较，看看两种方式得到的结果是否相同
		require.Equal(t, idx, int(task.Arg))
		if idx%2 == 0 {
			require.True(t, errors_example.IsServerDbError(task.Erk))
		} else {
			require.Equal(t, strconv.Itoa(idx), task.Res)
			erkrequire.NoError(t, task.Erk)
		}
	}
}
