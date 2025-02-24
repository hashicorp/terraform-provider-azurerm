package fileservice

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SmbSetting struct {
	AuthenticationMethods    *string       `json:"authenticationMethods,omitempty"`
	ChannelEncryption        *string       `json:"channelEncryption,omitempty"`
	KerberosTicketEncryption *string       `json:"kerberosTicketEncryption,omitempty"`
	Multichannel             *Multichannel `json:"multichannel,omitempty"`
	Versions                 *string       `json:"versions,omitempty"`
}
