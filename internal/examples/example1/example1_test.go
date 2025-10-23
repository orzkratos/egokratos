package example1_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/orzkratos/egokratos"
	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/egokratos/internal/errorspb"
	"github.com/orzkratos/egokratos/internal/examples/example1"
	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/must/erkmust"
	"github.com/yyle88/neatjson/neatjsons"
)

// TestRun demonstrates guest order processing with nested batch tasks
// Shows two-stage processing: guests -> orders with error handling at each stage
//
// TestRun 演示访客订单处理与嵌套批量任务
// 展示两阶段处理：访客 -> 订单，每个阶段都有错误处理
func TestRun(t *testing.T) {
	ctx := context.Background()
	guests := example1.NewGuests(10)
	taskResults := processGuests(ctx, guests)
	// Flatten results to avoid nested generic output
	// 展平结果避免嵌套泛型输出
	guestOrdersStates := taskResults.Flatten(func(guest *example1.Guest, erk *errkratos.Erk) *example1.GuestOrdersStates {
		return &example1.GuestOrdersStates{
			Guest:       guest,
			OrderStates: nil,
			Outline:     "",
			Erk:         erk,
		}
	})
	t.Log(neatjsons.S(guestOrdersStates))
}

func processGuests(ctx context.Context, guests []*example1.Guest) egokratos.Tasks[*example1.Guest, *example1.GuestOrdersStates] {
	taskBatch := egokratos.NewTaskBatch[*example1.Guest, *example1.GuestOrdersStates](guests)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *errkratos.Erk {
		return errorspb.ErrorWrongContext("wrong-ctx-can-not-invoke-process-guest-func. error=%v", err)
	})
	ego := erkgroup.NewGroup(ctx)
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, processGuestFunc)
	erkmust.Done(ego.Wait())
	return taskBatch.Tasks
}

func processGuestFunc(ctx context.Context, guest *example1.Guest) (*example1.GuestOrdersStates, *errkratos.Erk) {
	if rand.IntN(2) == 0 {
		return nil, errorspb.ErrorServerDbError("wrong-db")
	}
	orderCount := 1 + rand.IntN(5)
	orders := example1.NewOrders(guest, orderCount)

	taskResults := processOrders(ctx, orders)

	// Flatten task results to reduce nesting depth and improve code structure
	// 这里把数据降低维度，避免泛型套泛型，能够让逻辑更清楚些，直接返回这个 task-results 也是可以的
	orderStates := taskResults.Flatten(func(order *example1.Order, erk *errkratos.Erk) *example1.OrderState {
		return &example1.OrderState{
			Order: order,
			Erk:   erk,
		}
	})

	outline := createStatusOutline(orderStates)

	return &example1.GuestOrdersStates{
		Guest:       guest,
		OrderStates: orderStates,
		Outline:     outline,
		Erk:         nil,
	}, nil
}

func processOrders(ctx context.Context, orders []*example1.Order) egokratos.Tasks[*example1.Order, *example1.OrderState] {
	taskBatch := egokratos.NewTaskBatch[*example1.Order, *example1.OrderState](orders)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *errkratos.Erk {
		return errorspb.ErrorWrongContext("wrong-ctx-can-not-invoke-process-order-func. error=%v", err)
	})
	ego := erkgroup.NewGroup(ctx)
	ego.SetLimit(2)
	taskBatch.EgoRun(ego, processOrderFunc)
	erkmust.Done(ego.Wait())
	return taskBatch.Tasks
}

func processOrderFunc(ctx context.Context, order *example1.Order) (*example1.OrderState, *errkratos.Erk) {
	if rand.IntN(2) == 0 {
		return nil, errorspb.ErrorServerDbError("wrong-db")
	}
	return &example1.OrderState{
		Order: order,
		State: "OK",
		Erk:   nil,
	}, nil
}

func createStatusOutline(orderStates []*example1.OrderState) string {
	okCount := 0
	waCount := 0
	for _, state := range orderStates {
		if state.Erk != nil {
			waCount++
		} else {
			okCount++
		}
	}
	return fmt.Sprintf("ok-count:%d wa-count=%d", okCount, waCount)
}
