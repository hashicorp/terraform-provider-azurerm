// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package recoveryservices_test

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-01-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2024-04-01/vaultcertificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2024-04-01/replicationfabrics"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

const vaultResourceType string = "Vaults"
const vaultProviderNameSpace string = "Microsoft.RecoveryServices"
const xmlContentVersion string = "2.0"

func (HyperVHostTestResource) generateHyperVHostRegistrationCert(callbackFunc func(xmlContent string) error) func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
	return func(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) error {
		client := clients.RecoveryServices.VaultCertificatesClient

		FabricId, err := replicationfabrics.ParseReplicationFabricID(state.ID)
		if err != nil {
			return fmt.Errorf("parsing: %+v", err)
		}

		input := azuresdkhacks.CertificateRequest{
			CreateOptions: &azuresdkhacks.CertificateCreateOptions{
				ValidityInHours: 120,
			},
		}

		ctx2, cancel := context.WithTimeout(ctx, 15*time.Minute)
		defer cancel()

		date := time.Now().Format("2-01-2006")
		name := fmt.Sprintf(`CN=CB_%s-%s-vaultcredentials`, FabricId.VaultName, date)
		id := vaultcertificates.NewCertificateID(FabricId.SubscriptionId, FabricId.ResourceGroupName, FabricId.VaultName, name)

		resp, err := client.Create(ctx2, id, input)
		if err != nil {
			return fmt.Errorf("creating: %+v", err)
		}

		if resp.Model == nil {
			return fmt.Errorf("retrieving %s: Model was nil", id)
		}

		detail, ok := resp.Model.Properties.(vaultcertificates.ResourceCertificateAndAadDetails)
		if !ok {
			return fmt.Errorf("unexpected response type")
		}

		vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

		vault, err := fetchRecoveryServicesVaultDetails(ctx2, clients.RecoveryServices.VaultsClient, vaultId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		fabric, err := fetchRecoveryServicesFabricDetails(ctx2, clients.RecoveryServices.FabricClient, *FabricId)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		xmlContent, err := flattenHostRegistrationKeyXMLModel(id, vault, fabric, detail)
		if err != nil {
			return fmt.Errorf("flattening %s: %+v", id, err)
		}

		return callbackFunc(xmlContent)
	}
}

func fetchRecoveryServicesVaultDetails(ctx context.Context, client *vaults.VaultsClient, id vaults.VaultId) (*vaults.Vault, error) {
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resp.Model, nil
}

func fetchRecoveryServicesFabricDetails(ctx context.Context, client *replicationfabrics.ReplicationFabricsClient, id replicationfabrics.ReplicationFabricId) (*replicationfabrics.Fabric, error) {
	resp, err := client.Get(ctx, id, replicationfabrics.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return resp.Model, nil
}

const xmlHeader string = `<?xml version="1.0" encoding="utf-8"?>`

type hyperVHostRegistrationKeyXMLModel struct {
	XMLName                             xml.Name             `xml:"RSVaultAsrCreds"`
	XMLNS                               string               `xml:"xmlns,attr"`
	XMLNSI                              string               `xml:"xmlns:i,attr"`
	VaultDetails                        vaultDetailsXMLModel `xml:"VaultDetails"`
	ManagementCert                      string               `xml:"ManagementCert"`
	Version                             string               `xml:"Version"`
	AadDetails                          aadDetailsXMLModel   `xml:"AadDetails"`
	ChannelIntegrityKey                 string               `xml:"ChannelIntegrityKey"`
	SiteId                              string               `xml:"SiteId"`
	SiteName                            string               `xml:"SiteName"`
	PrivateEndpointStateForSiteRecovery string               `xml:"PrivateEndpointStateForSiteRecovery"`
}

type vaultDetailsXMLModel struct {
	SubscriptionId    string `xml:"SubscriptionId"`
	ResourceGroup     string `xml:"ResourceGroup"`
	ResourceName      string `xml:"ResourceName"`
	ResourceId        int64  `xml:"ResourceId"`
	Location          string `xml:"Location"`
	ResourceType      string `xml:"ResourceType"`
	ProviderNamespace string `xml:"ProviderNamespace"`
}

type aadDetailsXMLModel struct {
	AadAuthority             string `xml:"AadAuthority"`
	AadTenantId              string `xml:"AadTenantId"`
	ServicePrincipalClientId string `xml:"ServicePrincipalClientId"`
	AadVaultAudience         string `xml:"AadVaultAudience"`
	ArmManagementEndpoint    string `xml:"ArmManagementEndpoint"`
}

func flattenHostRegistrationKeyXMLModel(id vaultcertificates.CertificateId, vaultDetail *vaults.Vault, fabricDetails *replicationfabrics.Fabric, certDetails vaultcertificates.ResourceCertificateAndAadDetails) (string, error) {
	if vaultDetail == nil {
		return "", fmt.Errorf("vault was nil")
	}

	if fabricDetails == nil {
		return "", fmt.Errorf("fabric was nil")
	}

	xmlModel := hyperVHostRegistrationKeyXMLModel{
		XMLNSI: "http://www.w3.org/2001/XMLSchema-instance",
		XMLNS:  "http://schemas.datacontract.org/2004/07/Microsoft.Azure.Portal.RecoveryServices.Models.Common",
		VaultDetails: vaultDetailsXMLModel{
			SubscriptionId:    id.SubscriptionId,
			ResourceGroup:     id.ResourceGroupName,
			ResourceName:      id.VaultName,
			Location:          vaultDetail.Location,
			ResourceType:      vaultResourceType,
			ProviderNamespace: vaultProviderNameSpace,
		},
		ManagementCert: *certDetails.Certificate,
		Version:        xmlContentVersion,
		AadDetails: aadDetailsXMLModel{
			AadAuthority:             certDetails.AadAuthority,
			AadTenantId:              certDetails.AadTenantId,
			ServicePrincipalClientId: certDetails.ServicePrincipalClientId,
			ArmManagementEndpoint:    certDetails.AzureManagementEndpointAudience,
		},
		ChannelIntegrityKey: "",
	}

	if certDetails.ResourceId != nil {
		xmlModel.VaultDetails.ResourceId = *certDetails.ResourceId
	}

	if certDetails.AadAudience != nil {
		xmlModel.AadDetails.AadVaultAudience = *certDetails.AadAudience
	}

	if fabricDetails.Properties.InternalIdentifier != nil {
		xmlModel.SiteId = *fabricDetails.Properties.InternalIdentifier
	}

	if fabricDetails.Name != nil {
		xmlModel.SiteName = *fabricDetails.Name
	}

	if vaultDetail.Properties.PrivateEndpointStateForSiteRecovery != nil {
		xmlModel.PrivateEndpointStateForSiteRecovery = string(*vaultDetail.Properties.PrivateEndpointStateForSiteRecovery)
	}

	xmlResult, err := xml.Marshal(xmlModel)
	if err != nil {
		return "", fmt.Errorf("marshaling: %+v", err)
	}

	return xmlHeader + string(xmlResult), nil
}
