package v2025_02_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/attachednetworkconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/catalogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/checknameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/checkscopednameavailability"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/customizationtasks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devboxdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/devcenters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmentdefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/environmenttypes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/galleries"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/imagedefinitions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/images"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/imageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnection"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/networkconnections"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/pools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectcatalogs"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projects"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/projectskus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/schedules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devcenter/2025-02-01/usages"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AttachedNetworkConnections  *attachednetworkconnections.AttachedNetworkConnectionsClient
	Catalogs                    *catalogs.CatalogsClient
	CheckNameAvailability       *checknameavailability.CheckNameAvailabilityClient
	CheckScopedNameAvailability *checkscopednameavailability.CheckScopedNameAvailabilityClient
	CustomizationTasks          *customizationtasks.CustomizationTasksClient
	DevBoxDefinitions           *devboxdefinitions.DevBoxDefinitionsClient
	DevCenters                  *devcenters.DevCentersClient
	EnvironmentDefinitions      *environmentdefinitions.EnvironmentDefinitionsClient
	EnvironmentTypes            *environmenttypes.EnvironmentTypesClient
	Galleries                   *galleries.GalleriesClient
	ImageDefinitions            *imagedefinitions.ImageDefinitionsClient
	ImageVersions               *imageversions.ImageVersionsClient
	Images                      *images.ImagesClient
	NetworkConnection           *networkconnection.NetworkConnectionClient
	NetworkConnections          *networkconnections.NetworkConnectionsClient
	Pools                       *pools.PoolsClient
	ProjectCatalogs             *projectcatalogs.ProjectCatalogsClient
	ProjectPolicies             *projectpolicies.ProjectPoliciesClient
	ProjectSKUs                 *projectskus.ProjectSKUsClient
	Projects                    *projects.ProjectsClient
	SKUs                        *skus.SKUsClient
	Schedules                   *schedules.SchedulesClient
	Usages                      *usages.UsagesClient
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

	checkScopedNameAvailabilityClient, err := checkscopednameavailability.NewCheckScopedNameAvailabilityClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CheckScopedNameAvailability client: %+v", err)
	}
	configureFunc(checkScopedNameAvailabilityClient.Client)

	customizationTasksClient, err := customizationtasks.NewCustomizationTasksClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CustomizationTasks client: %+v", err)
	}
	configureFunc(customizationTasksClient.Client)

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

	environmentDefinitionsClient, err := environmentdefinitions.NewEnvironmentDefinitionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building EnvironmentDefinitions client: %+v", err)
	}
	configureFunc(environmentDefinitionsClient.Client)

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

	imageDefinitionsClient, err := imagedefinitions.NewImageDefinitionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ImageDefinitions client: %+v", err)
	}
	configureFunc(imageDefinitionsClient.Client)

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

	projectCatalogsClient, err := projectcatalogs.NewProjectCatalogsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProjectCatalogs client: %+v", err)
	}
	configureFunc(projectCatalogsClient.Client)

	projectPoliciesClient, err := projectpolicies.NewProjectPoliciesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProjectPolicies client: %+v", err)
	}
	configureFunc(projectPoliciesClient.Client)

	projectSKUsClient, err := projectskus.NewProjectSKUsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building ProjectSKUs client: %+v", err)
	}
	configureFunc(projectSKUsClient.Client)

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
		AttachedNetworkConnections:  attachedNetworkConnectionsClient,
		Catalogs:                    catalogsClient,
		CheckNameAvailability:       checkNameAvailabilityClient,
		CheckScopedNameAvailability: checkScopedNameAvailabilityClient,
		CustomizationTasks:          customizationTasksClient,
		DevBoxDefinitions:           devBoxDefinitionsClient,
		DevCenters:                  devCentersClient,
		EnvironmentDefinitions:      environmentDefinitionsClient,
		EnvironmentTypes:            environmentTypesClient,
		Galleries:                   galleriesClient,
		ImageDefinitions:            imageDefinitionsClient,
		ImageVersions:               imageVersionsClient,
		Images:                      imagesClient,
		NetworkConnection:           networkConnectionClient,
		NetworkConnections:          networkConnectionsClient,
		Pools:                       poolsClient,
		ProjectCatalogs:             projectCatalogsClient,
		ProjectPolicies:             projectPoliciesClient,
		ProjectSKUs:                 projectSKUsClient,
		Projects:                    projectsClient,
		SKUs:                        sKUsClient,
		Schedules:                   schedulesClient,
		Usages:                      usagesClient,
	}, nil
}
