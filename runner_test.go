package task_runner

import (
	"fmt"
	"testing"
)

type TestTask1 struct {
	Tasks
}

func (t *TestTask1) OnRequest(data string) TaskCode {
	fmt.Println("In TestTask1")
	return TaskFinish
}

type TestTask2 struct {
	Tasks
}

func (t *TestTask2) OnRequest(data string) TaskCode {
	fmt.Println("In TestTask2")
	return TaskFinish
}

func TestTryFinallyRunner(t *testing.T) {
	var executor = Executors{}
	executor.Init()
	executor.Try(
		&TestTask1{},
		&TestTask1{})
	executor.Finally(&TestTask2{})
	Run(&executor)

	fmt.Println("OK")
}

func TestBackgroundRunner(t *testing.T) {
	var executor = Executors{}
	executor.Init()
	executor.RunBackground(&TestTask1{}, &TestTask2{}, &TestTask1{})
	Run(&executor)

	fmt.Println("OK")
}

func TestParallelRunner(t *testing.T) {
	var executor = Executors{}
	executor.Init()
	executor.RunParallel(&TestTask1{}, &TestTask2{}, &TestTask1{})
	Run(&executor)

	fmt.Println("OK")
}
