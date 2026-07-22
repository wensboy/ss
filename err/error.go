package err

import (
	"errors"
	"fmt"
)

type Err struct {
	Errno string `json:"errno,omitempty"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	err   error  `json:"-"`
}

func (e Err) Error() string {
	if e.err == nil {
		return fmt.Sprintf("[%s#%d] %s", e.Errno, e.Code, e.Msg)
	}
	return fmt.Sprintf("[%s#%d] %s -> %s", e.Errno, e.Code, e.Msg, e.err.Error())
}

func (e Err) Wrap(err error) Err {
	return Err{
		Errno: e.Errno,
		Code:  e.Code,
		Msg:   e.Msg,
		err:   err,
	}
}

func (e Err) Unwrap() error {
	return e.err
}

func NewErr(errno string, code int, template string, args ...any) Err {
	return Err{
		Errno: errno,
		Code:  code,
		Msg:   fmt.Sprintf(template, args...),
	}
}

func Tracing(err error, target error) bool {
	return errors.Is(err, target)
}
