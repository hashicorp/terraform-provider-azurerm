package inputs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IoTHubStreamInputDataSourceProperties struct {
	ConsumerGroupName      *string `json:"consumerGroupName,omitempty"`
	Endpoint               *string `json:"endpoint,omitempty"`
	IotHubNamespace        *string `json:"iotHubNamespace,omitempty"`
	SharedAccessPolicyKey  *string `json:"sharedAccessPolicyKey,omitempty"`
	SharedAccessPolicyName *string `json:"sharedAccessPolicyName,omitempty"`
}
