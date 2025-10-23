package egokratos_test

import (
	"context"
	"math/rand/v2"
	"strconv"
	"testing"
	"time"

	"github.com/orzkratos/egokratos"
	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/egokratos/internal/errorspb"
	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/must/erkrequire"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func TestGroup_Go_TaskRun(t *testing.T) {
	ego := erkgroup.NewGroup(context.Background())
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	var taskBatch = egokratos.NewTaskBatch[uint64, string](args)
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(idx, taskRun))
	}
	erkrequire.Error(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erx)
	}
}

func taskRun(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
	if ctx.Err() != nil {
		zaplog.LOG.Info("task no", zap.Uint64("arg", arg))
		return "", errorspb.ErrorWrongContext("error=%v", ctx.Err())
	}
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // Simulate computation time // 模拟计算延迟
	if arg%10 == 3 {
		zaplog.LOG.Info("task wa", zap.Uint64("arg", arg))
		return "", errorspb.ErrorServerDbError("task wa %d", arg) // Simulate task execution failure // 模拟某个任务失败
	}
	zaplog.LOG.Info("task ok", zap.Uint64("arg", arg))

	res := strconv.FormatUint(arg, 10)
	return res, nil
}

// TestGroup_Go_SetGlide_TaskRun tests glide mode with concurrent task execution
// Validates that glide mode allows independent task execution without stopping on errors
//
// TestGroup_Go_SetGlide_TaskRun 测试平滑模式与并发任务执行
// 验证平滑模式允许独立任务执行，不会因错误而停止
func TestGroup_Go_SetGlide_TaskRun(t *testing.T) {
	ego := erkgroup.NewGroup(context.Background())
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	taskBatch := egokratos.NewTaskBatch[uint64, string](args)
	taskBatch.SetGlide(true)
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(idx, taskRun))
	}
	erkrequire.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erx)
	}
}

func TestGroup_Go_SetGlide_SetWaCtx_TaskRun(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond*20)
	defer cancelFunc()

	ego := erkgroup.NewGroup(ctx)
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	taskBatch := egokratos.NewTaskBatch[uint64, string](args)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *errkratos.Erk {
		return errorspb.ErrorWrongContext("ctx wrong reason=%v", err)
	})
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(idx, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
			time.Sleep(time.Millisecond * 10)
			res := strconv.FormatUint(arg, 10)
			return res, nil
		}))
	}
	erkrequire.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erx)
	}
}

func TestTaskBatch_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	var taskBatch = egokratos.NewTaskBatch[uint64, string](args)
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
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erx)
		require.Equal(t, strconv.Itoa(idx), task.Res)
		erkrequire.NoError(t, task.Erx)
	}
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"0", "1", "2", "3", "4", "5"}, results)
}

func TestTaskBatch_SetGlide_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	taskBatch := egokratos.NewTaskBatch[uint64, string](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}
	taskBatch.SetGlide(true)

	ctx := context.Background()
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(idx, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
			if arg%2 == 0 {
				return "", errorspb.ErrorServerDbError("wrong db")
			}
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run(ctx)
		t.Log(erk)
		erkrequire.NoError(t, erk) // When "glide" mode is set, no error returned here // 当设置 "平滑继续" 时这里不返回错误
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erx)
		if idx%2 == 0 {
			require.True(t, errorspb.IsServerDbError(task.Erx))
		} else {
			require.Equal(t, strconv.Itoa(idx), task.Res)
			erkrequire.NoError(t, task.Erx)
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
	taskBatch := egokratos.NewTaskBatch[uint64, string](args)
	for idx, task := range taskBatch.Tasks {
		require.Equal(t, idx, int(task.Arg))
	}
	taskBatch.SetGlide(true)

	ego := erkgroup.NewGroup(ctx)
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
		if arg%2 == 0 {
			return "", errorspb.ErrorServerDbError("wrong db")
		}
		res := strconv.FormatUint(arg, 10)
		return res, nil
	})
	erkrequire.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erx)
		if idx%2 == 0 {
			require.True(t, errorspb.IsServerDbError(task.Erx))
		} else {
			require.Equal(t, strconv.Itoa(idx), task.Res)
			erkrequire.NoError(t, task.Erx)
		}
	}
	results := taskBatch.Tasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
		return "wa-" + strconv.Itoa(int(arg))
	})
	t.Log(neatjsons.S(results))
	require.Equal(t, []string{"wa-0", "1", "wa-2", "3", "wa-4", "5"}, results)
}
