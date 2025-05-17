package errors_example_test

import (
	"testing"

	"github.com/orzkratos/synckratos/internal/errors_example"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestIsServerDbError(t *testing.T) {
	erk := errors_example.ErrorServerDbError("error=%v", errors.New("abc"))
	res := errors_example.IsServerDbError(erk)
	require.True(t, res)
}

func TestIsWrongContext(t *testing.T) {
	erk := errors_example.ErrorWrongContext("error=%v", errors.New("ctx"))
	res := errors_example.IsWrongContext(erk)
	require.True(t, res)
}
