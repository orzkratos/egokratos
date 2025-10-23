package erkgroup_test

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/egokratos/internal/errorspb"
	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/must/erkrequire"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

// TestGoErrGroup demonstrates context cancellation issue with standard errgroup
// Shows how erkgroup.NewGroup wraps context to avoid shadowing problems
//
// TestGoErrGroup 演示标准 errgroup 的上下文取消问题
// 说明 erkgroup.NewGroup 如何包装上下文以避免覆盖问题
func TestGoErrGroup(t *testing.T) {
	ctx := context.Background()

	ego, ctx := errgroup.WithContext(ctx) // Use same-name ctx to shadow the old ctx, common Go idiom // 使用同名的 ctx 覆盖旧的 ctx 这是 go 里面的习惯

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func() error {
			t.Log(num)
			return nil
		})
	}

	// In Wait, cancelFunc gets invoked to cancel ctx
	// 在 wait 里面会调用 cancelFunc 把 ctx 取消掉
	require.NoError(t, ego.Wait())
	// Here ctx reports canceled, because ctx was shadowed before, and shadowing ctx is Go idiom (means errgroup design needs adjustment)
	// 这里 ctx 报已取消，因为前面 ctx 被覆盖，而覆盖ctx又是go语言开发者的习惯（就是说 errgroup 设计的不太符合习惯，需要改改）
	t.Log("ctx-err-res:", ctx.Err())
	// This is unexpected, because ctx is still needed in subsequent logic
	// NewGroup wraps errgroup.WithContext to hide ctx, see next test case
	// 这里其实是不符合预期的，因为 ctx 还要被后续逻辑用到
	// 这样就容易导致BUG，因此我在该项目里使用 NewGroup 封装 errgroup.WithContext 把 ctx 隐藏起来，具体请看下面的测试用例
	require.ErrorIs(t, checkCtx(ctx), context.Canceled)
}

// TestNewGroup demonstrates erkgroup.NewGroup solving context shadowing issue
// Validates that outer context remains unaffected when errgroup cancels
//
// TestNewGroup 演示 erkgroup.NewGroup 解决上下文覆盖问题
// 验证当 errgroup 取消时外部上下文保持不受影响
func TestNewGroup(t *testing.T) {
	ctx := context.Background()

	ego := erkgroup.NewGroup(ctx)

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func(ctx context.Context) *errkratos.Erk {
			t.Log(num)
			return nil
		})
	}

	// In Wait, cancelFunc cancels internal ctx, not the outer one
	// 在 wait 里面会调用 cancelFunc 把 ctx 取消掉，但是取消的是内部的 ctx 而不是外部的
	erkrequire.NoError(t, ego.Wait())
	// Here outer ctx is not affected
	// 这里不受影响
	t.Log("ctx-err-res:", ctx.Err())
	// Here ctx can still be used, because it's the outer ctx, unaffected when internal cancelFunc invoked
	// Group ctx is hidden, so Go and TryGo run functions must accept ctx as param
	// 这里依然可以用 ctx， 因为它是最外层的 ctx，其不受内部的 cancelFunc 的影响
	// 这样不容易出BUG，但由于 group 的 ctx 被隐藏，group 的 Go 和 TryGo 的 run 都需要是带有 ctx 信息参数的
	require.NoError(t, checkCtx(ctx))
}

func checkCtx(ctx context.Context) error {
	return ctx.Err()
}

func TestNewGroup_StepRun(t *testing.T) {
	ego := erkgroup.NewGroup(context.Background())
	ego.SetLimit(10)

	for idx := 0; idx < 50; idx++ {
		num := idx
		ego.Go(func(ctx context.Context) *errkratos.Erk {
			return stepRun(ctx, num)
		})
	}

	erkrequire.Error(t, ego.Wait())
}

func stepRun(ctx context.Context, idx int) *errkratos.Erk {
	if ctx.Err() != nil {
		zaplog.LOG.Info("task no", zap.Int("num", idx))
		return errorspb.ErrorWrongContext("error=%v", ctx.Err())
	}
	time.Sleep(time.Duration(rand.IntN(1000)) * time.Millisecond) // Simulate computation time // 模拟计算延迟
	if idx%10 == 3 {
		zaplog.LOG.Info("task wa", zap.Int("num", idx))
		return errorspb.ErrorServerDbError("task wa %d", idx) // Simulate task execution failure // 模拟某个任务失败
	}
	zaplog.LOG.Info("task ok", zap.Int("num", idx))
	return nil
}
