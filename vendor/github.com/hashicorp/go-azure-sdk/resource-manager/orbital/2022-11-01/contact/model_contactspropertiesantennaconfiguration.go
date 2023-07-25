package contact

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContactsPropertiesAntennaConfiguration struct {
	DestinationIP *string   `json:"destinationIp,omitempty"`
	SourceIPs     *[]string `json:"sourceIps,omitempty"`
}
