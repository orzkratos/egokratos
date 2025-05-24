package example3_test

import (
	"context"
	"math/rand/v2"
	"strconv"
	"testing"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/erkmust"
	"github.com/orzkratos/synckratos/erkgroup"
	"github.com/orzkratos/synckratos/erkgroup/internal/examples/example3"
	"github.com/orzkratos/synckratos/internal/errors_example"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

type Step1Output = erkgroup.TaskOutput[*example3.Step1Param, *example3.Step1Result]
type Step2Output = erkgroup.TaskOutput[*example3.Step2Param, *example3.Step2Result]
type Step3Output = erkgroup.TaskOutput[*example3.Step3Param, *example3.Step3Result]

func TestTaskOutput(t *testing.T) {
	params := example3.NewStep1Params(5)

	outputs := processStep1s(t, params, zaplog.LOGGER)
	t.Log(neatjsons.S(outputs))

	require.Len(t, outputs.OkList(), 3)
	require.Len(t, outputs.WaList(), 2)

	require.Equal(t, 3, outputs.OkCount())
	require.Equal(t, 2, outputs.WaCount())

	t.Log(neatjsons.S(outputs.OkResults()))
	t.Log(neatjsons.S(outputs.WaReasons()))
}

func processStep1s(t *testing.T, params []*example3.Step1Param, zapLog *zaplog.Zap) erkgroup.TaskOutputList[*example3.Step1Param, *example3.Step1Result] {
	taskBatch := erkgroup.NewTaskBatch[*example3.Step1Param, *Step1Output](params)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(erx error) *errkratos.Erk {
		return errors_example.ErrorWrongContext("wrong-ctx. error=%v", erx)
	})
	ego := erkgroup.NewGroup(context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *example3.Step1Param) (*Step1Output, *errkratos.Erk) {
		return processStep1Func(t, ctx, arg, zapLog.SubZap(zap.Int("num_a", arg.NumA)))
	})
	erkmust.Done(ego.Wait())
	results := taskBatch.Tasks.Flatten(erkgroup.NewWaTaskOutput[*example3.Step1Param, *example3.Step1Result])
	require.Equal(t, len(params), len(results))

	outputs := erkgroup.TaskOutputList[*example3.Step1Param, *example3.Step1Result](results)
	return outputs
}

func processStep1Func(t *testing.T, ctx context.Context, arg *example3.Step1Param, zapLog *zaplog.Zap) (*Step1Output, *errkratos.Erk) {
	if arg.NumA%2 == 1 {
		zapLog.SUG.Debugln("wrong-a")
		return nil, errors_example.ErrorServerDbError("step-1-wrong-db")
	}
	zapLog.SUG.Debugln("process")
	res := &example3.Step1Result{
		StrA:         strconv.Itoa(arg.NumA),
		Step2Outputs: processStep2s(t, example3.NewStep2Params(1+rand.IntN(3)), zapLog),
	}
	return &Step1Output{
		Arg: arg,
		Res: res,
		Erk: nil,
	}, nil
}

func processStep2s(t *testing.T, params []*example3.Step2Param, zapLog *zaplog.Zap) erkgroup.TaskOutputList[*example3.Step2Param, *example3.Step2Result] {
	taskBatch := erkgroup.NewTaskBatch[*example3.Step2Param, *Step2Output](params)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(erx error) *errkratos.Erk {
		return errors_example.ErrorWrongContext("wrong-ctx. error=%v", erx)
	})
	ego := erkgroup.NewGroup(context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *example3.Step2Param) (*Step2Output, *errkratos.Erk) {
		return processStep2Func(t, ctx, arg, zapLog.SubZap(zap.Int("num_b", arg.NumB)))
	})
	erkmust.Done(ego.Wait())
	results := taskBatch.Tasks.Flatten(erkgroup.NewWaTaskOutput[*example3.Step2Param, *example3.Step2Result])
	require.Equal(t, len(params), len(results))

	outputs := erkgroup.TaskOutputList[*example3.Step2Param, *example3.Step2Result](results)
	return outputs
}

func processStep2Func(t *testing.T, ctx context.Context, arg *example3.Step2Param, zapLog *zaplog.Zap) (*Step2Output, *errkratos.Erk) {
	if rand.IntN(100) < 30 {
		zapLog.SUG.Debugln("wrong-b")
		return nil, errors_example.ErrorServerDbError("step-2-wrong-db")
	}
	zapLog.SUG.Debugln("process")
	res := &example3.Step2Result{
		StrB:         strconv.Itoa(arg.NumB),
		Step3Outputs: processStep3s(t, example3.NewStep3Params(1+rand.IntN(3)), zapLog),
	}
	return &Step2Output{
		Arg: arg,
		Res: res,
		Erk: nil,
	}, nil
}

func processStep3s(t *testing.T, params []*example3.Step3Param, zapLog *zaplog.Zap) erkgroup.TaskOutputList[*example3.Step3Param, *example3.Step3Result] {
	taskBatch := erkgroup.NewTaskBatch[*example3.Step3Param, *Step3Output](params)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(erx error) *errkratos.Erk {
		return errors_example.ErrorWrongContext("wrong-ctx. error=%v", erx)
	})
	ego := erkgroup.NewGroup(context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *example3.Step3Param) (*Step3Output, *errkratos.Erk) {
		return processStep3Func(t, ctx, arg, zapLog.SubZap(zap.Int("num_c", arg.NumC)))
	})
	erkmust.Done(ego.Wait())
	results := taskBatch.Tasks.Flatten(erkgroup.NewWaTaskOutput[*example3.Step3Param, *example3.Step3Result])
	require.Equal(t, len(params), len(results))

	outputs := erkgroup.TaskOutputList[*example3.Step3Param, *example3.Step3Result](results)
	return outputs
}

func processStep3Func(t *testing.T, ctx context.Context, arg *example3.Step3Param, zapLog *zaplog.Zap) (*Step3Output, *errkratos.Erk) {
	if rand.IntN(100) < 50 {
		zapLog.SUG.Debugln("wrong-c")
		return nil, errors_example.ErrorServerDbError("step-3-wrong-db")
	}
	zapLog.SUG.Debugln("process")
	res := &example3.Step3Result{
		StrC: strconv.Itoa(arg.NumC),
	}
	return &Step3Output{
		Arg: arg,
		Res: res,
		Erk: nil,
	}, nil
}
