package erkgroup_test

import (
	"context"
	"math/rand/v2"
	"strconv"
	"testing"
	"time"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/erkrequire"
	"github.com/orzkratos/synckratos/erkgroup"
	"github.com/orzkratos/synckratos/internal/errors_example"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func TestWithContext(t *testing.T) {
	ego, _ := erkgroup.WithContext(context.Background())

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func() *errkratos.Erk {
			t.Log(num)
			return nil
		})
	}

	erkrequire.NoError(t, ego.Wait())
}

func TestWithContext_StepRun(t *testing.T) {
	ego, ctx := erkgroup.WithContext(context.Background())
	ego.SetLimit(10)

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func() *errkratos.Erk {
			return stepRun(ctx, num)
		})
	}

	erkrequire.Error(t, ego.Wait())
}

func stepRun(ctx context.Context, idx int) *errkratos.Erk {
	if ctx.Err() != nil {
		zaplog.LOG.Info("task no", zap.Int("num", idx))
		return errors_example.ErrorWrongContext("error=%v", ctx.Err())
	}
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // 模拟计算延迟
	if idx%10 == 3 {
		zaplog.LOG.Info("task wa", zap.Int("num", idx))
		return errors_example.ErrorServerDbError("task wa %d", idx) // 模拟某个任务失败
	}
	zaplog.LOG.Info("task ok", zap.Int("num", idx))
	return nil
}

func TestGroup_Go_TaskRun(t *testing.T) {
	ego, ctx := erkgroup.WithContext(context.Background())
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	var taskBatch = erkgroup.NewTaskBatch[uint64, string](args)
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(ctx, taskRun))
	}
	erkrequire.Error(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)
	}
	require.Equal(t, int64(50), taskBatch.Index)
}

func taskRun(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
	if ctx.Err() != nil {
		zaplog.LOG.Info("task no", zap.Uint64("arg", arg))
		return "", errors_example.ErrorWrongContext("error=%v", ctx.Err())
	}
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // 模拟计算延迟
	if arg%10 == 3 {
		zaplog.LOG.Info("task wa", zap.Uint64("arg", arg))
		return "", errors_example.ErrorServerDbError("task wa %d", arg) // 模拟某个任务失败
	}
	zaplog.LOG.Info("task ok", zap.Uint64("arg", arg))

	res := strconv.FormatUint(arg, 10)
	return res, nil
}

func TestGroup_Go_SetGlide_TaskRun(t *testing.T) {
	ego, ctx := erkgroup.WithContext(context.Background())
	ego.SetLimit(10)

	args := make([]uint64, 0, 50)
	for num := uint64(0); num < 50; num++ {
		args = append(args, num)
	}

	var taskBatch = erkgroup.NewTaskBatch[uint64, string](args)
	taskBatch.SetGlide(true)
	for idx := 0; idx < 50; idx++ {
		ego.Go(taskBatch.GetRun(ctx, taskRun))
	}
	erkrequire.NoError(t, ego.Wait())

	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)
	}
	require.Equal(t, int64(50), taskBatch.Index)
}
