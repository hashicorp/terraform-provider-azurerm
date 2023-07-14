// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package domainservices

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DomainServiceTrustResource struct{}

var _ sdk.ResourceWithUpdate = DomainServiceTrustResource{}

type DomainServiceTrustModel struct {
	Name                string   `tfschema:"name"`
	DomainServiceId     string   `tfschema:"domain_service_id"`
	TrustedDomainFqdn   string   `tfschema:"trusted_domain_fqdn"`
	TrustedDomainDnsIPs []string `tfschema:"trusted_domain_dns_ips"`
	Password            string   `tfschema:"password"`
}

func (r DomainServiceTrustResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"domain_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.DomainServiceID,
		},
		"trusted_domain_fqdn": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"trusted_domain_dns_ips": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 2,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPAddress,
			},
		},
		"password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r DomainServiceTrustResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DomainServiceTrustResource) ResourceType() string {
	return "azurerm_active_directory_domain_service_trust"
}

func (r DomainServiceTrustResource) ModelObject() interface{} {
	return &DomainServiceTrustModel{}
}

func (r DomainServiceTrustResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.DomainServiceTrustID
}

func (r DomainServiceTrustResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DomainServices.DomainServicesClient

			var plan DomainServiceTrustModel
			if err := metadata.Decode(&plan); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			dsid, err := parse.DomainServiceID(plan.DomainServiceId)
			if err != nil {
				return err
			}

			id := parse.NewDomainServiceTrustID(dsid.SubscriptionId, dsid.ResourceGroup, dsid.Name, plan.Name)
			idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

			locks.ByName(id.DomainServiceName, DomainServiceResourceName)
			defer locks.UnlockByName(id.DomainServiceName, DomainServiceResourceName)

			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s: returned with null model", idsdk)
			}

			existingTrusts := []domainservices.ForestTrust{}
			if props := model.Properties; props != nil {
				if fsettings := props.ResourceForestSettings; fsettings != nil {
					if settings := fsettings.Settings; settings != nil {
						existingTrusts = *settings
					}
				}
			}
			for _, setting := range existingTrusts {
				if setting.FriendlyName != nil && *setting.FriendlyName == id.TrustName {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			existingTrusts = append(existingTrusts, domainservices.ForestTrust{
				TrustedDomainFqdn: utils.String(plan.TrustedDomainFqdn),
				TrustDirection:    utils.String("Inbound"),
				FriendlyName:      utils.String(id.TrustName),
				RemoteDnsIPs:      utils.String(strings.Join(plan.TrustedDomainDnsIPs, ",")),
				TrustPassword:     utils.String(plan.Password),
			})
			params := domainservices.DomainService{
				Properties: &domainservices.DomainServiceProperties{
					ResourceForestSettings: &domainservices.ResourceForestSettings{
						Settings: &existingTrusts,
					},
				},
			}

			if err := client.UpdateThenPoll(ctx, idsdk, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r DomainServiceTrustResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DomainServices.DomainServicesClient
			id, err := parse.DomainServiceTrustID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", idsdk, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s: returned with null model", idsdk)
			}

			props := model.Properties
			if props == nil {
				return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing properties", idsdk)
			}

			existingTrusts := []domainservices.ForestTrust{}
			if props != nil {
				if fsettings := props.ResourceForestSettings; fsettings != nil {
					if settings := fsettings.Settings; settings != nil {
						existingTrusts = *settings
					}
				}
			}
			var trust *domainservices.ForestTrust
			for _, setting := range existingTrusts {
				existingTrust := setting
				if setting.FriendlyName != nil && *setting.FriendlyName == id.TrustName {
					trust = &existingTrust
				}
			}
			if trust == nil {
				return metadata.MarkAsGone(id)
			}

			// Retrieve the initial replica set id to construct the domain service id.
			replicaSets := flattenDomainServiceReplicaSets(props.ReplicaSets)
			if len(replicaSets) == 0 {
				return fmt.Errorf("checking for presence of existing %s: API response contained nil or missing replica set details", idsdk)
			}
			initialReplicaSetId := replicaSets[0].(map[string]interface{})["id"].(string)
			dsid := parse.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName, initialReplicaSetId)

			var state DomainServiceTrustModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			data := DomainServiceTrustModel{
				DomainServiceId: dsid.ID(),
				Name:            id.TrustName,
				// Setting the password from state as it is not returned by API.
				Password: state.Password,
			}

			if trust.TrustedDomainFqdn != nil {
				data.TrustedDomainFqdn = *trust.TrustedDomainFqdn
			}

			if trust.RemoteDnsIPs != nil {
				data.TrustedDomainDnsIPs = strings.Split(*trust.RemoteDnsIPs, ",")
			}

			return metadata.Encode(&data)
		},
	}
}

