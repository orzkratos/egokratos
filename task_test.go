package egokratos_test

import (
	"strconv"
	"testing"

	"github.com/orzkratos/egokratos"
	"github.com/orzkratos/egokratos/internal/errors_example"
	"github.com/orzkratos/errkratos"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/neatjson/neatjsons"
)

func TestTasks_OkTasks(t *testing.T) {
	var tasks = make(egokratos.Tasks[uint64, string], 0, 10)
	for idx := 0; idx < 10; idx++ {
		if idx%2 == 0 {
			tasks = append(tasks, &egokratos.Task[uint64, string]{
				Arg: uint64(idx),
				Res: strconv.Itoa(idx),
				Erk: nil,
			})
		} else {
			tasks = append(tasks, &egokratos.Task[uint64, string]{
				Arg: uint64(idx),
				Res: "",
				Erk: errors_example.ErrorServerDbError("wrong-db"),
			})
		}
	}
	t.Run("ok", func(t *testing.T) {
		okTasks := tasks.OkTasks()
		t.Log(neatjsons.S(okTasks))
		require.Len(t, okTasks, 5)
		results := okTasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
			panic("impossible")
		})
		t.Log(neatjsons.S(results))
		require.Equal(t, []string{"0", "2", "4", "6", "8"}, results)
	})

	t.Run("wa", func(t *testing.T) {
		waTasks := tasks.WaTasks()
		t.Log(neatjsons.S(waTasks))
		require.Len(t, waTasks, 5)
		results := waTasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
			return "wa-" + strconv.FormatUint(arg, 10)
		})
		t.Log(neatjsons.S(results))
		require.Equal(t, []string{"wa-1", "wa-3", "wa-5", "wa-7", "wa-9"}, results)
	})

	t.Run("do", func(t *testing.T) {
		results := tasks.Flatten(func(arg uint64, erk *errkratos.Erk) string {
			return "wa-" + strconv.FormatUint(arg, 10)
		})
		require.Len(t, results, 10)
		t.Log(neatjsons.S(results))
		require.Equal(t, []string{"0", "wa-1", "2", "wa-3", "4", "wa-5", "6", "wa-7", "8", "wa-9"}, results)
	})
}
