package v2026_03_18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/communicationservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/domains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/emailservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/senderusernames"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/smtpusernames"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2026-03-18/suppressionlists"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	CommunicationServices *communicationservices.CommunicationServicesClient
	Domains               *domains.DomainsClient
	EmailServices         *emailservices.EmailServicesClient
	SenderUsernames       *senderusernames.SenderUsernamesClient
	SmtpUsernames         *smtpusernames.SmtpUsernamesClient
	SuppressionLists      *suppressionlists.SuppressionListsClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	communicationServicesClient, err := communicationservices.NewCommunicationServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CommunicationServices client: %+v", err)
	}
	configureFunc(communicationServicesClient.Client)

	domainsClient, err := domains.NewDomainsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Domains client: %+v", err)
	}
	configureFunc(domainsClient.Client)

	emailServicesClient, err := emailservices.NewEmailServicesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EmailServices client: %+v", err)
	}
	configureFunc(emailServicesClient.Client)

	senderUsernamesClient, err := senderusernames.NewSenderUsernamesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SenderUsernames client: %+v", err)
	}
	configureFunc(senderUsernamesClient.Client)

	smtpUsernamesClient, err := smtpusernames.NewSmtpUsernamesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SmtpUsernames client: %+v", err)
	}
	configureFunc(smtpUsernamesClient.Client)

	suppressionListsClient, err := suppressionlists.NewSuppressionListsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SuppressionLists client: %+v", err)
	}
	configureFunc(suppressionListsClient.Client)

	return &Client{
		CommunicationServices: communicationServicesClient,
		Domains:               domainsClient,
		EmailServices:         emailServicesClient,
		SenderUsernames:       senderUsernamesClient,
		SmtpUsernames:         smtpUsernamesClient,
		SuppressionLists:      suppressionListsClient,
	}, nil
}
