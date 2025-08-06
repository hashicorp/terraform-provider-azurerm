package serviceresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetTdeCertificatesSqlTaskInput struct {
	BackupFileShare      FileShare                  `json:"backupFileShare"`
	ConnectionInfo       SqlConnectionInfo          `json:"connectionInfo"`
	SelectedCertificates []SelectedCertificateInput `json:"selectedCertificates"`
}
