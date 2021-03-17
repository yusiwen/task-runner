package task_runner

type Task interface {
	OnRequest(data string) TaskCode
	Parse(params string)
	OnFork(executor interface{}, param interface{}) int
	GetErrCode() (ErrCode, string)
	OnStop() int
	WithName(name string)
	SetSerErrInfo(serErr *SerErrInfo)
}

type Tasks struct {
	serErr     *SerErrInfo
	errMsg     string
	Name       string
	Param      []string
	resultCode ErrCode
}

// set task base name
func (t *Tasks) WithName(name string) {
	t.Name = name
}

// task base parse params
func (t *Tasks) Parse(params string) {
	t.Param = append(t.Param, params)
}

// task base on fork
func (t *Tasks) OnFork(executor interface{}, param interface{}) int {
	return 0
}

// task base on stop
func (t *Tasks) OnStop() int {
	return 0
}

// task base on request
func (t *Tasks) OnRequest(data string) TaskCode {
	return TaskFinish
}

// set task base error code
func (t *Tasks) SetFirstErrorCode(code ErrCode, msg string) {
	if t.resultCode > TaskOK {
		return
	}
	t.resultCode = code
	t.errMsg = msg
}

// get error code
func (t *Tasks) GetErrCode() (ErrCode, string) {
	return t.resultCode, t.errMsg
}

// set error info
func (t *Tasks) SetSerErrInfo(serErr *SerErrInfo) {
	t.serErr = serErr
}

// get error info
func (t *Tasks) GetSerErrInfo() *SerErrInfo {
	return t.serErr
}
