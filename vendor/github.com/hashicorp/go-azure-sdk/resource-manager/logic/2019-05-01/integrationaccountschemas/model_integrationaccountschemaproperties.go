package integrationaccountschemas

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IntegrationAccountSchemaProperties struct {
	ChangedTime     *string      `json:"changedTime,omitempty"`
	Content         *string      `json:"content,omitempty"`
	ContentLink     *ContentLink `json:"contentLink,omitempty"`
	ContentType     *string      `json:"contentType,omitempty"`
	CreatedTime     *string      `json:"createdTime,omitempty"`
	DocumentName    *string      `json:"documentName,omitempty"`
	FileName        *string      `json:"fileName,omitempty"`
	Metadata        *interface{} `json:"metadata,omitempty"`
	SchemaType      SchemaType   `json:"schemaType"`
	TargetNamespace *string      `json:"targetNamespace,omitempty"`
}

func (o *IntegrationAccountSchemaProperties) GetChangedTimeAsTime() (*time.Time, error) {
	if o.ChangedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ChangedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountSchemaProperties) SetChangedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ChangedTime = &formatted
}

func (o *IntegrationAccountSchemaProperties) GetCreatedTimeAsTime() (*time.Time, error) {
	if o.CreatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IntegrationAccountSchemaProperties) SetCreatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedTime = &formatted
}
