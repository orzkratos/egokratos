package erkgroup_test

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/erkrequire"
	"github.com/orzkratos/synckratos/erkgroup"
	"github.com/orzkratos/synckratos/internal/errors_example"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

func TestWithContext(t *testing.T) {
	ego, _ := erkgroup.WithContext(context.Background())

	for idx := 1; idx <= 50; idx++ {
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

	for idx := 1; idx <= 50; idx++ {
		num := idx
		ego.Go(func() *errkratos.Erk {
			return stepRun(ctx, num)
		})
	}

	erkrequire.Error(t, ego.Wait())
}

func stepRun(ctx context.Context, num int) *errkratos.Erk {
	if ctx.Err() != nil {
		zaplog.LOG.Info("task no", zap.Int("num", num))
		return errors_example.ErrorWrongContext("error=%v", ctx.Err())
	}
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // 模拟计算延迟
	if num%10 == 3 {
		zaplog.LOG.Info("task wa", zap.Int("num", num))
		return errors_example.ErrorServerDbError("task wa %d", num) // 模拟某个任务失败
	}
	zaplog.LOG.Info("task ok", zap.Int("num", num))
	return nil
}
