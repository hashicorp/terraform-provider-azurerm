// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/machinelearningcomputes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	synapseValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSynapseSpark() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseSparkCreate,
		Read:   resourceSynapseSparkRead,
		Delete: resourceSynapseSparkDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := machinelearningcomputes.ParseComputeID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{2,16}$`),
					"It can include letters, digits and dashes. It must start with a letter, end with a letter or digit, and be between 2 and 16 characters in length."),
			},

			"machine_learning_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspaces.ValidateWorkspaceID,
			},

			"location": commonschema.Location(),

			"synapse_spark_pool_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: synapseValidate.SparkPoolID,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptionalForceNew(),

			"local_auth_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},

			"tags": commonschema.TagsForceNew(),
		},
	}
}

func resourceSynapseSparkCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspaceID, _ := workspaces.ParseWorkspaceID(d.Get("machine_learning_workspace_id").(string))
	id := machinelearningcomputes.NewComputeID(subscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.ComputeGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing Machine Learning Compute (%q): %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_machine_learning_synapse_spark", id.ID())
		}
	}

	identity, err := expandIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := machinelearningcomputes.ComputeResource{
		Properties: &machinelearningcomputes.SynapseSpark{
			Properties:       nil,
			ComputeLocation:  utils.String(d.Get("location").(string)),
			Description:      utils.String(d.Get("description").(string)),
			ResourceId:       utils.String(d.Get("synapse_spark_pool_id").(string)),
			DisableLocalAuth: utils.Bool(!d.Get("local_auth_enabled").(bool)),
		},
		Identity: identity,
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.ComputeCreateOrUpdate(ctx, id, parameters)
	if err != nil {
		return fmt.Errorf("creating Machine Learning Compute (%q): %+v", id, err)
	}
	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of Machine Learning Compute (%q): %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseSparkRead(d, meta)
}

func resourceSynapseSparkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.ComputeGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Machine Learning Compute %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Machine Learning Compute (%q): %+v", id, err)
	}

	d.Set("name", id.ComputeName)
	workspaceId := workspaces.NewWorkspaceID(subscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("machine_learning_workspace_id", workspaceId.ID())

	if location := resp.Model.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	identity, err := flattenIdentity(resp.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	props := resp.Model.Properties.(machinelearningcomputes.SynapseSpark)
	if props.DisableLocalAuth != nil {
		d.Set("local_auth_enabled", !*props.DisableLocalAuth)
	}
	d.Set("description", props.Description)
	d.Set("synapse_spark_pool_id", props.ResourceId)

	return tags.FlattenAndSet(d, resp.Model.Tags)
}

func resourceSynapseSparkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MachineLearning.MachineLearningComputes
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	id, err := machinelearningcomputes.ParseComputeID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.ComputeDelete(ctx, *id, machinelearningcomputes.ComputeDeleteOperationOptions{
		UnderlyingResourceAction: pointer.To(machinelearningcomputes.UnderlyingResourceActionDetach),
	})
	if err != nil {
		return fmt.Errorf("deleting Machine Learning Compute (%q): %+v", id, err)
	}

	if err := future.Poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for deletion of the Machine Learning Compute (%q): %+v", id, err)
	}
	return nil
}
