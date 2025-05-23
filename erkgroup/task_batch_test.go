package erkgroup_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/erkrequire"
	"github.com/orzkratos/synckratos/erkgroup"
	"github.com/orzkratos/synckratos/internal/errors_example"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
)

func TestTaskBatch_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	var taskBatch = erkgroup.NewTaskBatch[uint64, string](args)
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
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"0", "1", "2", "3", "4", "5"}, results)
}

func TestTaskBatch_SetGlide_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	taskBatch := erkgroup.NewTaskBatch[uint64, string](args)
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
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"wa-0", "1", "wa-2", "3", "wa-4", "5"}, results)
}

func TestTaskBatch_EgoRun(t *testing.T) {
	ctx := context.Background()

	var args = []uint64{0, 1, 2, 3, 4, 5}
	taskBatch := erkgroup.NewTaskBatch[uint64, string](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}
	taskBatch.SetGlide(true)

	ego := erkgroup.NewGroup(ctx)
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
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"wa-0", "1", "wa-2", "3", "wa-4", "5"}, results)
}
