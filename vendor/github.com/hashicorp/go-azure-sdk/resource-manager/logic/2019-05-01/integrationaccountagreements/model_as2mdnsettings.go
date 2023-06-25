package integrationaccountagreements

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AS2MdnSettings struct {
	DispositionNotificationTo  *string          `json:"dispositionNotificationTo,omitempty"`
	MdnText                    *string          `json:"mdnText,omitempty"`
	MicHashingAlgorithm        HashingAlgorithm `json:"micHashingAlgorithm"`
	NeedMDN                    bool             `json:"needMDN"`
	ReceiptDeliveryUrl         *string          `json:"receiptDeliveryUrl,omitempty"`
	SendInboundMDNToMessageBox bool             `json:"sendInboundMDNToMessageBox"`
	SendMDNAsynchronously      bool             `json:"sendMDNAsynchronously"`
	SignMDN                    bool             `json:"signMDN"`
	SignOutboundMDNIfOptional  bool             `json:"signOutboundMDNIfOptional"`
}
