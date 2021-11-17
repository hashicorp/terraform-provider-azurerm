package eventsubscriptions

import (
	"encoding/json"
	"fmt"
)

type DeadLetterWithResourceIdentity struct {
	DeadLetterDestination DeadLetterDestination      `json:"deadLetterDestination"`
	Identity              *EventSubscriptionIdentity `json:"identity,omitempty"`
}

var _ json.Unmarshaler = &DeadLetterWithResourceIdentity{}

func (s *DeadLetterWithResourceIdentity) UnmarshalJSON(bytes []byte) error {
	type alias DeadLetterWithResourceIdentity
	var decoded alias
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling into DeadLetterWithResourceIdentity: %+v", err)
	}

	s.Identity = decoded.Identity

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("unmarshaling DeadLetterWithResourceIdentity into map[string]json.RawMessage: %+v", err)
	}

	if v, ok := temp["deadLetterDestination"]; ok {
		impl, err := unmarshalDeadLetterDestinationImplementation(v)
		if err != nil {
			return fmt.Errorf("unmarshaling field 'DeadLetterDestination' for 'DeadLetterWithResourceIdentity': %+v", err)
		}
		s.DeadLetterDestination = impl
	}
	return nil
}
