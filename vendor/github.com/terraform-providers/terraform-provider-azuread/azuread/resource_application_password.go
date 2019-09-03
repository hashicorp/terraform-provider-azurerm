package azuread

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/ar"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func resourceApplicationPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationPasswordCreate,
		Read:   resourceApplicationPasswordRead,
		Delete: resourceApplicationPasswordDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		// Schema: graph.PasswordResourceSchema("application_object"), //todo switch back to this in 1.0
		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ValidateFunc:  validate.UUID,
				Deprecated:    "Deprecated in favour of `application_object_id` to prevent confusion",
				ConflictsWith: []string{"application_object_id"},
			},

			"application_object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ValidateFunc:  validate.UUID,
				ConflictsWith: []string{"application_id"},
			},

			"key_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUID,
			},

			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"start_date": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},

			"end_date": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"end_date_relative"},
				ValidateFunc:  validation.ValidateRFC3339TimeString,
			},

			"end_date_relative": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"end_date"},
				ValidateFunc:  validate.NoEmptyStrings,
			},
		},
	}
}

func resourceApplicationPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	objectId := d.Get("application_object_id").(string)
	if objectId == "" { // todo remove in 1.0
		objectId = d.Get("application_id").(string)
	}
	if objectId == "" {
		return fmt.Errorf("one of `application_object_id` or `application_id` must be specified")
	}

	cred, err := graph.PasswordCredentialForResource(d)
	if err != nil {
		return fmt.Errorf("Error generating Application Credentials for Object ID %q: %+v", objectId, err)
	}
	id := graph.PasswordCredentialIdFrom(objectId, *cred.KeyID)

	tf.LockByName(resourceApplicationName, id.ObjectId)
	defer tf.UnlockByName(resourceApplicationName, id.ObjectId)

	existingCreds, err := client.ListPasswordCredentials(ctx, id.ObjectId)
	if err != nil {
		return fmt.Errorf("Error Listing Application Credentials for Object ID %q: %+v", id.ObjectId, err)
	}

	newCreds, err := graph.PasswordCredentialResultAdd(existingCreds, cred, requireResourcesToBeImported)
	if err != nil {
		return tf.ImportAsExistsError("azuread_application_password", id.String())
	}

	if _, err = client.UpdatePasswordCredentials(ctx, id.ObjectId, graphrbac.PasswordCredentialsUpdateParameters{Value: newCreds}); err != nil {
		return fmt.Errorf("Error creating Application Credentials %q for Object ID %q: %+v", id.KeyId, id.ObjectId, err)
	}

	_, err = graph.WaitForPasswordCredentialReplication(id.KeyId, func() (graphrbac.PasswordCredentialListResult, error) {
		return client.ListPasswordCredentials(ctx, id.ObjectId)
	})
	if err != nil {
		return fmt.Errorf("Error waiting for Application Password replication (AppID %q, KeyID %q: %+v", id.ObjectId, id.KeyId, err)
	}

	d.SetId(id.String())

	return resourceApplicationPasswordRead(d, meta)
}

func resourceApplicationPasswordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := graph.ParsePasswordCredentialId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Application Password ID: %v", err)
	}
	// ensure the Application Object exists
	app, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		// the parent Service Principal has been removed - skip it
		if ar.ResponseWasNotFound(app.Response) {
			log.Printf("[DEBUG] Application with Object ID %q was not found - removing from state!", id.ObjectId)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Application ID %q: %+v", id.ObjectId, err)
	}

	credentials, err := client.ListPasswordCredentials(ctx, id.ObjectId)
	if err != nil {
		return fmt.Errorf("Error Listing Application Credentials for Application with Object ID %q: %+v", id.ObjectId, err)
	}

	credential := graph.PasswordCredentialResultFindByKeyId(credentials, id.KeyId)
	if credential == nil {
		log.Printf("[DEBUG] Application Credentials %q (ID %q) was not found - removing from state!", id.KeyId, id.ObjectId)
		d.SetId("")
		return nil
	}

	// todo, move this into a graph helper function?
	d.Set("application_object_id", id.ObjectId)
	d.Set("application_id", id.ObjectId) //todo remove in 2.0
	d.Set("key_id", id.KeyId)

	if endDate := credential.EndDate; endDate != nil {
		d.Set("end_date", endDate.Format(time.RFC3339))
	}

	if startDate := credential.StartDate; startDate != nil {
		d.Set("start_date", startDate.Format(time.RFC3339))
	}

	return nil
}

func resourceApplicationPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).applicationsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := graph.ParsePasswordCredentialId(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Application Password ID: %v", err)
	}

	tf.LockByName(resourceApplicationName, id.ObjectId)
	defer tf.UnlockByName(resourceApplicationName, id.ObjectId)

	// ensure the parent Application exists
	app, err := client.Get(ctx, id.ObjectId)
	if err != nil {
		// the parent Service Principal has been removed - skip it
		if ar.ResponseWasNotFound(app.Response) {
			log.Printf("[DEBUG] Application with Object ID %q was not found - removing from state!", id.ObjectId)
			return nil
		}
		return fmt.Errorf("Error retrieving Application ID %q: %+v", id.ObjectId, err)
	}

	existing, err := client.ListPasswordCredentials(ctx, id.ObjectId)
	if err != nil {
		return fmt.Errorf("Error Listing Application Credentials for %q: %+v", id.ObjectId, err)
	}

	newCreds := graph.PasswordCredentialResultRemoveByKeyId(existing, id.KeyId)
	if _, err = client.UpdatePasswordCredentials(ctx, id.ObjectId, graphrbac.PasswordCredentialsUpdateParameters{Value: newCreds}); err != nil {
		return fmt.Errorf("Error removing Application Credentials %q from Application Object ID %q: %+v", id.KeyId, id.ObjectId, err)
	}

	return nil
}
