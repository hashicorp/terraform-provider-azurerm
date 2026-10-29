package hcpopenshiftclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EtcdDataEncryptionProfile struct {
	CustomerManaged   *CustomerManagedEncryptionProfile       `json:"customerManaged,omitempty"`
	KeyManagementMode EtcdDataEncryptionKeyManagementModeType `json:"keyManagementMode"`
}
