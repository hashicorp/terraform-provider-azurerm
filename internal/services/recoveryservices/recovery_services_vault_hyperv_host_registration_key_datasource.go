package recoveryservices

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-10-01/vaultcertificates"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservices/2022-10-01/vaults"
	"github.com/hashicorp/go-azure-sdk/resource-manager/recoveryservicessiterecovery/2022-10-01/replicationfabrics"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Registration Key is a special resource, it only support PUT, not GET/DELETE.
// And it will be invalid with the param `certificateCreateOptions.validityInHours`

const vaultResourceType string = "Vaults"
const vaultProviderNameSpace string = "Microsoft.RecoveryServices"
const xmlContentVersion string = "2.0"

type HyperVHostRegistrationKeyModel struct {
	Name                        string `tfschema:"name"`
	HyperVSiteId                string `tfschema:"site_recovery_services_vault_hyperv_site_id"`
	ValidateInHours             int64  `tfschema:"validate_in_hours"`
	XmlContent                  string `tfschema:"xml_content"`
	ResourceId                  int64  `tfschema:"resource_id"`
	ManagementCert              string `tfschema:"management_cert"`
	AadTenantId                 string `tfschema:"aad_tenant_id"`
	AadAuthority                string `tfschema:"aad_authority"`
	ServicePrincipalClientId    string `tfschema:"service_principal_client_id"`
	AadVaultAudience            string `tfschema:"aad_vault_audience"`
	AadManagementEndpoint       string `tfschema:"aad_management_endpoint"`
	VaultPrivateEndpointEnabled string `tfschema:"vault_private_endpoint_enabled"`
	ValidateToDate              string `tfschema:"validate_to"`
}

type HyperVHostRegistrationKeyDataSource struct{}

var _ sdk.DataSource = HyperVHostRegistrationKeyDataSource{}

func (h HyperVHostRegistrationKeyDataSource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return vaultcertificates.ValidateCertificateID
}

func (h HyperVHostRegistrationKeyDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"site_recovery_services_vault_hyperv_site_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: replicationfabrics.ValidateReplicationFabricID,
		},

		"validate_in_hours": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			Default:      120,
			ValidateFunc: validation.IntAtLeast(0),
		},
	}
}

func (h HyperVHostRegistrationKeyDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"xml_content": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
		"resource_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"management_cert": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aad_tenant_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aad_authority": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"service_principal_client_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aad_vault_audience": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"aad_management_endpoint": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vault_private_endpoint_enabled": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"validate_to": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (h HyperVHostRegistrationKeyDataSource) ModelObject() interface{} {
	return &HyperVHostRegistrationKeyModel{}
}

func (h HyperVHostRegistrationKeyDataSource) ResourceType() string {
	return "azurerm_recovery_services_vault_hyperv_host_registration_key"
}

func (h HyperVHostRegistrationKeyDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var metaModel HyperVHostRegistrationKeyModel
			if err := metadata.Decode(&metaModel); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			FabricId, err := replicationfabrics.ParseReplicationFabricID(metaModel.HyperVSiteId)
			if err != nil {
				return fmt.Errorf("parsing: %+v", err)
			}

			date := time.Now().Format("2-01-2006")
			name := fmt.Sprintf(`CN=CB_%s-%s-vaultcredentials`, FabricId.VaultName, date)

			id := vaultcertificates.NewCertificateID(FabricId.SubscriptionId, FabricId.ResourceGroupName, FabricId.VaultName, name)

			metadata.SetID(id)

			return resourceRecoveryServicesVaultHyperVHostRegistrationKeyCreateInternal(ctx, id, metadata)
		},
	}
}

func resourceRecoveryServicesVaultHyperVHostRegistrationKeyCreateInternal(ctx context.Context, id vaultcertificates.CertificateId, metadata sdk.ResourceMetaData) error {
	var metaModel HyperVHostRegistrationKeyModel
	if err := metadata.Decode(&metaModel); err != nil {
		return fmt.Errorf("decoding: %+v", err)
	}

	client := metadata.Client.RecoveryServices.VaultCertificatesClient

	FabricId, err := replicationfabrics.ParseReplicationFabricID(metaModel.HyperVSiteId)
	if err != nil {
		return fmt.Errorf("parsing: %+v", err)
	}

	input := azuresdkhacks.CertificateRequest{
		CreateOptions: &azuresdkhacks.CertificateCreateOptions{
			ValidityInHours: metaModel.ValidateInHours,
		},
	}

	resp, err := client.Create(ctx, id, input)
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

	state := HyperVHostRegistrationKeyModel{
		Name:                     id.CertificateName,
		HyperVSiteId:             FabricId.ID(),
		AadTenantId:              detail.AadTenantId,
		AadAuthority:             detail.AadAuthority,
		ServicePrincipalClientId: detail.ServicePrincipalClientId,
		AadManagementEndpoint:    detail.AzureManagementEndpointAudience,
		ValidateInHours:          metaModel.ValidateInHours,
	}

	if detail.ResourceId != nil {
		state.ResourceId = *detail.ResourceId
	}

	if detail.Certificate != nil {
		state.ManagementCert = *detail.Certificate
	}

	if detail.AadAudience != nil {
		state.AadVaultAudience = *detail.AadAudience
	}

	if detail.ValidTo != nil {
		state.ValidateToDate = *detail.ValidTo
	}

	vaultId := vaults.NewVaultID(id.SubscriptionId, id.ResourceGroupName, id.VaultName)

	vault, err := fetchRecoveryServicesVaultDetails(ctx, metadata.Client.RecoveryServices.VaultsClient, vaultId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if vault.Properties.PrivateEndpointStateForSiteRecovery != nil {
		state.VaultPrivateEndpointEnabled = string(*vault.Properties.PrivateEndpointStateForSiteRecovery)
	}

	fabric, err := fetchRecoveryServicesFabricDetails(ctx, metadata.Client.RecoveryServices.FabricClient, *FabricId)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	state.XmlContent, err = flattenHostRegistrationKeyXMLModel(id, vault, fabric, detail)
	if err != nil {
		return fmt.Errorf("flattening %s: %+v", id, err)
	}

	return metadata.Encode(&state)
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
