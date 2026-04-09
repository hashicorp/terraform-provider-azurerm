package streamingjobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomClrSerializationProperties struct {
	SerializationClassName *string `json:"serializationClassName,omitempty"`
	SerializationDllPath   *string `json:"serializationDllPath,omitempty"`
}
