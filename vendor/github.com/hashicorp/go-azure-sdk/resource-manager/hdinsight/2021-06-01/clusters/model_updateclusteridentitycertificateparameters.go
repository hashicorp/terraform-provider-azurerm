package clusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateClusterIdentityCertificateParameters struct {
	ApplicationId       *string `json:"applicationId,omitempty"`
	Certificate         *string `json:"certificate,omitempty"`
	CertificatePassword *string `json:"certificatePassword,omitempty"`
}
