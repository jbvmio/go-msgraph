package msgraph

import (
	"encoding/json"
	"fmt"
	"strings"
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

// ParseAppRestartBehavior takes the given string and attempts to match it to a corresponding
// AppRestartBehavior ID. If the parse is unsuccessful, the default returned AppRestartBehavior
// will be of 0 (RestartBasedOnReturnCode).
func ParseAppRestartBehavior(v string) AppRestartBehavior {
	var rb AppRestartBehavior
	val := strings.ToLower(v)
	switch val {
	case `basedonreturncode`:
		rb = RestartBasedOnReturnCode
	case `allow`:
		rb = RestartAllow
	case `suppress`:
		rb = RestartSuppress
	case `force`:
		rb = RestartForce
	}
	return rb
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
