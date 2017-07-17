package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jen20/riviera/search"
)

func resourceArmSearchService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSearchServiceCreate,
		Read:   resourceArmSearchServiceRead,
		Update: resourceArmSearchServiceCreate,
		Delete: resourceArmSearchServiceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Required: true,
			},

			"replica_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"partition_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSearchServiceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	rivieraClient := client.rivieraClient

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	command := &search.CreateOrUpdateSearchService{
		Name:              d.Get("name").(string),
		Location:          d.Get("location").(string),
		ResourceGroupName: d.Get("resource_group_name").(string),
		Tags:              *expandedTags,
		Sku: search.Sku{
			Name: d.Get("sku").(string),
		},
	}

	if v, ok := d.GetOk("replica_count"); ok {
		replica_count := v.(int)
		command.ReplicaCount = &replica_count
	}

	if v, ok := d.GetOk("partition_count"); ok {
		partition_count := v.(int)
		command.PartitionCount = &partition_count
	}

	createRequest := rivieraClient.NewRequest()
	createRequest.Command = command

	createResponse, err := createRequest.Execute()
	if err != nil {
		return fmt.Errorf("Error creating Search Service: %+v", err)
	}
	if !createResponse.IsSuccessful() {
		return fmt.Errorf("Error creating Search Service: %+v", createResponse.Error)
	}

	getSearchServiceCommand := &search.GetSearchService{
		Name:              d.Get("name").(string),
		ResourceGroupName: d.Get("resource_group_name").(string),
	}

	readRequest := rivieraClient.NewRequest()
	readRequest.Command = getSearchServiceCommand

	readResponse, err := readRequest.Execute()
	if err != nil {
		return fmt.Errorf("Error reading Search Service: %+v", err)
	}
	if !readResponse.IsSuccessful() {
		return fmt.Errorf("Error reading Search Service: %+v", readResponse.Error)
	}
	resp := readResponse.Parsed.(*search.GetSearchServiceResponse)

	log.Printf("[DEBUG] Waiting for Search Service (%s) to become available", d.Get("name"))
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"provisioning"},
		Target:     []string{"succeeded"},
		Refresh:    azureStateRefreshFunc(*resp.ID, client, getSearchServiceCommand),
		Timeout:    30 * time.Minute,
		MinTimeout: 15 * time.Second,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Search Service (%s) to become available: %+v", d.Get("name"), err)
	}

	d.SetId(*resp.ID)

	return resourceArmSearchServiceRead(d, meta)
}

func resourceArmSearchServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	rivieraClient := client.rivieraClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["searchServices"]

	readRequest := rivieraClient.NewRequestForURI(d.Id())
	readRequest.Command = &search.GetSearchService{}

	readResponse, err := readRequest.Execute()
	if err != nil {
		return fmt.Errorf("Error reading Search Service: %+v", err)
	}
	if !readResponse.IsSuccessful() {
		log.Printf("[INFO] Error reading Search Service %q - removing from state", d.Id())
		d.SetId("")
		return fmt.Errorf("Error reading Search Service: %+v", readResponse.Error)
	}

	resp := readResponse.Parsed.(*search.GetSearchServiceResponse)

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", azureRMNormalizeLocation(resp.Location))
	d.Set("sku", string(resp.Sku.Name))

	if resp.PartitionCount != nil {
		d.Set("partition_count", resp.PartitionCount)
	}

	if resp.ReplicaCount != nil {
		d.Set("replica_count", resp.ReplicaCount)
	}

	flattenAndSetTags(d, &resp.Tags)

	return nil
}

func resourceArmSearchServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	rivieraClient := client.rivieraClient

	deleteRequest := rivieraClient.NewRequestForURI(d.Id())
	deleteRequest.Command = &search.DeleteSearchService{}

	deleteResponse, err := deleteRequest.Execute()
	if err != nil {
		return fmt.Errorf("Error deleting Search Service: %+v", err)
	}
	if !deleteResponse.IsSuccessful() {
		return fmt.Errorf("Error deleting Search Service: %+v", deleteResponse.Error)
	}

	return nil
}
