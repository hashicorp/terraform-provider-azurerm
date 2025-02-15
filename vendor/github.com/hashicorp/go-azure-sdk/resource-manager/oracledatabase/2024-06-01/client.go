package v2024_06_01

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasebackups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasecharactersets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabasenationalcharactersets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabases"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/autonomousdatabaseversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudexadatainfrastructures"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/cloudvmclusters"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbnodes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbservers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dbsystemshapes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dnsprivateviews"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/dnsprivatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/giversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/oraclesubscriptions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/systemversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2024-06-01/virtualnetworkaddresses"
	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

type Client struct {
	AutonomousDatabaseBackups               *autonomousdatabasebackups.AutonomousDatabaseBackupsClient
	AutonomousDatabaseCharacterSets         *autonomousdatabasecharactersets.AutonomousDatabaseCharacterSetsClient
	AutonomousDatabaseNationalCharacterSets *autonomousdatabasenationalcharactersets.AutonomousDatabaseNationalCharacterSetsClient
	AutonomousDatabaseVersions              *autonomousdatabaseversions.AutonomousDatabaseVersionsClient
	AutonomousDatabases                     *autonomousdatabases.AutonomousDatabasesClient
	CloudExadataInfrastructures             *cloudexadatainfrastructures.CloudExadataInfrastructuresClient
	CloudVMClusters                         *cloudvmclusters.CloudVMClustersClient
	DbNodes                                 *dbnodes.DbNodesClient
	DbServers                               *dbservers.DbServersClient
	DbSystemShapes                          *dbsystemshapes.DbSystemShapesClient
	DnsPrivateViews                         *dnsprivateviews.DnsPrivateViewsClient
	DnsPrivateZones                         *dnsprivatezones.DnsPrivateZonesClient
	GiVersions                              *giversions.GiVersionsClient
	OracleSubscriptions                     *oraclesubscriptions.OracleSubscriptionsClient
	SystemVersions                          *systemversions.SystemVersionsClient
	VirtualNetworkAddresses                 *virtualnetworkaddresses.VirtualNetworkAddressesClient
}

func NewClientWithBaseURI(sdkApi sdkEnv.Api, configureFunc func(c *resourcemanager.Client)) (*Client, error) {
	autonomousDatabaseBackupsClient, err := autonomousdatabasebackups.NewAutonomousDatabaseBackupsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutonomousDatabaseBackups client: %+v", err)
	}
	configureFunc(autonomousDatabaseBackupsClient.Client)

	autonomousDatabaseCharacterSetsClient, err := autonomousdatabasecharactersets.NewAutonomousDatabaseCharacterSetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutonomousDatabaseCharacterSets client: %+v", err)
	}
	configureFunc(autonomousDatabaseCharacterSetsClient.Client)

	autonomousDatabaseNationalCharacterSetsClient, err := autonomousdatabasenationalcharactersets.NewAutonomousDatabaseNationalCharacterSetsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutonomousDatabaseNationalCharacterSets client: %+v", err)
	}
	configureFunc(autonomousDatabaseNationalCharacterSetsClient.Client)

	autonomousDatabaseVersionsClient, err := autonomousdatabaseversions.NewAutonomousDatabaseVersionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutonomousDatabaseVersions client: %+v", err)
	}
	configureFunc(autonomousDatabaseVersionsClient.Client)

	autonomousDatabasesClient, err := autonomousdatabases.NewAutonomousDatabasesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building AutonomousDatabases client: %+v", err)
	}
	configureFunc(autonomousDatabasesClient.Client)

	cloudExadataInfrastructuresClient, err := cloudexadatainfrastructures.NewCloudExadataInfrastructuresClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CloudExadataInfrastructures client: %+v", err)
	}
	configureFunc(cloudExadataInfrastructuresClient.Client)

	cloudVMClustersClient, err := cloudvmclusters.NewCloudVMClustersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building CloudVMClusters client: %+v", err)
	}
	configureFunc(cloudVMClustersClient.Client)

	dbNodesClient, err := dbnodes.NewDbNodesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DbNodes client: %+v", err)
	}
	configureFunc(dbNodesClient.Client)

	dbServersClient, err := dbservers.NewDbServersClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DbServers client: %+v", err)
	}
	configureFunc(dbServersClient.Client)

	dbSystemShapesClient, err := dbsystemshapes.NewDbSystemShapesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DbSystemShapes client: %+v", err)
	}
	configureFunc(dbSystemShapesClient.Client)

	dnsPrivateViewsClient, err := dnsprivateviews.NewDnsPrivateViewsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DnsPrivateViews client: %+v", err)
	}
	configureFunc(dnsPrivateViewsClient.Client)

	dnsPrivateZonesClient, err := dnsprivatezones.NewDnsPrivateZonesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building DnsPrivateZones client: %+v", err)
	}
	configureFunc(dnsPrivateZonesClient.Client)

	giVersionsClient, err := giversions.NewGiVersionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building GiVersions client: %+v", err)
	}
	configureFunc(giVersionsClient.Client)

	oracleSubscriptionsClient, err := oraclesubscriptions.NewOracleSubscriptionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building OracleSubscriptions client: %+v", err)
	}
	configureFunc(oracleSubscriptionsClient.Client)

	systemVersionsClient, err := systemversions.NewSystemVersionsClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building SystemVersions client: %+v", err)
	}
	configureFunc(systemVersionsClient.Client)

	virtualNetworkAddressesClient, err := virtualnetworkaddresses.NewVirtualNetworkAddressesClientWithBaseURI(sdkApi)
	if err != nil {
		return nil, fmt.Errorf("building VirtualNetworkAddresses client: %+v", err)
	}
	configureFunc(virtualNetworkAddressesClient.Client)

	return &Client{
		AutonomousDatabaseBackups:               autonomousDatabaseBackupsClient,
		AutonomousDatabaseCharacterSets:         autonomousDatabaseCharacterSetsClient,
		AutonomousDatabaseNationalCharacterSets: autonomousDatabaseNationalCharacterSetsClient,
		AutonomousDatabaseVersions:              autonomousDatabaseVersionsClient,
		AutonomousDatabases:                     autonomousDatabasesClient,
		CloudExadataInfrastructures:             cloudExadataInfrastructuresClient,
		CloudVMClusters:                         cloudVMClustersClient,
		DbNodes:                                 dbNodesClient,
		DbServers:                               dbServersClient,
		DbSystemShapes:                          dbSystemShapesClient,
		DnsPrivateViews:                         dnsPrivateViewsClient,
		DnsPrivateZones:                         dnsPrivateZonesClient,
		GiVersions:                              giVersionsClient,
		OracleSubscriptions:                     oracleSubscriptionsClient,
		SystemVersions:                          systemVersionsClient,
		VirtualNetworkAddresses:                 virtualNetworkAddressesClient,
	}, nil
}
