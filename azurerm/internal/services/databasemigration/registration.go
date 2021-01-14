package databasemigration

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Database Migration"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Database Migration",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_database_migration_service": dataSourceDatabaseMigrationService(),
		"azurerm_database_migration_project": dataSourceDatabaseMigrationProject(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	resources := map[string]*schema.Resource{
		"azurerm_database_migration_service": resourceDatabaseMigrationService(),
		"azurerm_database_migration_project": resourceDatabaseMigrationProject(),
	}

	return resources
}
