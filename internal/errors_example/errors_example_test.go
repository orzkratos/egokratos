package errors_example_test

import (
	"testing"

	"github.com/orzkratos/synckratos/internal/errors_example"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/erero"
)

func TestIsServerDbError(t *testing.T) {
	erk := errors_example.ErrorServerDbError("error=%v", erero.New("abc"))
	res := errors_example.IsServerDbError(erk)
	require.True(t, res)
}

func TestIsWrongContext(t *testing.T) {
	erk := errors_example.ErrorWrongContext("error=%v", erero.New("ctx"))
	res := errors_example.IsWrongContext(erk)
	require.True(t, res)
}
