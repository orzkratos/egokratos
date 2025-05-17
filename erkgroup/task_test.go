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
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}

	ctx := context.Background()
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(idx, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run(ctx)
		t.Log(erk)
		erkrequire.NoError(t, erk)
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)
		require.Equal(t, strconv.Itoa(idx), task.Res)
		erkrequire.NoError(t, task.Erk)
	}
}

func TestTaskBatch_SetGlide_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	taskBatch := NewTaskBatch[uint64, string](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}
	taskBatch.SetGlide(true)

	ctx := context.Background()
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(idx, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
			if arg%2 == 0 {
				return "", errors_example.ErrorServerDbError("wrong db")
			}
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run(ctx)
		t.Log(erk)
		erkrequire.NoError(t, erk) //当设置 "平滑继续" 时这里不返回错误
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)
		if idx%2 == 0 {
			require.True(t, errors_example.IsServerDbError(task.Erk))
		} else {
			require.Equal(t, strconv.Itoa(idx), task.Res)
			erkrequire.NoError(t, task.Erk)
		}
	}
}

func TestTaskBatch_EgoRun(t *testing.T) {
	ctx := context.Background()

	var args = []uint64{0, 1, 2, 3, 4, 5}
	taskBatch := NewTaskBatch[uint64, string](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}
	taskBatch.SetGlide(true)

	ego := NewGroup(ctx)
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
		if arg%2 == 0 {
			return "", errors_example.ErrorServerDbError("wrong db")
		}
		res := strconv.FormatUint(arg, 10)
		return res, nil
	})
	erkrequire.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)
		if idx%2 == 0 {
			require.True(t, errors_example.IsServerDbError(task.Erk))
		} else {
			require.Equal(t, strconv.Itoa(idx), task.Res)
			erkrequire.NoError(t, task.Erk)
		}
	}
}
