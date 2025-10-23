package egokratos

import (
	"github.com/orzkratos/errkratos"
	"github.com/yyle88/egobatch"
)

// TaskOutput is a type alias using egobatch generic TaskOutput with *errkratos.Erk error type
// Represents individual task execution result containing argument, result value, and Kratos error
// Provides structured output making it easy to inspect success/failure status
//
// TaskOutput 是使用 egobatch 泛型 TaskOutput 的类型别名，错误类型为 *errkratos.Erk
// 代表单个任务执行结果，包含参数、结果值和 Kratos 错误
// 提供结构化输出，便于检查成功/失败状态
type TaskOutput[ARG any, RES any] = egobatch.TaskOutput[ARG, RES, *errkratos.Erk]

// NewOkTaskOutput creates success task output with result value
// Marks task as successful with zero Kratos error
// Returns TaskOutput ready to aggregate in result collections
//
// NewOkTaskOutput 创建带有结果值的成功任务输出
// 将任务标记成功，Kratos 错误为零值
// 返回准备就绪的 TaskOutput，可以聚合到结果集合中
func NewOkTaskOutput[ARG any, RES any](arg ARG, res RES) *TaskOutput[ARG, RES] {
	return egobatch.NewOkTaskOutput[ARG, RES, *errkratos.Erk](arg, res)
}

// NewWaTaskOutput creates failed task output with Kratos error
// Marks task as failed with zero result value
// Returns TaskOutput capturing error state in Kratos format
//
// NewWaTaskOutput 创建带有 Kratos 错误的失败任务输出
// 将任务标记成失败，结果值为零值
// 返回捕获 Kratos 格式错误状态的 TaskOutput
func NewWaTaskOutput[ARG any, RES any](arg ARG, erk *errkratos.Erk) *TaskOutput[ARG, RES] {
	return egobatch.NewWaTaskOutput[ARG, RES, *errkratos.Erk](arg, erk)
}

// TaskOutputList is a type alias using egobatch generic TaskOutputList with *errkratos.Erk error type
// Represents collection of task outputs supporting batch result aggregation
// Enables filtering success/failure cases and transforming results with Kratos errors
//
// TaskOutputList 是使用 egobatch 泛型 TaskOutputList 的类型别名，错误类型为 *errkratos.Erk
// 代表支持批量结果聚合的任务输出集合
// 支持过滤成功/失败情况并使用 Kratos 错误转换结果
type TaskOutputList[ARG any, RES any] = egobatch.TaskOutputList[ARG, RES, *errkratos.Erk]
