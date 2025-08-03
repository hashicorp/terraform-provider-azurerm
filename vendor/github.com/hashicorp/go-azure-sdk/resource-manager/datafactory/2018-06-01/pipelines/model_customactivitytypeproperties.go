package pipelines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomActivityTypeProperties struct {
	AutoUserSpecification *interface{}                   `json:"autoUserSpecification,omitempty"`
	Command               interface{}                    `json:"command"`
	ExtendedProperties    *map[string]interface{}        `json:"extendedProperties,omitempty"`
	FolderPath            *interface{}                   `json:"folderPath,omitempty"`
	ReferenceObjects      *CustomActivityReferenceObject `json:"referenceObjects,omitempty"`
	ResourceLinkedService *LinkedServiceReference        `json:"resourceLinkedService,omitempty"`
	RetentionTimeInDays   *float64                       `json:"retentionTimeInDays,omitempty"`
}
