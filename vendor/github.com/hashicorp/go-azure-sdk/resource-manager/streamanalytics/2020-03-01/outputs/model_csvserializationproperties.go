package outputs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CsvSerializationProperties struct {
	Encoding       *Encoding `json:"encoding,omitempty"`
	FieldDelimiter *string   `json:"fieldDelimiter,omitempty"`
}
