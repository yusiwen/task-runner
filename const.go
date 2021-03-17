package task_runner

type TaskPolicy int

const (
	_ TaskPolicy = iota
	TaskParallel
	TaskBackground
	TaskSerial
)

const DataIn string = "in"
const DataOut string = "out"

type TaskCode int

const (
	TaskFinish TaskCode = iota
)

type ErrCode int

const (
	TaskOK ErrCode = iota
	TaskFail
)

const FINALLY string = "finally"
