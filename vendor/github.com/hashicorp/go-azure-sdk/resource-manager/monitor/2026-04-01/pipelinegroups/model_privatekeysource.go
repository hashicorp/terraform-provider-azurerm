package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateKeySource struct {
	Location    string               `json:"location"`
	SubLocation string               `json:"subLocation"`
	Type        PrivateKeySourceType `json:"type"`
}
