package certificate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateCreateOrUpdateParameters struct {
	Name       string                              `json:"name"`
	Properties CertificateCreateOrUpdateProperties `json:"properties"`
}
