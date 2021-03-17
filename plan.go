package task_runner

import (
	"sync"
)

type Plan interface {
	SetErrorStep(stepName string)
	UsingExecutor(executorName string)
	RunBackground(task ...interface{})
	RunParallel(task ...interface{})
	RunSerial(task ...interface{})
}

type SubGrp struct {
	Policy     TaskPolicy
	CurStepIdx int
	StepNames  []string
	StepObjs   []interface{}
}

func (s *SubGrp) install(task []interface{}) bool {
	for _, stepObj := range task {
		if stepObj == nil {
			return false
		}
		s.StepObjs = append(s.StepObjs, stepObj)
	}
	return true
}

type SerErrInfo struct {
	ErrCode int
	Message string
	GrpIdx  int
	TaskIdx int
}

type Plans struct {
	SerError     *SerErrInfo
	PlanName     string
	ExecutorName string
	ErrStep      string
	PlanGrp      []SubGrp
	CurGrpIdx    int
	WtPlan       sync.WaitGroup
}

// set error step
func (b *Plans) SetErrorStep(stepName string) {
	b.ErrStep = stepName
}

// set executor name
func (b *Plans) UsingExecutor(executorName string) {
	b.ExecutorName = executorName
}

// run task in background
func (b *Plans) RunBackground(task ...interface{}) {
	b.loadTask(TaskBackground, task)
}

// run task in parallel
func (b *Plans) RunParallel(task ...interface{}) {
	b.loadTask(TaskParallel, task)
}

// run task in serial
func (b *Plans) RunSerial(task ...interface{}) {
	b.loadTask(TaskSerial, task)
}

// run task in serial with name
func (b *Plans) RunSerialName(task interface{}, name string) {
	var subGrp SubGrp
	subGrp.Policy = TaskSerial
	subGrp.install([]interface{}{task})
	b.loadData([]interface{}{task})
	subGrp.StepNames = []string{name}
	b.PlanGrp = append(b.PlanGrp, subGrp)
}

// run a task group
func (b *Plans) Try(task ...interface{}) {
	b.SetErrorStep(FINALLY)
	b.loadTask(TaskSerial, task)
}

// finally handle exception
func (b *Plans) Finally(task interface{}) {
	var subGrp SubGrp
	subGrp.Policy = TaskSerial
	subGrp.install([]interface{}{task})
	b.loadData([]interface{}{task})
	subGrp.StepNames = []string{FINALLY}
	b.PlanGrp = append(b.PlanGrp, subGrp)
}

func (b *Plans) loadTask(policy TaskPolicy, task []interface{}) {
	var subGrp SubGrp
	subGrp.Policy = policy
	subGrp.install(task)
	b.loadData(task)
	b.PlanGrp = append(b.PlanGrp, subGrp)
}

func (b *Plans) loadData(task []interface{}) bool {
	for _, stepObj := range task {
		if stepObj == nil {
			return false
		}
		stepIf, ok := stepObj.(Task)
		if !ok {
			continue
		}
		stepIf.SetSerErrInfo(b.SerError)
	}
	return true
}

// go to error handling step
func GotoErrorStep(curPlan *Plans, grpNum int) bool {
	if grpNum >= len(curPlan.PlanGrp) {
		return false
	}
	for idx, stepName := range curPlan.PlanGrp[grpNum].StepNames {
		if stepName == curPlan.ErrStep {
			if curPlan.CurGrpIdx > grpNum {
				return true
			}
			curPlan.CurGrpIdx = grpNum
			if curPlan.PlanGrp[grpNum].CurStepIdx < idx {
				curPlan.PlanGrp[grpNum].CurStepIdx = idx
			}
			return true
		}
	}
	return false
}

// record error info
func RecordErrInfo(curPlan *Plans, stepIdx int) {
	if curPlan.CurGrpIdx >= len(curPlan.PlanGrp) {
		return
	}
	curGrp := curPlan.PlanGrp[curPlan.CurGrpIdx]
	if stepIdx < 0 || stepIdx >= len(curGrp.StepObjs) {
		return
	}

	curPlan.SerError.GrpIdx = curPlan.CurGrpIdx
	curPlan.SerError.TaskIdx = stepIdx
	curStep := curGrp.StepObjs[stepIdx]
	stepIf, ok := curStep.(Task)
	if !ok {
		return
	}
	errCode, msg := stepIf.GetErrCode()
	curPlan.SerError.ErrCode = int(errCode)
	curPlan.SerError.Message = msg
}

// go to error handling process
func GotoErrorProc(curPlan *Plans) {
	curPlan.CurGrpIdx++
	for grpNum := 0; grpNum < len(curPlan.PlanGrp); grpNum++ {
		done := GotoErrorStep(curPlan, grpNum)
		if done {
			break
		}
	}
}
