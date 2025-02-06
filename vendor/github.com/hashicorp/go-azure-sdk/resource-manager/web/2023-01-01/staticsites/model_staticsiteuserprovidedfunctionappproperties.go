package staticsites

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteUserProvidedFunctionAppProperties struct {
	CreatedOn             *string `json:"createdOn,omitempty"`
	FunctionAppRegion     *string `json:"functionAppRegion,omitempty"`
	FunctionAppResourceId *string `json:"functionAppResourceId,omitempty"`
}

func (o *StaticSiteUserProvidedFunctionAppProperties) GetCreatedOnAsTime() (*time.Time, error) {
	if o.CreatedOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedOn, "2006-01-02T15:04:05Z07:00")
}

func (o *StaticSiteUserProvidedFunctionAppProperties) SetCreatedOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedOn = &formatted
}
