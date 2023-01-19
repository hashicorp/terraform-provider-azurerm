package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2ValidationSettings struct {
	CheckCertificateRevocationListOnReceive bool                `json:"checkCertificateRevocationListOnReceive"`
	CheckCertificateRevocationListOnSend    bool                `json:"checkCertificateRevocationListOnSend"`
	CheckDuplicateMessage                   bool                `json:"checkDuplicateMessage"`
	CompressMessage                         bool                `json:"compressMessage"`
	EncryptMessage                          bool                `json:"encryptMessage"`
	EncryptionAlgorithm                     EncryptionAlgorithm `json:"encryptionAlgorithm"`
	InterchangeDuplicatesValidityDays       int64               `json:"interchangeDuplicatesValidityDays"`
	OverrideMessageProperties               bool                `json:"overrideMessageProperties"`
	SignMessage                             bool                `json:"signMessage"`
	SigningAlgorithm                        *SigningAlgorithm   `json:"signingAlgorithm,omitempty"`
}
