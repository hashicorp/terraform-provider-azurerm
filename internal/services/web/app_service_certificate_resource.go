// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2023-12-01/certificates"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultSuppress "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCertificateCreate,
		Read:   resourceAppServiceCertificateRead,
		Update: resourceAppServiceCertificateUpdate,
		Delete: resourceAppServiceCertificateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := certificates.ParseCertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"pfx_blob": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsBase64,
			},

			"password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"key_vault_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateKeyVaultID,
				RequiredWith: []string{"key_vault_secret_id"},
			},

			"key_vault_secret_id": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: keyVaultSuppress.DiffSuppressIgnoreKeyVaultKeyVersion,
				ValidateFunc:     keyvault.ValidateNestedItemID(keyvault.VersionTypeAny, keyvault.NestedItemTypeAny),
				ConflictsWith:    []string{"pfx_blob", "password"},
				ExactlyOneOf:     []string{"key_vault_secret_id", "pfx_blob"},
			},

			"app_service_plan_id": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"friendly_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"subject_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"host_names": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"issuer": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"issue_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"expiration_date": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"thumbprint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"hosting_environment_profile_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceAppServiceCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).Web.CertificatesClient

	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := certificates.NewCertificateID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if !response.WasNotFound(existing.HttpResponse) {
		if err != nil {
			return fmt.Errorf("checking for presence of existing %s: %w", id, err)
		}
		return tf.ImportAsExistsError("azurerm_app_service_certificate", id.ID())
	}

	certificate := certificates.Certificate{
		Properties: &certificates.CertificateProperties{
			Password: pointer.To(d.Get("password").(string)),
		},
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{})),
	}

	if v := d.Get("app_service_plan_id").(string); v != "" {
		certificate.Properties.ServerFarmId = &v
	}

	if v := d.Get("pfx_blob").(string); v != "" {
		certificate.Properties.PfxBlob = pointer.To(v)
	}

	if v := d.Get("key_vault_secret_id").(string); v != "" {
		parsedSecretId, err := keyvault.ParseNestedItemID(v, keyvault.VersionTypeAny, keyvault.NestedItemTypeAny)
		if err != nil {
			return err
		}

		var keyVaultId *string
		if v := d.Get("key_vault_id").(string); v != "" {
			keyVaultId = pointer.To(v)
		} else {
			keyVaultBaseUrl := parsedSecretId.KeyVaultBaseURL

			subscriptionResourceId := commonids.NewSubscriptionID(id.SubscriptionId)
			keyVaultId, err = keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, subscriptionResourceId, keyVaultBaseUrl)
			if err != nil {
				return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %s", keyVaultBaseUrl, err)
			}
			if keyVaultId == nil {
				return fmt.Errorf("unable to determine the Resource ID for the Key Vault at URL %q", keyVaultBaseUrl)
			}
		}

		certificate.Properties.KeyVaultId = keyVaultId
		certificate.Properties.KeyVaultSecretName = pointer.To(parsedSecretId.Name)
	}

	if _, err := client.CreateOrUpdate(ctx, id, certificate); err != nil {
		return fmt.Errorf("creating %s: %s", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateRead(d, meta)
}

func resourceAppServiceCertificateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient

	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := certificates.NewCertificateID(meta.(*clients.Client).Account.SubscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %w", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	existing.Model.Tags = utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{}))

	if _, err := client.CreateOrUpdate(ctx, id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %s", id, err)
	}

	return resourceAppServiceCertificateRead(d, meta)
}

func resourceAppServiceCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(resp.Model.Location))

		if props := model.Properties; props != nil {
			d.Set("friendly_name", props.FriendlyName)
			d.Set("subject_name", props.SubjectName)
			d.Set("host_names", props.HostNames)
			d.Set("issuer", props.Issuer)
			d.Set("tags", model.Tags)
			d.Set("issue_date", props.IssueDate)
			d.Set("expiration_date", props.ExpirationDate)
			d.Set("thumbprint", props.Thumbprint)

			if props.HostingEnvironmentProfile != nil && props.HostingEnvironmentProfile.Id != nil {
				envId, err := commonids.ParseAppServiceEnvironmentID(*props.HostingEnvironmentProfile.Id)
				if err != nil {
					return err
				}
				d.Set("hosting_environment_profile_id", envId.ID())
			}

			if props.ServerFarmId != nil {
				sfID, err := commonids.ParseAppServicePlanID(*props.ServerFarmId)
				if err != nil {
					return err
				}
				d.Set("app_service_plan_id", sfID.ID())
			}
		}
	}

	return nil
}

func resourceAppServiceCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificates.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %w", id, err)
	}

	return nil
}
