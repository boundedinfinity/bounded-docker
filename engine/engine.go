package engine

type StateMachine struct {
}

type State interface {
	Run() (State, error)
}
