package rwlock

import (
	"fmt"
	"github.com/anishathalye/porcupine"
	"reflect"
)

const (
	GetOp = iota
	PutOp = iota
	DelOp = iota
)

type Input struct {
	Operation uint8
	Key       int
	Val       int
}

type Output struct {
	Key   int
	Val   int
	Found bool
}

type State struct {
	state map[int]int
}

var Model = porcupine.Model{
	Init: func() interface{} {
		return State{
			state: map[int]int{},
		}
	},
	Step: func(state interface{}, input interface{}, output interface{}) (bool, interface{}) {
		s := state.(State)
		i := input.(Input)
		o := output.(Output)

		stateVal, found := s.state[i.Key]

		switch i.Operation {
		case GetOp:
			if !o.Found {
				return !found, s
			} else if stateVal == o.Val {
				return true, s
			}
			break
		case PutOp:
			s.state[i.Key] = i.Val
			return true, s
		case DelOp:
			delete(s.state, i.Key)
			return true, s
		}

		return false, s
	},
	Equal: func(a, b interface{}) bool {
		return reflect.DeepEqual(a, b)
	},
	DescribeOperation: func(input interface{}, output interface{}) string {
		i := input.(Input)
		o := output.(Output)

		opName := ""
		switch i.Operation {
		case GetOp:
			opName = "Get"
			break
		case PutOp:
			opName = "Put"
			break
		case DelOp:
			opName = "Del"
			break
		}

		return fmt.Sprintf("%s(%d) -> %d", opName, i.Key, o.Val)
	},
	DescribeState: func(state interface{}) string {
		s := state.(State)
		return fmt.Sprintf("%v", s.state)
	},
}
