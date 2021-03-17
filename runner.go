package task_runner

import (
	"sync"
)

// run plans in executor
func Run(executor Executor) ErrCode {
	curPlan := executor.getPlan()
	for {
		if curPlan.CurGrpIdx >= len(curPlan.PlanGrp) {
			break
		}
		curSub := &curPlan.PlanGrp[curPlan.CurGrpIdx]
		retCode, stepIdx := grpRun(curSub, executor, &curPlan.WtPlan)
		if retCode <= TaskOK {
			curPlan.CurGrpIdx++
			continue
		}
		RecordErrInfo(curPlan, stepIdx)
		GotoErrorProc(curPlan)
	}
	// wait background job finish
	curPlan.WtPlan.Wait()
	return TaskOK
}

// task runner
func taskRunner(executor interface{}, task Task) int {
	for {
		loadObjByInd(task, executor, DataIn)
		retCode := task.OnRequest("")
		if retCode <= TaskFinish {
			loadObjByInd(task, executor, DataOut)
			break
		}
	}
	return 0
}

// step run policy: background, parallel, serial
func step(wg *sync.WaitGroup, curSub *SubGrp, executor Executor, wtPlan *sync.WaitGroup, task Task) ErrCode {
	taskRet := TaskOK
	switch curSub.Policy {
	case TaskBackground:
		wtPlan.Add(1)
		go func() {
			taskRunner(executor, task)
			wtPlan.Done()
		}()

	case TaskParallel:
		wg.Add(1)
		go func() {
			taskRunner(executor, task)
			wg.Done()
		}()

	default:
		taskRunner(executor, task)
		taskRet, _ = task.GetErrCode()
	}

	return taskRet
}

// run one step
func grpOneStep(wg *sync.WaitGroup, curSub *SubGrp, executor Executor, wtPlan *sync.WaitGroup) bool {
	if curSub.CurStepIdx >= len(curSub.StepObjs) {
		return false
	}
	curStep := curSub.StepObjs[curSub.CurStepIdx]
	if curStep == nil {
		curSub.CurStepIdx++
		return true
	}
	task, ok := curStep.(Task)
	if !ok {
		return false
	}
	taskRet := step(wg, curSub, executor, wtPlan, task)
	curSub.CurStepIdx++

	return taskRet <= TaskOK
}

func grpGetRetCode(curSub *SubGrp) (ErrCode, int) {
	for idx, curStep := range curSub.StepObjs {
		task, ok := curStep.(Task)
		if !ok {
			continue
		}
		errCode, _ := task.GetErrCode()
		if errCode > TaskOK {
			return errCode, idx
		}
	}

	return TaskOK, -1
}

func grpRun(curSub *SubGrp, executor Executor, wtPlan *sync.WaitGroup) (ErrCode, int) {
	var wg sync.WaitGroup
	for {
		nextStep := grpOneStep(&wg, curSub, executor, wtPlan)
		if !nextStep {
			break
		}
	}
	if curSub.Policy == TaskParallel {
		wg.Wait()
	}
	return grpGetRetCode(curSub)
}
