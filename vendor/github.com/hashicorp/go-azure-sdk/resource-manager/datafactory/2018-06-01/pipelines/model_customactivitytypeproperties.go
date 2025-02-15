package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomActivityTypeProperties struct {
	AutoUserSpecification *string                        `json:"autoUserSpecification,omitempty"`
	Command               string                         `json:"command"`
	ExtendedProperties    *map[string]string             `json:"extendedProperties,omitempty"`
	FolderPath            *string                        `json:"folderPath,omitempty"`
	ReferenceObjects      *CustomActivityReferenceObject `json:"referenceObjects,omitempty"`
	ResourceLinkedService *LinkedServiceReference        `json:"resourceLinkedService,omitempty"`
	RetentionTimeInDays   *float64                       `json:"retentionTimeInDays,omitempty"`
}
