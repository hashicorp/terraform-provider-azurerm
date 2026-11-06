package openapis

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CassandraKeyspaceGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p CassandraKeyspaceGetResultsOperationPredicate) Matches(input CassandraKeyspaceGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type CassandraRoleAssignmentResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p CassandraRoleAssignmentResourceOperationPredicate) Matches(input CassandraRoleAssignmentResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type CassandraRoleDefinitionResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p CassandraRoleDefinitionResourceOperationPredicate) Matches(input CassandraRoleDefinitionResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type CassandraTableGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p CassandraTableGetResultsOperationPredicate) Matches(input CassandraTableGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type ClientEncryptionKeyGetResultsOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ClientEncryptionKeyGetResultsOperationPredicate) Matches(input ClientEncryptionKeyGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type ClusterResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p ClusterResourceOperationPredicate) Matches(input ClusterResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type DataCenterResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p DataCenterResourceOperationPredicate) Matches(input DataCenterResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type DatabaseAccountGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p DatabaseAccountGetResultsOperationPredicate) Matches(input DatabaseAccountGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type FleetResourceOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p FleetResourceOperationPredicate) Matches(input FleetResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && *p.Location != input.Location {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type FleetspaceAccountResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p FleetspaceAccountResourceOperationPredicate) Matches(input FleetspaceAccountResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type FleetspaceResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p FleetspaceResourceOperationPredicate) Matches(input FleetspaceResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type GremlinDatabaseGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p GremlinDatabaseGetResultsOperationPredicate) Matches(input GremlinDatabaseGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type GremlinGraphGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p GremlinGraphGetResultsOperationPredicate) Matches(input GremlinGraphGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type GremlinRoleAssignmentResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p GremlinRoleAssignmentResourceOperationPredicate) Matches(input GremlinRoleAssignmentResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type GremlinRoleDefinitionResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p GremlinRoleDefinitionResourceOperationPredicate) Matches(input GremlinRoleDefinitionResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type LocationGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p LocationGetResultOperationPredicate) Matches(input LocationGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type MetricOperationPredicate struct {
	EndTime   *string
	StartTime *string
	TimeGrain *string
}

func (p MetricOperationPredicate) Matches(input Metric) bool {

	if p.EndTime != nil && (input.EndTime == nil || *p.EndTime != *input.EndTime) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	if p.TimeGrain != nil && (input.TimeGrain == nil || *p.TimeGrain != *input.TimeGrain) {
		return false
	}

	return true
}

type MetricDefinitionOperationPredicate struct {
	ResourceUri *string
}

func (p MetricDefinitionOperationPredicate) Matches(input MetricDefinition) bool {

	if p.ResourceUri != nil && (input.ResourceUri == nil || *p.ResourceUri != *input.ResourceUri) {
		return false
	}

	return true
}

type MongoDBCollectionGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p MongoDBCollectionGetResultsOperationPredicate) Matches(input MongoDBCollectionGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type MongoDBDatabaseGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p MongoDBDatabaseGetResultsOperationPredicate) Matches(input MongoDBDatabaseGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type MongoMIRoleAssignmentResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p MongoMIRoleAssignmentResourceOperationPredicate) Matches(input MongoMIRoleAssignmentResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type MongoMIRoleDefinitionResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p MongoMIRoleDefinitionResourceOperationPredicate) Matches(input MongoMIRoleDefinitionResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type MongoRoleDefinitionGetResultsOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p MongoRoleDefinitionGetResultsOperationPredicate) Matches(input MongoRoleDefinitionGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type MongoUserDefinitionGetResultsOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p MongoUserDefinitionGetResultsOperationPredicate) Matches(input MongoUserDefinitionGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type PartitionMetricOperationPredicate struct {
	EndTime             *string
	PartitionId         *string
	PartitionKeyRangeId *string
	StartTime           *string
	TimeGrain           *string
}

func (p PartitionMetricOperationPredicate) Matches(input PartitionMetric) bool {

	if p.EndTime != nil && (input.EndTime == nil || *p.EndTime != *input.EndTime) {
		return false
	}

	if p.PartitionId != nil && (input.PartitionId == nil || *p.PartitionId != *input.PartitionId) {
		return false
	}

	if p.PartitionKeyRangeId != nil && (input.PartitionKeyRangeId == nil || *p.PartitionKeyRangeId != *input.PartitionKeyRangeId) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	if p.TimeGrain != nil && (input.TimeGrain == nil || *p.TimeGrain != *input.TimeGrain) {
		return false
	}

	return true
}

type PartitionUsageOperationPredicate struct {
	CurrentValue        *int64
	Limit               *int64
	PartitionId         *string
	PartitionKeyRangeId *string
	QuotaPeriod         *string
}

func (p PartitionUsageOperationPredicate) Matches(input PartitionUsage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.PartitionId != nil && (input.PartitionId == nil || *p.PartitionId != *input.PartitionId) {
		return false
	}

	if p.PartitionKeyRangeId != nil && (input.PartitionKeyRangeId == nil || *p.PartitionKeyRangeId != *input.PartitionKeyRangeId) {
		return false
	}

	if p.QuotaPeriod != nil && (input.QuotaPeriod == nil || *p.QuotaPeriod != *input.QuotaPeriod) {
		return false
	}

	return true
}

type PercentileMetricOperationPredicate struct {
	EndTime   *string
	StartTime *string
	TimeGrain *string
}

func (p PercentileMetricOperationPredicate) Matches(input PercentileMetric) bool {

	if p.EndTime != nil && (input.EndTime == nil || *p.EndTime != *input.EndTime) {
		return false
	}

	if p.StartTime != nil && (input.StartTime == nil || *p.StartTime != *input.StartTime) {
		return false
	}

	if p.TimeGrain != nil && (input.TimeGrain == nil || *p.TimeGrain != *input.TimeGrain) {
		return false
	}

	return true
}

type RestorableDatabaseAccountGetResultOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p RestorableDatabaseAccountGetResultOperationPredicate) Matches(input RestorableDatabaseAccountGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableGremlinDatabaseGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableGremlinDatabaseGetResultOperationPredicate) Matches(input RestorableGremlinDatabaseGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableGremlinGraphGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableGremlinGraphGetResultOperationPredicate) Matches(input RestorableGremlinGraphGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableGremlinResourcesGetResultOperationPredicate struct {
	DatabaseName *string
	Id           *string
	Name         *string
	Type         *string
}

func (p RestorableGremlinResourcesGetResultOperationPredicate) Matches(input RestorableGremlinResourcesGetResult) bool {

	if p.DatabaseName != nil && (input.DatabaseName == nil || *p.DatabaseName != *input.DatabaseName) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableMongodbCollectionGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableMongodbCollectionGetResultOperationPredicate) Matches(input RestorableMongodbCollectionGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableMongodbDatabaseGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableMongodbDatabaseGetResultOperationPredicate) Matches(input RestorableMongodbDatabaseGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableMongodbResourcesGetResultOperationPredicate struct {
	DatabaseName *string
	Id           *string
	Name         *string
	Type         *string
}

func (p RestorableMongodbResourcesGetResultOperationPredicate) Matches(input RestorableMongodbResourcesGetResult) bool {

	if p.DatabaseName != nil && (input.DatabaseName == nil || *p.DatabaseName != *input.DatabaseName) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableSqlContainerGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableSqlContainerGetResultOperationPredicate) Matches(input RestorableSqlContainerGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableSqlDatabaseGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableSqlDatabaseGetResultOperationPredicate) Matches(input RestorableSqlDatabaseGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableSqlResourcesGetResultOperationPredicate struct {
	DatabaseName *string
	Id           *string
	Name         *string
	Type         *string
}

func (p RestorableSqlResourcesGetResultOperationPredicate) Matches(input RestorableSqlResourcesGetResult) bool {

	if p.DatabaseName != nil && (input.DatabaseName == nil || *p.DatabaseName != *input.DatabaseName) {
		return false
	}

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableTableGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableTableGetResultOperationPredicate) Matches(input RestorableTableGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type RestorableTableResourcesGetResultOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p RestorableTableResourcesGetResultOperationPredicate) Matches(input RestorableTableResourcesGetResult) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type ServiceResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p ServiceResourceOperationPredicate) Matches(input ServiceResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type SqlContainerGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p SqlContainerGetResultsOperationPredicate) Matches(input SqlContainerGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type SqlDatabaseGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p SqlDatabaseGetResultsOperationPredicate) Matches(input SqlDatabaseGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type SqlRoleAssignmentGetResultsOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p SqlRoleAssignmentGetResultsOperationPredicate) Matches(input SqlRoleAssignmentGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type SqlRoleDefinitionGetResultsOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p SqlRoleDefinitionGetResultsOperationPredicate) Matches(input SqlRoleDefinitionGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type SqlStoredProcedureGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p SqlStoredProcedureGetResultsOperationPredicate) Matches(input SqlStoredProcedureGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type SqlTriggerGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p SqlTriggerGetResultsOperationPredicate) Matches(input SqlTriggerGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type SqlUserDefinedFunctionGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p SqlUserDefinedFunctionGetResultsOperationPredicate) Matches(input SqlUserDefinedFunctionGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type TableGetResultsOperationPredicate struct {
	Id       *string
	Location *string
	Name     *string
	Type     *string
}

func (p TableGetResultsOperationPredicate) Matches(input TableGetResults) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Location != nil && (input.Location == nil || *p.Location != *input.Location) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type TableRoleAssignmentResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p TableRoleAssignmentResourceOperationPredicate) Matches(input TableRoleAssignmentResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type TableRoleDefinitionResourceOperationPredicate struct {
	Id   *string
	Name *string
	Type *string
}

func (p TableRoleDefinitionResourceOperationPredicate) Matches(input TableRoleDefinitionResource) bool {

	if p.Id != nil && (input.Id == nil || *p.Id != *input.Id) {
		return false
	}

	if p.Name != nil && (input.Name == nil || *p.Name != *input.Name) {
		return false
	}

	if p.Type != nil && (input.Type == nil || *p.Type != *input.Type) {
		return false
	}

	return true
}

type UsageOperationPredicate struct {
	CurrentValue *int64
	Limit        *int64
	QuotaPeriod  *string
}

func (p UsageOperationPredicate) Matches(input Usage) bool {

	if p.CurrentValue != nil && (input.CurrentValue == nil || *p.CurrentValue != *input.CurrentValue) {
		return false
	}

	if p.Limit != nil && (input.Limit == nil || *p.Limit != *input.Limit) {
		return false
	}

	if p.QuotaPeriod != nil && (input.QuotaPeriod == nil || *p.QuotaPeriod != *input.QuotaPeriod) {
		return false
	}

	return true
}
