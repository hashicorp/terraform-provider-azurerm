package azurerm

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/sql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/satori/uuid"
)

func resourceArmSqlDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSqlDatabaseCreateOrUpdate,
		Read:   resourceArmSqlDatabaseRead,
		Update: resourceArmSqlDatabaseCreateOrUpdate,
		Delete: resourceArmSqlDatabaseDelete,

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

			"server_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"create_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "Default",
				ValidateFunc: validateArmSqlDatabaseCreateMode,
			},

			"source_database_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"restore_point_in_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"edition": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateArmSqlDatabaseEdition,
			},

			"collation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_size_bytes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"requested_service_objective_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"requested_service_objective_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source_database_deletion_date": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"elastic_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"encryption": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"default_secondary_location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmSqlDatabaseCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	sqlDatabasesClient := meta.(*ArmClient).sqlDatabasesClient

	databaseName := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	serverName := d.Get("server_name").(string)
	sourceDatabaseID := d.Get("source_database_id").(string)
	collation := d.Get("collation").(string)
	elasticPoolName := d.Get("elastic_pool_name").(string)

	createMode := sql.CreateMode(d.Get("create_mode").(string))

	maxSizeBytes := d.Get("max_size_bytes").(string)

	edition := sql.DatabaseEdition(d.Get("edition").(string))

	requestedServiceObjectiveIDString := d.Get("requested_service_objective_id").(string)
	var requestedServiceObjectiveID uuid.UUID

	requestedServiceObjectiveNameString := d.Get("requested_service_objective_name").(string)
	var requestedServiceObjectiveName sql.ServiceObjectiveName

	if requestedServiceObjectiveIDString != "" && requestedServiceObjectiveNameString == "" {

		var errUUID error
		requestedServiceObjectiveID, errUUID = uuid.FromString(requestedServiceObjectiveIDString)
		if errUUID != nil {
			return errUUID
		}

	} else if requestedServiceObjectiveIDString == "" && requestedServiceObjectiveNameString != "" {

		requestedServiceObjectiveName = sql.ServiceObjectiveName(requestedServiceObjectiveNameString)

	} else {

		return fmt.Errorf("either service objective name or id must be specified")

	}

	tags := d.Get("tags").(map[string]interface{})
	metadata := expandTags(tags)

	props := sql.DatabaseProperties{
		Collation:  &collation,
		Edition:    edition,
		CreateMode: createMode,
	}

	if createMode == sql.Default && maxSizeBytes != "" {
		props.MaxSizeBytes = &maxSizeBytes
		props.SourceDatabaseID = nil
	} else if sourceDatabaseID != "" {
		props.MaxSizeBytes = nil
		props.SourceDatabaseID = &sourceDatabaseID
	}

	if requestedServiceObjectiveName != "" {
		props.RequestedServiceObjectiveName = requestedServiceObjectiveName
	} else {
		props.RequestedServiceObjectiveID = &requestedServiceObjectiveID
	}

	if requestedServiceObjectiveNameString == "ElasticPool" {
		props.ElasticPoolName = &elasticPoolName
	}

	parameters := sql.Database{
		Name:               &databaseName,
		DatabaseProperties: &props,
		Tags:               metadata,
		Location:           &location,
	}

	resultChan, errChan := sqlDatabasesClient.CreateOrUpdate(resGroup, serverName, databaseName, parameters, make(chan struct{}))
	resultServer := <-resultChan
	errServer := <-errChan
	if errServer != nil {
		return errServer
	}

	if resultServer.StatusCode != http.StatusOK {
		return fmt.Errorf("Cannot create Sql database %s (resource group %s) ID", databaseName, resGroup)
	}

	resultDb, errDb := sqlDatabasesClient.Get(resGroup, serverName, databaseName, "")
	if errDb != nil {
		d.SetId("")
		return fmt.Errorf("Error reading SQL database %s: %v", databaseName, errDb)
	}

	d.SetId(*resultDb.ID)

	return resourceArmSqlDatabaseRead(d, meta)
}

func resourceArmSqlDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	sqlDatabasesClient := meta.(*ArmClient).sqlDatabasesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	databaseName := id.Path["databases"]

	result, err := sqlDatabasesClient.Get(resGroup, serverName, databaseName, "")
	if err != nil {
		return fmt.Errorf("Error reading SQL database %s: %v", databaseName, err)
	}
	if result.Response.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	databaseProperties := *result.DatabaseProperties

	d.Set("name", databaseName)
	d.Set("resource_group_name", resGroup)
	d.Set("location", *result.Location)
	d.Set("server_name", serverName)
	d.Set("creation_date", *result.CreationDate)
	d.Set("default_secondary_location", *result.DefaultSecondaryLocation)

	if databaseProperties.ElasticPoolName != nil {
		d.Set("elastic_pool_name", *databaseProperties.ElasticPoolName)
	}

	flattenAndSetTags(d, result.Tags)

	return nil
}

func resourceArmSqlDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	sqlDatabasesClient := meta.(*ArmClient).sqlDatabasesClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	serverName := id.Path["servers"]
	databaseName := id.Path["databases"]

	result, error := sqlDatabasesClient.Delete(resGroup, serverName, databaseName)
	if result.Response.StatusCode != http.StatusOK {
		return fmt.Errorf("Error deleting SQL database %s: %s", databaseName, error)
	}

	return nil
}

func validateArmSqlDatabaseEdition(v interface{}, k string) (ws []string, errors []error) {
	editions := map[string]bool{
		"Basic":         true,
		"Standard":      true,
		"Premium":       true,
		"DataWarehouse": true,
	}
	if !editions[v.(string)] {
		errors = append(errors, fmt.Errorf("SQL Database Edition can only be Basic, Standard, Premium or DataWarehouse"))
	}
	return
}

func validateArmSqlDatabaseCreateMode(v interface{}, k string) (ws []string, errors []error) {
	modes := map[string]bool{
		"Copy":                           true,
		"Default":                        true,
		"NonReadableSecondary":           true,
		"OnlineSecondary":                true,
		"PointInTimeRestore":             true,
		"Recovery":                       true,
		"Restore":                        true,
		"RestoreLongTermRetentionBackup": true,
	}
	if !modes[v.(string)] {
		errors = append(errors, fmt.Errorf("SQL Database Create Mode can only be Copy, Default, NonReadableSecondary, OnlineSecondary, PointInTimeRestore, Recovery, Restore, RestoreLongTermRetentionBackup"))
	}
	return
}
