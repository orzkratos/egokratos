package egokratos

import (
	"github.com/orzkratos/errkratos"
	"github.com/yyle88/egobatch"
)

// TaskBatch is a type alias using egobatch generic TaskBatch with *errkratos.Erk error type
// Manages batch task execution with concurrent processing and Kratos error handling
// Supports glide mode enabling independent task execution and fail-fast mode
// Provides context error handling and result aggregation with Kratos-specific errors
//
// TaskBatch 是使用 egobatch 泛型 TaskBatch 的类型别名，错误类型为 *errkratos.Erk
// 管理批量任务的并发执行和 Kratos 错误处理
// 支持平滑模式实现独立任务执行和快速失败模式
// 提供上下文错误处理和使用 Kratos 特定错误的结果聚合
type TaskBatch[A any, R any] = egobatch.TaskBatch[A, R, *errkratos.Erk]

// NewTaskBatch creates batch task engine with starting arguments
// Each argument becomes a task with zero-initialized result and Kratos error
// Returns TaskBatch prepared compatible with erkgroup concurrent execution
//
// NewTaskBatch 使用初始参数创建批量任务处理器
// 每个参数成为一个任务，结果和 Kratos 错误初始化成零值
// 返回准备就绪的 TaskBatch，兼容 erkgroup 并发执行
func NewTaskBatch[A any, R any](args []A) *TaskBatch[A, R] {
	return egobatch.NewTaskBatch[A, R, *errkratos.Erk](args)
}
