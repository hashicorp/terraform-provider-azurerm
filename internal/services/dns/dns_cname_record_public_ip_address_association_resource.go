// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01/recordsets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DnsCnameRecordPublicIpAddressAssociationResource struct{}

var _ sdk.Resource = DnsCnameRecordPublicIpAddressAssociationResource{}

type DnsCnameRecordPublicIpAddressAssociationModel struct {
	CnameRecordId     string `tfschema:"dns_cname_record_id"`
	PublicIpAddressId string `tfschema:"public_ip_address_id"`
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dns_cname_record_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: recordsets.ValidateRecordTypeID,
		},
		"public_ip_address_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidatePublicIPAddressID,
		},
	}
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) ModelObject() interface{} {
	return &DnsCnameRecordPublicIpAddressAssociationModel{}
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) ResourceType() string {
	return "azurerm_dns_cname_record_public_ip_address_association"
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return parse.DnsCnameRecordPublicIpAddressAssociationIDValidation
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model DnsCnameRecordPublicIpAddressAssociationModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			dnsClient := metadata.Client.Dns.RecordSets
			publicIpClient := metadata.Client.Network.PublicIPAddresses

			cnameRecordId, err := recordsets.ParseRecordTypeID(model.CnameRecordId)
			if err != nil {
				return err
			}

			if cnameRecordId.RecordType != recordsets.RecordTypeCNAME {
				return fmt.Errorf("expected record type to be CNAME but got %s", cnameRecordId.RecordType)
			}

			publicIpAddressId, err := commonids.ParsePublicIPAddressID(model.PublicIpAddressId)
			if err != nil {
				return err
			}

			locks.ByName(cnameRecordId.RelativeRecordSetName, "azurerm_dns_cname_record")
			defer locks.UnlockByName(cnameRecordId.RelativeRecordSetName, "azurerm_dns_cname_record")

			locks.ByName(publicIpAddressId.PublicIPAddressesName, "azurerm_public_ip")
			defer locks.UnlockByName(publicIpAddressId.PublicIPAddressesName, "azurerm_public_ip")

			cnameRecord, err := dnsClient.Get(ctx, *cnameRecordId)
			if err != nil {
				if response.WasNotFound(cnameRecord.HttpResponse) {
					return fmt.Errorf("CNAME Record %s was not found", cnameRecordId)
				}
				return fmt.Errorf("retrieving CNAME Record %s: %+v", cnameRecordId, err)
			}

			if cnameRecord.Model == nil || cnameRecord.Model.Properties == nil {
				return fmt.Errorf("model/properties was nil for CNAME Record %s", cnameRecordId)
			}

			fqdn := pointer.From(cnameRecord.Model.Properties.Fqdn)
			if fqdn == "" {
				return fmt.Errorf("FQDN was empty for CNAME Record %s", cnameRecordId)
			}

			publicIp, err := publicIpClient.Get(ctx, *publicIpAddressId, publicipaddresses.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(publicIp.HttpResponse) {
					return fmt.Errorf("Public IP Address %s was not found", publicIpAddressId)
				}
				return fmt.Errorf("retrieving Public IP Address %s: %+v", publicIpAddressId, err)
			}

			if publicIp.Model == nil || publicIp.Model.Properties == nil {
				return fmt.Errorf("model/properties was nil for Public IP Address %s", publicIpAddressId)
			}

			if publicIp.Model.Properties.DnsSettings != nil && publicIp.Model.Properties.DnsSettings.ReverseFqdn != nil {
				existingReverseFqdn := *publicIp.Model.Properties.DnsSettings.ReverseFqdn
				if existingReverseFqdn != "" && existingReverseFqdn != fqdn {
					return fmt.Errorf("Public IP Address %s already has a reverse_fqdn set to %q, cannot set to %q", publicIpAddressId, existingReverseFqdn, fqdn)
				}
			}

			if publicIp.Model.Properties.DnsSettings == nil {
				publicIp.Model.Properties.DnsSettings = &publicipaddresses.PublicIPAddressDnsSettings{}
			}
			publicIp.Model.Properties.DnsSettings.ReverseFqdn = pointer.To(fqdn)

			if err := publicIpClient.CreateOrUpdateThenPoll(ctx, *publicIpAddressId, *publicIp.Model); err != nil {
				return fmt.Errorf("updating reverse_fqdn for Public IP Address %s: %+v", publicIpAddressId, err)
			}

			resourceId := parse.NewDnsCnameRecordPublicIpAddressAssociationId(*cnameRecordId, *publicIpAddressId)
			metadata.SetID(&resourceId)

			return nil
		},
	}
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			publicIpClient := metadata.Client.Network.PublicIPAddresses
			dnsClient := metadata.Client.Dns.RecordSets

			resourceId, err := parse.DnsCnameRecordPublicIpAddressAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			cnameRecord, err := dnsClient.Get(ctx, resourceId.CnameRecordId)
			if err != nil {
				if response.WasNotFound(cnameRecord.HttpResponse) {
					log.Printf("[DEBUG] CNAME Record %s was not found - removing from state", resourceId.CnameRecordId)
					return metadata.MarkAsGone(&resourceId.CnameRecordId)
				}
				return fmt.Errorf("retrieving CNAME Record %s: %+v", resourceId.CnameRecordId, err)
			}

			publicIp, err := publicIpClient.Get(ctx, resourceId.PublicIpAddressId, publicipaddresses.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(publicIp.HttpResponse) {
					log.Printf("[DEBUG] Public IP Address %s was not found - removing from state", resourceId.PublicIpAddressId)
					return metadata.MarkAsGone(&resourceId.PublicIpAddressId)
				}
				return fmt.Errorf("retrieving Public IP Address %s: %+v", resourceId.PublicIpAddressId, err)
			}

			if publicIp.Model == nil || publicIp.Model.Properties == nil {
				log.Printf("[DEBUG] Public IP Address %s has no properties - removing from state", resourceId.PublicIpAddressId)
				return metadata.MarkAsGone(&resourceId.PublicIpAddressId)
			}

			if publicIp.Model.Properties.DnsSettings == nil || publicIp.Model.Properties.DnsSettings.ReverseFqdn == nil {
				log.Printf("[DEBUG] Public IP Address %s has no reverse_fqdn set - removing from state", resourceId.PublicIpAddressId)
				return metadata.MarkAsGone(&resourceId.PublicIpAddressId)
			}

			if cnameRecord.Model == nil || cnameRecord.Model.Properties == nil {
				log.Printf("[DEBUG] CNAME Record %s has no properties - removing from state", resourceId.CnameRecordId)
				return metadata.MarkAsGone(&resourceId.CnameRecordId)
			}

			expectedFqdn := pointer.From(cnameRecord.Model.Properties.Fqdn)
			actualReverseFqdn := pointer.From(publicIp.Model.Properties.DnsSettings.ReverseFqdn)

			if expectedFqdn != actualReverseFqdn {
				log.Printf("[DEBUG] Public IP Address %s reverse_fqdn (%s) does not match CNAME record FQDN (%s) - removing from state", resourceId.PublicIpAddressId, actualReverseFqdn, expectedFqdn)
				return metadata.MarkAsGone(&resourceId.PublicIpAddressId)
			}

			state := DnsCnameRecordPublicIpAddressAssociationModel{
				CnameRecordId:     resourceId.CnameRecordId.ID(),
				PublicIpAddressId: resourceId.PublicIpAddressId.ID(),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			publicIpClient := metadata.Client.Network.PublicIPAddresses

			resourceId, err := parse.DnsCnameRecordPublicIpAddressAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(resourceId.PublicIpAddressId.PublicIPAddressesName, "azurerm_public_ip")
			defer locks.UnlockByName(resourceId.PublicIpAddressId.PublicIPAddressesName, "azurerm_public_ip")

			publicIp, err := publicIpClient.Get(ctx, resourceId.PublicIpAddressId, publicipaddresses.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(publicIp.HttpResponse) {
					return nil
				}
				return fmt.Errorf("retrieving Public IP Address %s: %+v", resourceId.PublicIpAddressId, err)
			}

			if publicIp.Model == nil || publicIp.Model.Properties == nil {
				return nil
			}

			if publicIp.Model.Properties.DnsSettings != nil {
				publicIp.Model.Properties.DnsSettings.ReverseFqdn = nil
			}

			if err := publicIpClient.CreateOrUpdateThenPoll(ctx, resourceId.PublicIpAddressId, *publicIp.Model); err != nil {
				return fmt.Errorf("removing reverse_fqdn from Public IP Address %s: %+v", resourceId.PublicIpAddressId, err)
			}

			return nil
		},
	}
}
