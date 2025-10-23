[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/orzkratos/egokratos/release.yml?branch=main&label=BUILD)](https://github.com/orzkratos/egokratos/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/orzkratos/egokratos)](https://pkg.go.dev/github.com/orzkratos/egokratos)
[![Coverage Status](https://img.shields.io/coveralls/github/orzkratos/egokratos/main.svg)](https://coveralls.io/github/orzkratos/egokratos?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.25%2B-lightgrey.svg)](https://github.com/orzkratos/egokratos)
[![GitHub Release](https://img.shields.io/github/release/orzkratos/egokratos.svg)](https://github.com/orzkratos/egokratos/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/orzkratos/egokratos)](https://goreportcard.com/report/github.com/orzkratos/egokratos)

# egokratos

为 Kratos 提供类型安全的批量任务处理，使用 `*errkratos.Erk` 错误处理。

基于 [egobatch](https://github.com/yyle88/egobatch) 泛型基础库构建。

---

## 特性

🎯 **Kratos 集成**: 专门为 `*errkratos.Erk` 错误类型定制
⚡ **批量处理**: 并发任务执行，类型安全的错误处理
🔄 **灵活模式**: 平滑模式和快速失败模式
🌍 **上下文支持**: 完整的上下文传播和超时处理
📋 **结果过滤**: OkTasks/WaTasks 方法聚合结果

## 安装

```bash
go get github.com/orzkratos/egokratos
```

## 快速开始

### 基础 errgroup 使用 Kratos 错误

```go
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

	// 添加任务 1：需要 100ms 完成
	ego.Go(func(ctx context.Context) *errkratos.Erk {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("任务 1 完成")
		return nil
	})

	// 添加任务 2：需要 50ms 完成
	ego.Go(func(ctx context.Context) *errkratos.Erk {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("任务 2 完成")
		return nil
	})

	// 添加任务 3：需要 80ms 完成
	ego.Go(func(ctx context.Context) *errkratos.Erk {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("任务 3 完成")
		return nil
	})

	// 等待所有任务完成并获取第一个错误（如果存在）
	if erk := ego.Wait(); erk != nil {
		fmt.Printf("发生错误: %s\n", erk.Error())
	} else {
		fmt.Println("任务完成")
	}
}
```

⬆️ **源码:** [源码](internal/demos/demo1x/main.go)

### 批量任务处理

```go
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
	// 创建批量任务
	args := []int{1, 2, 3, 4, 5}
	batch := egokratos.NewTaskBatch[int, string](args)

	// 配置平滑模式 - 即使出现错误也继续处理
	batch.SetGlide(true)

	// 执行批量任务
	ctx := context.Background()
	ego := erkgroup.NewGroup(ctx)

	batch.EgoRun(ego, func(ctx context.Context, num int) (string, *errors.Error) {
		if num%2 == 0 {
			// 偶数处理完成
			return fmt.Sprintf("even-%d", num), nil
		}
		// 奇数出现错误
		return "", errors.BadRequest("ODD_NUMBER", "odd number")
	})

	// 在平滑模式下，ego.Wait() 返回 nil 因为错误已被捕获在任务中
	erkmust.Done(ego.Wait())

	// 获取和处理任务结果
	okTasks := batch.Tasks.OkTasks()
	waTasks := batch.Tasks.WaTasks()

	fmt.Printf("成功: %d, 失败: %d\n", len(okTasks), len(waTasks))

	// 显示成功结果
	for _, task := range okTasks {
		fmt.Printf("参数: %d -> 结果: %s\n", task.Arg, task.Res)
	}

	// 显示失败结果
	for _, task := range waTasks {
		fmt.Printf("参数: %d -> 错误: %s\n", task.Arg, task.Erx.Error())
	}
}
```

⬆️ **源码:** [源码](internal/demos/demo2x/main.go)

## 核心组件

### erkgroup.Group

Kratos 的类型安全 errgroup：

```go
type Group = erxgroup.Group[*errkratos.Erk]

func NewGroup(ctx context.Context) *Group
```

### TaskBatch[A, R]

批量任务执行：

```go
type TaskBatch[A, R] = egobatch.TaskBatch[A, R, *errkratos.Erk]

func NewTaskBatch[A, R](args []A) *TaskBatch[A, R]
```

方法：
- `SetGlide(bool)` - 配置执行模式
- `SetWaCtx(func(error) *errkratos.Erk)` - 处理上下文错误
- `EgoRun(ego, func)` - 使用 errgroup 运行批量任务

### Tasks[A, R]

任务集合，支持过滤：

```go
type Tasks[A, R] = egobatch.Tasks[A, R, *errkratos.Erk]
```

方法：
- `OkTasks()` - 获取成功任务
- `WaTasks()` - 获取失败任务
- `Flatten(func)` - 转换结果

## 示例

查看 [examples](internal/examples/) 获取完整示例：

- [example1](internal/examples/example1) - 访客订单处理
- [example2](internal/examples/example2) - 学生成绩处理
- [example3](internal/examples/example3) - 多步骤流水线

## 与 egobatch 的关系

egokratos 基于 [egobatch](https://github.com/yyle88/egobatch) 使用类型别名构建：

```go
// egokratos 提供 Kratos 专用类型
type Task[A, R] = egobatch.Task[A, R, *errkratos.Erk]
type Tasks[A, R] = egobatch.Tasks[A, R, *errkratos.Erk]
type TaskBatch[A, R] = egobatch.TaskBatch[A, R, *errkratos.Erk]
```

这种方式：
- ✅ 减少代码重复
- ✅ 保持类型安全
- ✅ 提供 Kratos 友好的 API
- ✅ 受益于 egobatch 的改进

## 许可证

MIT License. 参见 [LICENSE](../LICENSE).

## 贡献

欢迎提交 Issue 和 Pull Request！
