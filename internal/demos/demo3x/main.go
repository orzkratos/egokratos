package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/orzkratos/egokratos"
	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/errkratos/must/erkmust"
)

func main() {
	// Create context with 150ms timeout
	// 创建带 150ms 超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	// Create batch with task arguments
	// 创建批量任务参数
	args := []int{1, 2, 3, 4, 5}
	batch := egokratos.NewTaskBatch[int, string](args)

	// Use glide mode to see which tasks finish vs timeout
	// 使用平滑模式观察哪些任务完成、哪些超时
	batch.SetGlide(true)

	// Convert context errors to Kratos error type
	// 将上下文错误转换为 Kratos 错误类型
	batch.SetWaCtx(func(err error) *errors.Error {
		return errors.GatewayTimeout("CONTEXT_TIMEOUT", err.Error())
	})

	ego := erkgroup.NewGroup(ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *errors.Error) {
		// Each task needs different time: 50ms, 100ms, 150ms, 200ms, 250ms
		// 每个任务需要不同时间：50ms、100ms、150ms、200ms、250ms
		taskTime := time.Duration(num*50) * time.Millisecond

		timer := time.NewTimer(taskTime)
		defer timer.Stop()

		select {
		case <-timer.C:
			// Task finishes within timeout
			// 任务在超时前完成
			fmt.Printf("Task %d finished (%dms)\n", num, num*50)
			return fmt.Sprintf("task-%d", num), nil
		case <-ctx.Done():
			// Task cancelled due to timeout
			// 任务因超时而取消
			fmt.Printf("Task %d cancelled (%dms needed)\n", num, num*50)
			return "", errors.GatewayTimeout("TASK_CANCELLED", "context cancelled")
		}
	})

	// In glide mode, ego.Wait() returns nil because errors are captured in tasks
	// 在平滑模式下，ego.Wait() 返回 nil 因为错误已被捕获在任务中
	erkmust.Done(ego.Wait())

	// Show task results
	// 显示任务结果
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("\nSuccess: %d, Timeout: %d\n", len(okTasks), len(waTasks))

	// Show finished tasks
	// 显示完成的任务
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Result: %s\n", task.Arg, task.Res)
	}

	// Show timed-out tasks
	// 显示超时的任务
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Error: %s\n", task.Arg, task.Erx.Error())
	}
}
