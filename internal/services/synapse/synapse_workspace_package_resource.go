package synapse

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/mitchellh/go-homedir"
)

func resourceSynapseWorkspacePackage() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseWorkspacePackageCreate,
		Read:   resourceSynapseWorkspacePackageRead,
		Delete: resourceSynapseWorkspacePackageDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WorkspacePackageID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"source": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"source_md5": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"container_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"path": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSynapseWorkspacePackageCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.LibraryClient(workspaceId.Name, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	id := parse.NewWorkspacePackageID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.LibraryName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_workspace_package", id.ID())
		}
	}

	filepath := d.Get("source").(string)
	filepath, err = homedir.Expand(filepath)
	if err != nil {
		return fmt.Errorf("expanding homedir in source (%s): %+v", filepath, err)
	}
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("opening file: %+v", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("[WARN] Error closing filepath (%s): %+v", filepath, err)
		}
	}()

	future, err := client.Create(ctx, id.LibraryName)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	var tmp = make([]byte, 1024*1024*10)
	for {
		n, err := file.Read(tmp)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading file: %+v", err)
		}
		_, err = client.Append(ctx, id.LibraryName, string(tmp[:n]), nil)
		if err != nil {
			return fmt.Errorf("uploading file: %+v", err)
		}
	}

	flushFuture, err := client.Flush(ctx, id.LibraryName)
	if err != nil {
		return fmt.Errorf("flushing %s: %+v", id, err)
	}

	if err = flushFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on flushing for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSynapseWorkspacePackageRead(d, meta)
}

func resourceSynapseWorkspacePackageRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.WorkspacePackageID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.LibraryClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.LibraryName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.LibraryName)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	d.Set("source", d.Get("source"))
	d.Set("source_md5", d.Get("source_md5"))
	if props := resp.Properties; props != nil {
		d.Set("type", props.Type)
		d.Set("container_name", props.ContainerName)
		d.Set("path", props.Path)
	}

	return nil
}

func resourceSynapseWorkspacePackageDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment

	id, err := parse.WorkspacePackageID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.LibraryClient(id.WorkspaceName, environment.SynapseEndpointSuffix)
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.LibraryName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}
