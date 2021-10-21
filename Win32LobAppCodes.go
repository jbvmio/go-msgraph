package msgraph

import (
	"encoding/json"
	"fmt"
)

type AppRestartBehavior int

func (A AppRestartBehavior) String() string {
	return appRestartBehaviorStrings[A]
}

const (
	RestartBasedOnReturnCode AppRestartBehavior = iota
	RestartAllow
	RestartSuppress
	RestartForce
)

var appRestartBehaviorStrings = [...]string{
	`basedOnReturnCode`,
	`allow`,
	`suppress`,
	`force`,
}

func (A *AppRestartBehavior) UnmarshalJSON(data []byte) error {
	if string(data) == `null` {
		*A = 0
		return nil
	}
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	switch tmp {
	case `null`, `basedOnReturnCode`:
		*A = 0
	case `allow`:
		*A = 1
	case `suppress`:
		*A = 2
	case `force`:
		*A = 3
	default:
		err = fmt.Errorf("invalid restart behavior - must be one of either %s, %s, %s or %s", RestartBasedOnReturnCode, RestartAllow, RestartSuppress, RestartForce)
	}
	return err
}

func (A AppRestartBehavior) MarshalJSON() ([]byte, error) {
	switch A {
	case 0, 1, 2, 3:
		return json.Marshal(A.String())
	default:
		return []byte{}, fmt.Errorf("invalid restart behavior - must be one of either %s, %s, %s or %s", RestartBasedOnReturnCode, RestartAllow, RestartSuppress, RestartForce)
	}
}
