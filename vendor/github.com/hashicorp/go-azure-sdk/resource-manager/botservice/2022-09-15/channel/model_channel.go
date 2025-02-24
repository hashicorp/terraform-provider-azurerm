package channel

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Channel interface {
	Channel() BaseChannelImpl
}

var _ Channel = BaseChannelImpl{}

type BaseChannelImpl struct {
	ChannelName       string  `json:"channelName"`
	Etag              *string `json:"etag,omitempty"`
	Location          *string `json:"location,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
}

func (s BaseChannelImpl) Channel() BaseChannelImpl {
	return s
}

var _ Channel = RawChannelImpl{}

// RawChannelImpl is returned when the Discriminated Value doesn't match any of the defined types
// NOTE: this should only be used when a type isn't defined for this type of Object (as a workaround)
// and is used only for Deserialization (e.g. this cannot be used as a Request Payload).
type RawChannelImpl struct {
	channel BaseChannelImpl
	Type    string
	Values  map[string]interface{}
}

func (s RawChannelImpl) Channel() BaseChannelImpl {
	return s.channel
}

func UnmarshalChannelImplementation(input []byte) (Channel, error) {
	if input == nil {
		return nil, nil
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(input, &temp); err != nil {
		return nil, fmt.Errorf("unmarshaling Channel into map[string]interface: %+v", err)
	}

	var value string
	if v, ok := temp["channelName"]; ok {
		value = fmt.Sprintf("%v", v)
	}

	if strings.EqualFold(value, "AcsChatChannel") {
		var out AcsChatChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AcsChatChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "AlexaChannel") {
		var out AlexaChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into AlexaChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DirectLineChannel") {
		var out DirectLineChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DirectLineChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "DirectLineSpeechChannel") {
		var out DirectLineSpeechChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into DirectLineSpeechChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "EmailChannel") {
		var out EmailChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into EmailChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "FacebookChannel") {
		var out FacebookChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into FacebookChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "KikChannel") {
		var out KikChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into KikChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "LineChannel") {
		var out LineChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into LineChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "M365Extensions") {
		var out M365Extensions
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into M365Extensions: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "MsTeamsChannel") {
		var out MsTeamsChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into MsTeamsChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "Omnichannel") {
		var out Omnichannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into Omnichannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "OutlookChannel") {
		var out OutlookChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into OutlookChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SearchAssistant") {
		var out SearchAssistant
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SearchAssistant: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SkypeChannel") {
		var out SkypeChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SkypeChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SlackChannel") {
		var out SlackChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SlackChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "SmsChannel") {
		var out SmsChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into SmsChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TelegramChannel") {
		var out TelegramChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TelegramChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "TelephonyChannel") {
		var out TelephonyChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into TelephonyChannel: %+v", err)
		}
		return out, nil
	}

	if strings.EqualFold(value, "WebChatChannel") {
		var out WebChatChannel
		if err := json.Unmarshal(input, &out); err != nil {
			return nil, fmt.Errorf("unmarshaling into WebChatChannel: %+v", err)
		}
		return out, nil
	}

	var parent BaseChannelImpl
	if err := json.Unmarshal(input, &parent); err != nil {
		return nil, fmt.Errorf("unmarshaling into BaseChannelImpl: %+v", err)
	}

	return RawChannelImpl{
		channel: parent,
		Type:    value,
		Values:  temp,
	}, nil

}
