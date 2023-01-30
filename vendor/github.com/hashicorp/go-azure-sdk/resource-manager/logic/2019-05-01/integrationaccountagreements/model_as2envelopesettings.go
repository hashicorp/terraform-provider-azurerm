package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2EnvelopeSettings struct {
	AutogenerateFileName                    bool   `json:"autogenerateFileName"`
	FileNameTemplate                        string `json:"fileNameTemplate"`
	MessageContentType                      string `json:"messageContentType"`
	SuspendMessageOnFileNameGenerationError bool   `json:"suspendMessageOnFileNameGenerationError"`
	TransmitFileNameInMimeHeader            bool   `json:"transmitFileNameInMimeHeader"`
}
