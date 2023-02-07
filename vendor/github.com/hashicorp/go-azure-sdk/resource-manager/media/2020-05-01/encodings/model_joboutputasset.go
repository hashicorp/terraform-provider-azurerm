package encodings

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

var _ JobOutput = JobOutputAsset{}

type JobOutputAsset struct {
	AssetName string `json:"assetName"`

	// Fields inherited from JobOutput
	EndTime   *string   `json:"endTime,omitempty"`
	Error     *JobError `json:"error,omitempty"`
	Label     *string   `json:"label,omitempty"`
	Progress  *int64    `json:"progress,omitempty"`
	StartTime *string   `json:"startTime,omitempty"`
	State     *JobState `json:"state,omitempty"`
}

func (o *JobOutputAsset) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobOutputAsset) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *JobOutputAsset) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobOutputAsset) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}

var _ json.Marshaler = JobOutputAsset{}

func (s JobOutputAsset) MarshalJSON() ([]byte, error) {
	type wrapper JobOutputAsset
	wrapped := wrapper(s)
	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, fmt.Errorf("marshaling JobOutputAsset: %+v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling JobOutputAsset: %+v", err)
	}
	decoded["@odata.type"] = "#Microsoft.Media.JobOutputAsset"

	encoded, err = json.Marshal(decoded)
	if err != nil {
		return nil, fmt.Errorf("re-marshaling JobOutputAsset: %+v", err)
	}

	return encoded, nil
}
