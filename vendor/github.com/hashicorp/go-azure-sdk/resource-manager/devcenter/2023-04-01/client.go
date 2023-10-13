package v2023_04_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/attachednetworkconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/catalogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/checknameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/devboxdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/devcenters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/environmenttypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/imageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/networkconnection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/networkconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/pools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/schedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2023-04-01/usages"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AttachedNetworkConnections *attachednetworkconnections.AttachedNetworkConnectionsClient
	Catalogs                   *catalogs.CatalogsClient
	CheckNameAvailability      *checknameavailability.CheckNameAvailabilityClient
	DevBoxDefinitions          *devboxdefinitions.DevBoxDefinitionsClient
	DevCenters                 *devcenters.DevCentersClient
	EnvironmentTypes           *environmenttypes.EnvironmentTypesClient
	Galleries                  *galleries.GalleriesClient
	ImageVersions              *imageversions.ImageVersionsClient
	Images                     *images.ImagesClient
	NetworkConnection          *networkconnection.NetworkConnectionClient
	NetworkConnections         *networkconnections.NetworkConnectionsClient
	Pools                      *pools.PoolsClient
	Projects                   *projects.ProjectsClient
	SKUs                       *skus.SKUsClient
	Schedules                  *schedules.SchedulesClient
	Usages                     *usages.UsagesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	attachedNetworkConnectionsClient, err := attachednetworkconnections.NewAttachedNetworkConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AttachedNetworkConnections client: %+v", err)
	}
	configureFunc(attachedNetworkConnectionsClient.Client)

	catalogsClient, err := catalogs.NewCatalogsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Catalogs client: %+v", err)
	}
	configureFunc(catalogsClient.Client)

	checkNameAvailabilityClient, err := checknameavailability.NewCheckNameAvailabilityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CheckNameAvailability client: %+v", err)
	}
	configureFunc(checkNameAvailabilityClient.Client)

	devBoxDefinitionsClient, err := devboxdefinitions.NewDevBoxDefinitionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DevBoxDefinitions client: %+v", err)
	}
	configureFunc(devBoxDefinitionsClient.Client)

	devCentersClient, err := devcenters.NewDevCentersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DevCenters client: %+v", err)
	}
	configureFunc(devCentersClient.Client)

	environmentTypesClient, err := environmenttypes.NewEnvironmentTypesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EnvironmentTypes client: %+v", err)
	}
	configureFunc(environmentTypesClient.Client)

	galleriesClient, err := galleries.NewGalleriesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Galleries client: %+v", err)
	}
	configureFunc(galleriesClient.Client)

	imageVersionsClient, err := imageversions.NewImageVersionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ImageVersions client: %+v", err)
	}
	configureFunc(imageVersionsClient.Client)

	imagesClient, err := images.NewImagesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Images client: %+v", err)
	}
	configureFunc(imagesClient.Client)

	networkConnectionClient, err := networkconnection.NewNetworkConnectionClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NetworkConnection client: %+v", err)
	}
	configureFunc(networkConnectionClient.Client)

	networkConnectionsClient, err := networkconnections.NewNetworkConnectionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building NetworkConnections client: %+v", err)
	}
	configureFunc(networkConnectionsClient.Client)

	poolsClient, err := pools.NewPoolsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Pools client: %+v", err)
	}
	configureFunc(poolsClient.Client)

	projectsClient, err := projects.NewProjectsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Projects client: %+v", err)
	}
	configureFunc(projectsClient.Client)

	sKUsClient, err := skus.NewSKUsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SKUs client: %+v", err)
	}
	configureFunc(sKUsClient.Client)

	schedulesClient, err := schedules.NewSchedulesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Schedules client: %+v", err)
	}
	configureFunc(schedulesClient.Client)

	usagesClient, err := usages.NewUsagesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building Usages client: %+v", err)
	}
	configureFunc(usagesClient.Client)

	return &Client{
		AttachedNetworkConnections: attachedNetworkConnectionsClient,
		Catalogs:                   catalogsClient,
		CheckNameAvailability:      checkNameAvailabilityClient,
		DevBoxDefinitions:          devBoxDefinitionsClient,
		DevCenters:                 devCentersClient,
		EnvironmentTypes:           environmentTypesClient,
		Galleries:                  galleriesClient,
		ImageVersions:              imageVersionsClient,
		Images:                     imagesClient,
		NetworkConnection:          networkConnectionClient,
		NetworkConnections:         networkConnectionsClient,
		Pools:                      poolsClient,
		Projects:                   projectsClient,
		SKUs:                       sKUsClient,
		Schedules:                  schedulesClient,
		Usages:                     usagesClient,
	}, nil
}
