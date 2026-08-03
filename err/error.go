package err

import (
	"errors"
	"fmt"
	"sync"
)

const (
	Wrap   = 1
	Unwrap = 0

	PackageConfig = iota + 1
	PackageServer
	PackageDB
	PackageLog
	PackageErr
	PackageCmd
)

var (
	_global_errhub *ErrHub
)

func init() {
	_global_errhub = NewErrHub()
}

func GetGErrHub() *ErrHub {
	return _global_errhub
}

func Tracing(err error, target error) bool {
	return errors.Is(err, target)
}

type Err struct {
	Errno string `json:"errno,omitempty"`
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	err   error  `json:"-"`
}

func NewErr(errno string, code int, template string, args ...any) *Err {
	return &Err{
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

type ErrEntry struct {
	msg  string
	perr *Err
}

type ErrCode interface {
	Code() int
}

type NormalErrCode struct {
	wrapBase      int
	packageBase   int
	componentBase int
	code          int
}

func NewNormalErrCode(encode [4]int) NormalErrCode {
	nec := NormalErrCode{
		wrapBase:      10_000_000,
		packageBase:   100_000,
		componentBase: 1_000,
	}
	nec.code = nec.wrapBase*(encode[0]%2) + nec.packageBase*(encode[1]%100) + nec.componentBase*(encode[2]%100) + encode[3]
	return nec
}

func (nec NormalErrCode) Code() int {
	return nec.code
}

type BusinessErrCode struct {
	wrapBase     int
	httpBase     int
	businessBase int
	code         int
}

func NewBusinessErrCode(encode [4]int) BusinessErrCode {
	bec := BusinessErrCode{
		wrapBase:     10_000_000,
		httpBase:     10_000,
		businessBase: 100,
	}
	bec.code = bec.wrapBase*(encode[0]%2) + bec.httpBase*(encode[1]%600) + bec.businessBase*(encode[2]%100) + encode[3]
	return bec
}

func (nec BusinessErrCode) Code() int {
	return nec.code
}

type ErrHub struct {
	mu   sync.RWMutex
	errs map[int]*ErrEntry
}

func NewErrHub() *ErrHub {
	return &ErrHub{
		errs: make(map[int]*ErrEntry),
	}
}

func (eh *ErrHub) GetEntry(code ErrCode) *ErrEntry {
	return eh.errs[code.Code()]
}

func (eh *ErrHub) GetMsg(code ErrCode) (string, bool) {
	eh.mu.RLock()
	defer eh.mu.RUnlock()
	entry := eh.GetEntry(code)
	if entry == nil {
		return "", false
	}
	return entry.msg, true
}

func (eh *ErrHub) GetErr(code ErrCode) (error, bool) {
	eh.mu.RLock()
	defer eh.mu.RUnlock()
	entry := eh.GetEntry(code)
	if entry == nil {
		return Err{}, false
	}
	if entry.perr == nil {
		// lazy set
		entry.perr = &Err{
			Code: code.Code(),
			Msg:  entry.msg,
		}
	}
	return entry.perr, true
}

func (eh *ErrHub) Set(code ErrCode, err *Err) {
	eh.mu.Lock()
	defer eh.mu.Unlock()
	eh.errs[code.Code()] = &ErrEntry{
		msg:  err.Msg,
		perr: err,
	}
}

func (eh *ErrHub) LazySet(code ErrCode, msg string) {
	eh.mu.Lock()
	defer eh.mu.Unlock()
	eh.errs[code.Code()] = &ErrEntry{
		msg: msg,
	}
}

func (eh *ErrHub) GetOrSet(code ErrCode, err *Err) error {
	entry := eh.GetEntry(code)
	if entry != nil {
		return entry.perr
	}
	eh.Set(code, err)
	return err
}
