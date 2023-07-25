// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package web

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2021-02-01/web" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	keyVaultParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
	keyVaultSuppress "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/suppress"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/web/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAppServiceCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAppServiceCertificateCreateUpdate,
		Read:   resourceAppServiceCertificateRead,
		Update: resourceAppServiceCertificateCreateUpdate,
		Delete: resourceAppServiceCertificateDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.CertificateID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceAppServiceCertificateSchema(),
	}
}

func resourceAppServiceCertificateCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	keyVaultsClient := meta.(*clients.Client).KeyVault
	client := meta.(*clients.Client).Web.CertificatesClient
	resourcesClient := meta.(*clients.Client).Resource
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for App Service Certificate creation.")

	id := parse.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	location := azure.NormalizeLocation(d.Get("location").(string))
	pfxBlob := d.Get("pfx_blob").(string)
	password := d.Get("password").(string)
	keyVaultSecretId := d.Get("key_vault_secret_id").(string)
	appServicePlanId := d.Get("app_service_plan_id").(string)
	t := d.Get("tags").(map[string]interface{})

	if pfxBlob == "" && keyVaultSecretId == "" {
		return fmt.Errorf("Either `pfx_blob` or `key_vault_secret_id` must be set")
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_app_service_certificate", id.ID())
		}
	}

	certificate := web.Certificate{
		CertificateProperties: &web.CertificateProperties{
			Password: utils.String(password),
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	if appServicePlanId != "" {
		certificate.CertificateProperties.ServerFarmID = &appServicePlanId
	}

	if pfxBlob != "" {
		decodedPfxBlob, err := base64.StdEncoding.DecodeString(pfxBlob)
		if err != nil {
			return fmt.Errorf("Could not decode PFX blob: %+v", err)
		}
		certificate.CertificateProperties.PfxBlob = &decodedPfxBlob
	}

	if keyVaultSecretId != "" {
		parsedSecretId, err := keyVaultParse.ParseNestedItemID(keyVaultSecretId)
		if err != nil {
			return err
		}

		keyVaultBaseUrl := parsedSecretId.KeyVaultBaseUrl

		keyVaultId, err := keyVaultsClient.KeyVaultIDFromBaseUrl(ctx, resourcesClient, keyVaultBaseUrl)
		if err != nil {
			return fmt.Errorf("retrieving the Resource ID for the Key Vault at URL %q: %s", keyVaultBaseUrl, err)
		}
		if keyVaultId == nil {
			return fmt.Errorf("Unable to determine the Resource ID for the Key Vault at URL %q", keyVaultBaseUrl)
		}

		certificate.CertificateProperties.KeyVaultID = keyVaultId
		certificate.CertificateProperties.KeyVaultSecretName = utils.String(parsedSecretId.Name)
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, certificate); err != nil {
		return fmt.Errorf("creating/updating %s: %s", id, err)
	}

	d.SetId(id.ID())

	return resourceAppServiceCertificateRead(d, meta)
}

func resourceAppServiceCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Certificate %q (Resource Group %q) was not found - removing from state", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on App Service Certificate %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.CertificateProperties; props != nil {
		d.Set("friendly_name", props.FriendlyName)
		d.Set("subject_name", props.SubjectName)
		d.Set("host_names", props.HostNames)
		d.Set("issuer", props.Issuer)
		if props.HostingEnvironmentProfile != nil && props.HostingEnvironmentProfile.ID != nil {
			envId, err := parse.AppServiceEnvironmentID(*props.HostingEnvironmentProfile.ID)
			if err != nil {
				return fmt.Errorf("parsing hosting environment error: %+v", err)
			}
			d.Set("hosting_environment_profile_id", envId.ID())
		}
		issueDate := ""
		if props.IssueDate != nil {
			issueDate = props.IssueDate.Format(time.RFC3339)
		}
		d.Set("issue_date", issueDate)
		expirationDate := ""
		if props.ExpirationDate != nil {
			expirationDate = props.ExpirationDate.Format(time.RFC3339)
		}
		d.Set("expiration_date", expirationDate)
		d.Set("thumbprint", props.Thumbprint)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAppServiceCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.CertificatesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.CertificateID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting App Service Certificate %q (Resource Group %q)", id.Name, id.ResourceGroup)

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting App Service Certificate %q (Resource Group %q): %s)", id.Name, id.ResourceGroup, err)
		}
	}

	return nil
}

func resourceAppServiceCertificateSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
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

		"key_vault_secret_id": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: keyVaultSuppress.DiffSuppressIgnoreKeyVaultKeyVersion,
			ValidateFunc:     keyVaultValidate.NestedItemId,
			ConflictsWith:    []string{"pfx_blob", "password"},
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

		"tags": tags.Schema(),
	}
}
