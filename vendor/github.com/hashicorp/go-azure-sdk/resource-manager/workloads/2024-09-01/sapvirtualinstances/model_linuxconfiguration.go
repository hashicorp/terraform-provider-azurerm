package sapvirtualinstances

import (
	"encoding/json"
	"fmt"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ OSConfiguration = LinuxConfiguration{}

type LinuxConfiguration struct {
	DisablePasswordAuthentication *bool             `json:"disablePasswordAuthentication,omitempty"`
	Ssh                           *SshConfiguration `json:"ssh,omitempty"`
	SshKeyPair                    *SshKeyPair       `json:"sshKeyPair,omitempty"`

	// Fields inherited from OSConfiguration

	OsType OSType `json:"osType"`
}

func (s LinuxConfiguration) OSConfiguration() BaseOSConfigurationImpl {
	return BaseOSConfigurationImpl{
		OsType: s.OsType,
	}
}

var _ json.Marshaler = LinuxConfiguration{}

func (s LinuxConfiguration) MarshalJSON() ([]byte, error) {
	type wrapper LinuxConfiguration
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling LinuxConfiguration: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling LinuxConfiguration: %+v", err)
	}

	decoded["osType"] = "Linux"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling LinuxConfiguration: %+v", err)
	}

	return encoded, nil
}
