// Package egokratos: Type-safe batch task processing specialized with Kratos framework error handling
// Built on egobatch generic foundation using type aliases with *errkratos.Erk custom error type
// Provides zero-cost abstraction while maintaining complete egobatch capabilities and Kratos error semantics
// Supports concurrent task execution, result filtering (OkTasks/WaTasks), and flexible error handling modes
//
// egokratos: 专门与 Kratos 框架错误处理集成的类型安全批量任务处理
// 基于 egobatch 泛型基础，使用类型别名和 *errkratos.Erk 自定义错误类型构建
// 提供零成本抽象，同时保持完整的 egobatch 能力和 Kratos 错误语义
// 支持并发任务执行、结果过滤（OkTasks/WaTasks）和灵活的错误处理模式
package egokratos

import (
	"github.com/orzkratos/errkratos"
	"github.com/yyle88/egobatch"
)

// Task is a type alias using egobatch generic Task with *errkratos.Erk error type
// Represents single task containing argument, result, and Kratos-specific error
// Inherits each egobatch.Task method while providing Kratos framework integration
//
// Task 是使用 egobatch 泛型 Task 的类型别名，错误类型为 *errkratos.Erk
// 代表包含参数、结果和 Kratos 特定错误的单个任务
// 继承所有 egobatch.Task 方法，同时提供 Kratos 框架集成
type Task[A any, R any] = egobatch.Task[A, R, *errkratos.Erk]

// Tasks is a type alias using egobatch generic Tasks with *errkratos.Erk error type
// Provides task collection supporting filtering and transformation operations
// Enables OkTasks/WaTasks filtering and Flatten transformation with Kratos errors
//
// Tasks 是使用 egobatch 泛型 Tasks 的类型别名，错误类型为 *errkratos.Erk
// 提供支持过滤和转换操作的任务集合
// 支持 OkTasks/WaTasks 过滤和使用 Kratos 错误的 Flatten 转换
type Tasks[A any, R any] = egobatch.Tasks[A, R, *errkratos.Erk]
