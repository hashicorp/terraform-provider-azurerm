// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

// The following data comes from the results of "az maintenance public-configuration list --query "[?contains(name, `SQL`) && contains(name, `DB`)]".name --output table"
func ResourceMsSqlDatabaseMaintenanceNames() []string {
	return []string{
		"SQL_Default", "SQL_EastUS_DB_1", "SQL_EastUS2_DB_1", "SQL_SoutheastAsia_DB_1", "SQL_AustraliaEast_DB_1", "SQL_NorthEurope_DB_1", "SQL_SouthCentralUS_DB_1", "SQL_WestUS2_DB_1",
		"SQL_UKSouth_DB_1", "SQL_WestEurope_DB_1", "SQL_EastUS_DB_2", "SQL_EastUS2_DB_2", "SQL_WestUS2_DB_2", "SQL_SoutheastAsia_DB_2", "SQL_AustraliaEast_DB_2", "SQL_NorthEurope_DB_2", "SQL_SouthCentralUS_DB_2",
		"SQL_UKSouth_DB_2", "SQL_WestEurope_DB_2", "SQL_AustraliaSoutheast_DB_1", "SQL_BrazilSouth_DB_1", "SQL_CanadaCentral_DB_1", "SQL_CanadaEast_DB_1", "SQL_CentralUS_DB_1", "SQL_EastAsia_DB_1",
		"SQL_FranceCentral_DB_1", "SQL_GermanyWestCentral_DB_1", "SQL_CentralIndia_DB_1", "SQL_SouthIndia_DB_1", "SQL_JapanEast_DB_1", "SQL_JapanWest_DB_1", "SQL_NorthCentralUS_DB_1", "SQL_UKWest_DB_1",
		"SQL_WestUS_DB_1", "SQL_AustraliaSoutheast_DB_2", "SQL_BrazilSouth_DB_2", "SQL_CanadaCentral_DB_2", "SQL_CanadaEast_DB_2", "SQL_CentralUS_DB_2", "SQL_EastAsia_DB_2", "SQL_FranceCentral_DB_2",
		"SQL_GermanyWestCentral_DB_2", "SQL_CentralIndia_DB_2", "SQL_SouthIndia_DB_2", "SQL_JapanEast_DB_2", "SQL_JapanWest_DB_2", "SQL_NorthCentralUS_DB_2", "SQL_UKWest_DB_2", "SQL_WestUS_DB_2",
		"SQL_WestCentralUS_DB_1", "SQL_FranceSouth_DB_1", "SQL_WestCentralUS_DB_2", "SQL_FranceSouth_DB_2", "SQL_SwitzerlandNorth_DB_1", "SQL_SwitzerlandNorth_DB_2", "SQL_BrazilSoutheast_DB_1",
		"SQL_UAENorth_DB_1", "SQL_BrazilSoutheast_DB_2", "SQL_UAENorth_DB_2", "SQL_SouthAfricaNorth_DB_1", "SQL_SouthAfricaNorth_DB_2", "SQL_WestUS3_DB_1", "SQL_WestUS3_DB_2", "SQL_SwedenCentral_DB_1",
		"SQL_SwedenCentral_DB_2",
	}
}
