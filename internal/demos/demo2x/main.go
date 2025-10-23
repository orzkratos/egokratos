package main

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/orzkratos/egokratos"
	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/errkratos/must/erkmust"
)

func main() {
	// Create batch with arguments
	// 使用参数创建批量任务
	args := []int{1, 2, 3, 4, 5}
	batch := egokratos.NewTaskBatch[int, string](args)

	// Configure glide mode - keep going even when errors happen
	// 配置平滑模式 - 即使出现错误也继续处理
	batch.SetGlide(true)

	// Execute batch tasks
	// 执行批量任务
	ctx := context.Background()
	ego := erkgroup.NewGroup(ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *errors.Error) {
		if num%2 == 0 {
			// Even numbers finish OK
			// 偶数处理完成
			return fmt.Sprintf("even-%d", num), nil
		}
		// Odd numbers have errors
		// 奇数出现错误
		return "", errors.BadRequest("ODD_NUMBER", "odd number")
	})

	// In glide mode, ego.Wait() returns nil because errors are captured in tasks
	// 在平滑模式下，ego.Wait() 返回 nil 因为错误已被捕获在任务中
	erkmust.Done(ego.Wait())

	// Get and handle task results
	// 获取和处理任务结果
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("Success: %d, Failed: %d\n", len(okTasks), len(waTasks))

	// Show OK results
	// 显示成功结果
	for _, task := range okTasks {
		fmt.Printf("Arg: %d -> Result: %s\n", task.Arg, task.Res)
	}

	// Show failed results
	// 显示失败结果
	for _, task := range waTasks {
		fmt.Printf("Arg: %d -> Error: %s\n", task.Arg, task.Erx.Error())
	}
}
