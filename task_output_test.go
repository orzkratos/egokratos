package egokratos_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/orzkratos/egokratos"
	"github.com/orzkratos/egokratos/erkgroup"
	"github.com/orzkratos/egokratos/internal/errors_example"
	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/erkmust"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
)

func TestTaskOutput(t *testing.T) {
	type Param struct {
		Value int
	}

	type Result struct {
		Value string
	}

	var args []*Param
	for _, num := range []int{0, 1, 2, 3, 4, 5} {
		args = append(args, &Param{Value: num})
	}

	taskBatch := egokratos.NewTaskBatch[*Param, *egokratos.TaskOutput[*Param, *Result]](args)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(erx error) *errkratos.Erk {
		return errors_example.ErrorWrongContext("wrong-ctx. error=%v", erx)
	})
	ego := erkgroup.NewGroup(context.Background())
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, func(ctx context.Context, arg *Param) (*egokratos.TaskOutput[*Param, *Result], *errkratos.Erk) {
		if arg.Value%3 == 2 {
			return nil, errors_example.ErrorServerDbError("wrong-db")
		}
		res := &Result{Value: strconv.Itoa(arg.Value)}
		return egokratos.NewOkTaskOutput[*Param, *Result](arg, res), nil
	})
	erkmust.Done(ego.Wait())
	results := taskBatch.Tasks.Flatten(egokratos.NewWaTaskOutput[*Param, *Result])

	ops := egokratos.TaskOutputList[*Param, *Result](results)
	t.Log(neatjsons.S(ops))

	require.Len(t, ops.OkList(), 4)
	require.Len(t, ops.WaList(), 2)

	require.Equal(t, 4, ops.OkCount())
	require.Equal(t, 2, ops.WaCount())

	t.Log(neatjsons.S(ops.OkResults()))
	t.Log(neatjsons.S(ops.WaReasons()))
}