func (r DomainServiceTrustResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DomainServices.DomainServicesClient

			id, err := parse.DomainServiceTrustID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

			locks.ByName(id.DomainServiceName, DomainServiceResourceName)
			defer locks.UnlockByName(id.DomainServiceName, DomainServiceResourceName)

			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", idsdk, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s: returned with null model", idsdk)
			}
			existingTrusts := []domainservices.ForestTrust{}
			if props := model.Properties; props != nil {
				if fsettings := props.ResourceForestSettings; fsettings != nil {
					if settings := fsettings.Settings; settings != nil {
						existingTrusts = *settings
					}
				}
			}
			var found bool
			newTrusts := []domainservices.ForestTrust{}
			for _, trust := range existingTrusts {
				if trust.FriendlyName != nil && *trust.FriendlyName == id.TrustName {
					found = true
					continue
				}
				newTrusts = append(newTrusts, trust)
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			params := domainservices.DomainService{
				Properties: &domainservices.DomainServiceProperties{
					ResourceForestSettings: &domainservices.ResourceForestSettings{
						Settings: &newTrusts,
					},
				},
			}

			if err := client.UpdateThenPoll(ctx, idsdk, params); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r DomainServiceTrustResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DomainServices.DomainServicesClient

			id, err := parse.DomainServiceTrustID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			idsdk := domainservices.NewDomainServiceID(id.SubscriptionId, id.ResourceGroup, id.DomainServiceName)

			locks.ByName(id.DomainServiceName, DomainServiceResourceName)
			defer locks.UnlockByName(id.DomainServiceName, DomainServiceResourceName)

			existing, err := client.Get(ctx, idsdk)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", idsdk, err)
			}

			model := existing.Model
			if model == nil {
				return fmt.Errorf("reading %s: returned with null model", idsdk)
			}

			existingTrusts := []domainservices.ForestTrust{}
			if props := model.Properties; props != nil {
				if fsettings := props.ResourceForestSettings; fsettings != nil {
					if settings := fsettings.Settings; settings != nil {
						existingTrusts = *settings
					}
				}
			}

			var plan DomainServiceTrustModel
			if err := metadata.Decode(&plan); err != nil {
				return err
			}

			var found bool
			newTrusts := []domainservices.ForestTrust{}
			for _, trust := range existingTrusts {
				if trust.FriendlyName != nil && *trust.FriendlyName == id.TrustName {
					found = true
					if metadata.ResourceData.HasChange("trusted_domain_fqdn") {
						trust.TrustedDomainFqdn = utils.String(plan.TrustedDomainFqdn)
					}
					if metadata.ResourceData.HasChange("trusted_domain_dns_ips") {
						trust.RemoteDnsIPs = utils.String(strings.Join(plan.TrustedDomainDnsIPs, ","))
					}
					trust.TrustPassword = utils.String(plan.Password)
				}
				newTrusts = append(newTrusts, trust)
			}
			if !found {
				return fmt.Errorf("%s not exists: %+v", id, err)
			}

			params := domainservices.DomainService{
				Properties: &domainservices.DomainServiceProperties{
					ResourceForestSettings: &domainservices.ResourceForestSettings{
						Settings: &newTrusts,
					},
				},
			}

			if err := client.UpdateThenPoll(ctx, idsdk, params); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}
