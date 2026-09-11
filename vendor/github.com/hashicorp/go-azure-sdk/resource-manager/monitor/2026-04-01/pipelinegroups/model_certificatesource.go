package pipelinegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateSource struct {
	Location    string                `json:"location"`
	SubLocation string                `json:"subLocation"`
	Type        CertificateSourceType `json:"type"`
}
