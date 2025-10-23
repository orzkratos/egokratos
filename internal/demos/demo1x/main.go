package main

import (
	"context"
	"fmt"
	"time"

	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/errkratos"
)

func main() {
	ctx := context.Background()
	ego := erkgroup.NewGroup(ctx)

	// Add task 1: takes 100ms to finish
	// 添加任务 1：需要 100ms 完成
	ego.Go(func(ctx context.Context) *errkratos.Erk {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("Task 1 finished OK")
		return nil
	})

	// Add task 2: takes 50ms to finish
	// 添加任务 2：需要 50ms 完成
	ego.Go(func(ctx context.Context) *errkratos.Erk {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Task 2 finished OK")
		return nil
	})

	// Add task 3: takes 80ms to finish
	// 添加任务 3：需要 80ms 完成
	ego.Go(func(ctx context.Context) *errkratos.Erk {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("Task 3 finished OK")
		return nil
	})

	// Wait until tasks finish and get the first error
	// 等待所有任务完成并获取第一个错误（如果存在）
	if erk := ego.Wait(); erk != nil {
		fmt.Printf("Got error: %s\n", erk.Error())
	} else {
		fmt.Println("Tasks finished OK")
	}
}
