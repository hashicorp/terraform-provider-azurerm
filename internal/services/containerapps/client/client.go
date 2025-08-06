// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/certificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerapps"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/containerappsrevisions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/daprcomponents"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/jobs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-01-01/managedenvironmentsstorages"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
)

type Client struct {
	CertificatesClient         *certificates.CertificatesClient
	ContainerAppClient         *containerapps.ContainerAppsClient
	ContainerAppRevisionClient *containerappsrevisions.ContainerAppsRevisionsClient
	DaprComponentsClient       *daprcomponents.DaprComponentsClient
	ManagedEnvironmentClient   *managedenvironments.ManagedEnvironmentsClient
	StorageClient              *managedenvironmentsstorages.ManagedEnvironmentsStoragesClient
	JobClient                  *jobs.JobsClient
}

func NewClient(o *common.ClientOptions) (*Client, error) {
	certificatesClient, err := certificates.NewCertificatesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Certificates client : %+v", err)
	}
	o.Configure(certificatesClient.Client, o.Authorizers.ResourceManager)

	containerAppsClient, err := containerapps.NewContainerAppsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Container Apps client : %+v", err)
	}
	o.Configure(containerAppsClient.Client, o.Authorizers.ResourceManager)

	containerAppsRevisionsClient, err := containerappsrevisions.NewContainerAppsRevisionsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Container Apps Revisions client : %+v", err)
	}
	o.Configure(containerAppsRevisionsClient.Client, o.Authorizers.ResourceManager)

	managedEnvironmentClient, err := managedenvironments.NewManagedEnvironmentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Environment client : %+v", err)
	}
	o.Configure(managedEnvironmentClient.Client, o.Authorizers.ResourceManager)

	managedEnvironmentStoragesClient, err := managedenvironmentsstorages.NewManagedEnvironmentsStoragesClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Managed Environment Storages client : %+v", err)
	}
	o.Configure(managedEnvironmentStoragesClient.Client, o.Authorizers.ResourceManager)

	daprComponentClient, err := daprcomponents.NewDaprComponentsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Dapr Components client : %+v", err)
	}
	o.Configure(daprComponentClient.Client, o.Authorizers.ResourceManager)

	jobsClient, err := jobs.NewJobsClientWithBaseURI(o.Environment.ResourceManager)
	if err != nil {
		return nil, fmt.Errorf("building Jobs client : %+v", err)
	}
	o.Configure(jobsClient.Client, o.Authorizers.ResourceManager)

	return &Client{
		CertificatesClient:         certificatesClient,
		ContainerAppClient:         containerAppsClient,
		ContainerAppRevisionClient: containerAppsRevisionsClient,
		DaprComponentsClient:       daprComponentClient,
		ManagedEnvironmentClient:   managedEnvironmentClient,
		StorageClient:              managedEnvironmentStoragesClient,
		JobClient:                  jobsClient,
	}, nil
}
