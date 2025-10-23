// Package erkgroup provides type-safe errgroup wrapping specialized with Kratos framework error handling
// Built on erxgroup generic foundation using type alias with *errkratos.Erk custom error type
// Enables concurrent task execution with Kratos-specific error propagation and context cancellation
// Maintains errgroup semantics while providing zero-cost Kratos framework integration
//
// Package erkgroup 提供专门与 Kratos 框架错误处理集成的类型安全 errgroup 包装
// 基于 erxgroup 泛型基础，使用类型别名和 *errkratos.Erk 自定义错误类型构建
// 支持并发任务执行，包含 Kratos 特定的错误传播和上下文取消
// 保持 errgroup 语义，同时提供零成本的 Kratos 框架集成
package erkgroup

import (
	"context"

	"github.com/orzkratos/errkratos"
	"github.com/yyle88/egobatch/erxgroup"
)

// Group is a type alias using erxgroup generic errgroup with *errkratos.Erk error type
// Wraps errgroup.Group with type-safe Kratos error handling and context propagation
// Maintains goroutine synchronization semantics while enabling Kratos error integration
//
// Group 是使用 erxgroup 泛型 errgroup 的类型别名，错误类型为 *errkratos.Erk
// 使用类型安全的 Kratos 错误处理和上下文传播包装 errgroup.Group
// 保持协程同步语义，同时支持 Kratos 错误集成
type Group = erxgroup.Group[*errkratos.Erk]

// NewGroup creates generic errgroup with Kratos error type and context cancellation
// Context cancels when first error occurs or parent context cancels
// Returns Group prepared with concurrent task execution using Kratos errors
//
// NewGroup 创建带有 Kratos 错误类型和上下文取消的泛型 errgroup
// 当第一个错误发生或父上下文取消时，上下文会被取消
// 返回准备就绪的 Group，可以使用 Kratos 错误执行并发任务
func NewGroup(ctx context.Context) *Group {
	return erxgroup.NewGroup[*errkratos.Erk](ctx)
}
