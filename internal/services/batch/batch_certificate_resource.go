// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package batch

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2023-05-01/certificate"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceBatchCertificate() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceBatchCertificateCreate,
		Read:   resourceBatchCertificateRead,
		Update: resourceBatchCertificateUpdate,
		Delete: resourceBatchCertificateDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := certificate.ParseCertificateID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AccountName,
			},

			// TODO: make this case sensitive once this API bug has been fixed:
			// https://github.com/Azure/azure-rest-api-specs/issues/5574
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"certificate": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(1, 10000),
			},

			"format": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(certificate.CertificateFormatCer),
					string(certificate.CertificateFormatPfx),
				}, false),
			},

			"password": {
				Type:      pluginsdk.TypeString,
				Optional:  true, // Cannot be used when `format` is "Cer"
				Sensitive: true,
			},

			"thumbprint": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"thumbprint_algorithm": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringInSlice([]string{"SHA1"}, false),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"public_data": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceBatchCertificateCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch certificate creation.")

	cert := d.Get("certificate").(string)
	format := d.Get("format").(string)
	password := d.Get("password").(string)
	thumbprint := d.Get("thumbprint").(string)
	thumbprintAlgorithm := d.Get("thumbprint_algorithm").(string)
	name := thumbprintAlgorithm + "-" + thumbprint
	id := certificate.NewCertificateID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), name)

	if err := validateBatchCertificateFormatAndPassword(format, password); err != nil {
		return err
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_batch_certificate", id.ID())
		}
	}
	certificateProperties := certificate.CertificateCreateOrUpdateProperties{
		Data:                cert,
		Format:              pointer.To(certificate.CertificateFormat(format)),
		Thumbprint:          &thumbprint,
		ThumbprintAlgorithm: &thumbprintAlgorithm,
	}
	if password != "" {
		certificateProperties.Password = &password
	}
	parameters := certificate.CertificateCreateOrUpdateParameters{
		Name:       &name,
		Properties: &certificateProperties,
	}

	_, err := client.Create(ctx, id, parameters, certificate.CreateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceBatchCertificateRead(d, meta)
}

func resourceBatchCertificateRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			log.Printf("%s was not found - removing from state!", *id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.CertificateName)
	d.Set("account_name", id.BatchAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			format := ""
			if v := props.Format; v != nil {
				format = string(*v)
			}
			d.Set("format", format)

			publicData := ""
			if v := props.PublicData; v != nil {
				publicData = *v
			}
			d.Set("public_data", publicData)

			thumbprint := ""
			if v := props.Thumbprint; v != nil {
				thumbprint = *v
			}
			d.Set("thumbprint", thumbprint)

			thumbprintAlgorithm := ""
			if v := props.ThumbprintAlgorithm; v != nil {
				thumbprintAlgorithm = *v
			}
			d.Set("thumbprint_algorithm", thumbprintAlgorithm)
		}
	}

	return nil
}

func resourceBatchCertificateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure Batch certificate update.")

	id, err := certificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	cert := d.Get("certificate").(string)
	format := d.Get("format").(string)
	password := d.Get("password").(string)
	thumbprint := d.Get("thumbprint").(string)
	thumbprintAlgorithm := d.Get("thumbprint_algorithm").(string)

	if err := validateBatchCertificateFormatAndPassword(format, password); err != nil {
		return err
	}

	parameters := certificate.CertificateCreateOrUpdateParameters{
		Name: &id.CertificateName,
		Properties: &certificate.CertificateCreateOrUpdateProperties{
			Data:                cert,
			Format:              pointer.To(certificate.CertificateFormat(format)),
			Password:            &password,
			Thumbprint:          &thumbprint,
			ThumbprintAlgorithm: &thumbprintAlgorithm,
		},
	}

	if _, err = client.Update(ctx, *id, parameters, certificate.UpdateOperationOptions{}); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	_, err = client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return resourceBatchCertificateRead(d, meta)
}

func resourceBatchCertificateDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Batch.CertificateClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := certificate.ParseCertificateID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func validateBatchCertificateFormatAndPassword(format string, password string) error {
	if format == "Cer" && password != "" {
		return fmt.Errorf(" Batch Certificate Password must not be specified when Format is `Cer`")
	}
	return nil
}
