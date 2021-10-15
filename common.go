package msgraph

import (
	"encoding/json"
	"fmt"
	"time"
)

// TypeAndID only contains the ODataType, ODataContext and ID.
type TypeAndID struct {
	ODataType    string `json:"@odata.type,omitempty"`
	ODataContext string `json:"@odata.context,omitempty"`
	ID           string `json:"id,omitempty"`
}

type DateTimeOffset time.Time

func (dto *DateTimeOffset) String() string {
	return time.Time(*dto).Format(time.RFC3339Nano)
}

func (dto *DateTimeOffset) UnmarshalJSON(data []byte) error {
	if string(data) == `null` {
		*dto = DateTimeOffset(time.Time{})
		return nil
	}
	var tmp string
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339Nano, tmp)
	if err != nil {
		return fmt.Errorf("cannot parse DateTimeOffset %q with RFC3339: %v", tmp, err)
	}
	*dto = DateTimeOffset(t)
	return nil
}

func (dto DateTimeOffset) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dto).Format(time.RFC3339Nano))
}
