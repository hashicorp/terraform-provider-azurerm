package cosmosdb

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountKeyMetadata struct {
	GenerationTime *string `json:"generationTime,omitempty"`
}

func (o *AccountKeyMetadata) GetGenerationTimeAsTime() (*time.Time, error) {
	if o.GenerationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.GenerationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AccountKeyMetadata) SetGenerationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.GenerationTime = &formatted
}
