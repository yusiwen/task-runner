package task_runner

type Executor interface {
	Plan
	getPlan() *Plans
}

type Executors struct {
	Plans
}

// init executor
func (s *Executors) Init() {
	s.SerError = &SerErrInfo{}
}

func (s *Executors) getPlan() *Plans {
	return &s.Plans
}
