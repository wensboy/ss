package err

import (
	"errors"
	"fmt"
)

func Tracing(err error, target error) bool {
	return errors.Is(err, target)
}

type Err struct {
	Errno string `json:"errno,omitempty"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	err   error  `json:"-"`
}

func NewErr(errno string, code int, template string, args ...any) Err {
	return Err{
		Errno: errno,
		Code:  code,
		Msg:   fmt.Sprintf(template, args...),
	}
}

func (e Err) Error() string {
	str := fmt.Sprintf("[%s#%d] %s", e.Errno, e.Code, e.Msg)
	if e.Errno == "" {
		str = fmt.Sprintf("[%d] %s", e.Code, e.Msg)
	}
	if e.err == nil {
		return str
	}
	return fmt.Sprintf("%s -> %s", str, e.err.Error())
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
