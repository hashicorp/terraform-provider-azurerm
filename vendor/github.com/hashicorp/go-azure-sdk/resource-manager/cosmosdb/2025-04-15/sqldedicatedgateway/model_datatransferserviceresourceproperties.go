package sqldedicatedgateway

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ ServiceResourceProperties = DataTransferServiceResourceProperties{}

type DataTransferServiceResourceProperties struct {
	Locations *[]RegionalServiceResource `json:"locations,omitempty"`

	// Fields inherited from ServiceResourceProperties

	CreationTime  *string        `json:"creationTime,omitempty"`
	InstanceCount *int64         `json:"instanceCount,omitempty"`
	InstanceSize  *ServiceSize   `json:"instanceSize,omitempty"`
	ServiceType   ServiceType    `json:"serviceType"`
	Status        *ServiceStatus `json:"status,omitempty"`
}

func (s DataTransferServiceResourceProperties) ServiceResourceProperties() BaseServiceResourcePropertiesImpl {
	return BaseServiceResourcePropertiesImpl{
		CreationTime:  s.CreationTime,
		InstanceCount: s.InstanceCount,
		InstanceSize:  s.InstanceSize,
		ServiceType:   s.ServiceType,
		Status:        s.Status,
	}
}

func (o *DataTransferServiceResourceProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DataTransferServiceResourceProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

var _ json.Marshaler = DataTransferServiceResourceProperties{}

func (s DataTransferServiceResourceProperties) MarshalJSON() ([]byte, error) {
	type wrapper DataTransferServiceResourceProperties
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling DataTransferServiceResourceProperties: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling DataTransferServiceResourceProperties: %+v", err)
	}

	decoded["serviceType"] = "DataTransfer"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling DataTransferServiceResourceProperties: %+v", err)
	}

	return encoded, nil
}
