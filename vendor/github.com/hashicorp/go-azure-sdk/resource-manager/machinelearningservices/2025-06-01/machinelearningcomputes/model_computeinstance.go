package machinelearningcomputes

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ Compute = ComputeInstance{}

type ComputeInstance struct {
	Properties *ComputeInstanceProperties `json:"properties,omitempty"`

	// Fields inherited from Compute

	ComputeLocation    *string            `json:"computeLocation,omitempty"`
	ComputeType        ComputeType        `json:"computeType"`
	CreatedOn          *string            `json:"createdOn,omitempty"`
	Description        *string            `json:"description,omitempty"`
	DisableLocalAuth   *bool              `json:"disableLocalAuth,omitempty"`
	IsAttachedCompute  *bool              `json:"isAttachedCompute,omitempty"`
	ModifiedOn         *string            `json:"modifiedOn,omitempty"`
	ProvisioningErrors *[]ErrorResponse   `json:"provisioningErrors,omitempty"`
	ProvisioningState  *ProvisioningState `json:"provisioningState,omitempty"`
	ResourceId         *string            `json:"resourceId,omitempty"`
}

func (s ComputeInstance) Compute() BaseComputeImpl {
	return BaseComputeImpl{
		ComputeLocation:    s.ComputeLocation,
		ComputeType:        s.ComputeType,
		CreatedOn:          s.CreatedOn,
		Description:        s.Description,
		DisableLocalAuth:   s.DisableLocalAuth,
		IsAttachedCompute:  s.IsAttachedCompute,
		ModifiedOn:         s.ModifiedOn,
		ProvisioningErrors: s.ProvisioningErrors,
		ProvisioningState:  s.ProvisioningState,
		ResourceId:         s.ResourceId,
	}
}

func (o *ComputeInstance) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *ComputeInstance) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}

func (o *ComputeInstance) GetModifiedOnAsTime() (*time.Time, error) {
	if o.ModifiedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ModifiedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *ComputeInstance) SetModifiedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ModifiedOn = &formatted
}

var _ json.Marshaler = ComputeInstance{}

func (s ComputeInstance) MarshalJSON() ([]byte, error) {
	type wrapper ComputeInstance
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling ComputeInstance: %+v", err)
	}

	var decoded map[string]interface{}
	if err = json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling ComputeInstance: %+v", err)
	}

	decoded["computeType"] = "ComputeInstance"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling ComputeInstance: %+v", err)
	}

	return encoded, nil
}
