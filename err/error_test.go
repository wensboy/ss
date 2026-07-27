package err

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrRoot = NewErr("E000", 0, "Root error")

func Test_Err(t *testing.T) {
	err := NewErr("E001", 400, "Invalid input: %s", "missing field").Wrap(ErrRoot)
	assert.Equal(t, "[E001#400] Invalid input: missing field -> [E000#0] Root error", err.Error())
}

func Test_Tracing(t *testing.T) {
	err := NewErr("E001", 400, "Invalid input: %s", "missing field").Wrap(ErrRoot)
	assert.True(t, Tracing(err, ErrRoot))
}

func Test_Unwrap(t *testing.T) {
	err := NewErr("E001", 400, "Invalid input: %s", "missing field").Wrap(ErrRoot)
	assert.Equal(t, ErrRoot, errors.Unwrap(err))
}

func Test_ErrCode(t *testing.T) {
	assert.Equal(t, 101001, NewNormalErrCode([4]int{0, 1, 1, 1}).Code())
	assert.Equal(t, 10101001, NewNormalErrCode([4]int{1, 1, 1, 1}).Code())
	assert.Equal(t, 2000101, NewBusinessErrCode([4]int{0, 200, 1, 1}).Code())
	assert.Equal(t, 12000101, NewBusinessErrCode([4]int{1, 200, 1, 1}).Code())
}
