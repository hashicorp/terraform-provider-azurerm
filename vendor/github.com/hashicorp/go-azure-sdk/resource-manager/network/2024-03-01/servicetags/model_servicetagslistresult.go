package servicetags

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceTagsListResult struct {
	ChangeNumber *string                  `json:"changeNumber,omitempty"`
	Cloud        *string                  `json:"cloud,omitempty"`
	Id           *string                  `json:"id,omitempty"`
	Name         *string                  `json:"name,omitempty"`
	NextLink     *string                  `json:"nextLink,omitempty"`
	Type         *string                  `json:"type,omitempty"`
	Values       *[]ServiceTagInformation `json:"values,omitempty"`
}
