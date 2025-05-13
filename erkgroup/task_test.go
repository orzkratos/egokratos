package erkgroup

import (
	"context"
	"strconv"
	"testing"

	"github.com/orzkratos/errkratos"
	"github.com/orzkratos/errkratos/erkrequire"
	"github.com/stretchr/testify/require"
)

func TestTaskBatch_GetRun(t *testing.T) {
	var args = []uint64{0, 1, 2, 3, 4, 5}
	var taskBatch = NewTaskBatch[uint64, string](args)
	for idx := 0; idx < len(args); idx++ {
		run := taskBatch.GetRun(context.Background(), func(ctx context.Context, arg uint64) (string, *errkratos.Erk) {
			res := strconv.FormatUint(arg, 10)
			return res, nil
		})
		erk := run()
		erkrequire.NoError(t, erk)
	}
	for idx, task := range taskBatch.Tasks {
		t.Log("idx:", idx, "arg:", task.Arg, "res:", task.Res, "erk:", task.Erk)

		//因为数据是 0 1 2 3 4 5 因此测试时恰好可以和下标比较，看看两种方式得到的结果是否相同
		require.Equal(t, idx, int(task.Arg))
		require.Equal(t, strconv.Itoa(idx), task.Res)
		erkrequire.NoError(t, task.Erk)
	}
}
