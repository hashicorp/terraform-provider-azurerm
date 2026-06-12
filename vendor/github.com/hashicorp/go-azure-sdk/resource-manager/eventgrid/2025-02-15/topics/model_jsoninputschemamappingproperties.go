package topics

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JsonInputSchemaMappingProperties struct {
	DataVersion *JsonFieldWithDefault `json:"dataVersion,omitempty"`
	EventTime   *JsonField            `json:"eventTime,omitempty"`
	EventType   *JsonFieldWithDefault `json:"eventType,omitempty"`
	Id          *JsonField            `json:"id,omitempty"`
	Subject     *JsonFieldWithDefault `json:"subject,omitempty"`
	Topic       *JsonField            `json:"topic,omitempty"`
}
