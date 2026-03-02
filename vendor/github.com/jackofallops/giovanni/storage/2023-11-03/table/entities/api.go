package entities

import (
	"context"
)

type StorageTableEntity interface {
	Delete(ctx context.Context, tableName string, input DeleteEntityInput) (resp DeleteEntityResponse, err error)
	Insert(ctx context.Context, tableName string, input InsertEntityInput) (resp InsertResponse, err error)
	InsertOrReplace(ctx context.Context, tableName string, input InsertOrReplaceEntityInput) (resp InsertOrReplaceResponse, err error)
	InsertOrMerge(ctx context.Context, tableName string, input InsertOrMergeEntityInput) (resp InsertOrMergeResponse, err error)
	Query(ctx context.Context, tableName string, input QueryEntitiesInput) (resp QueryEntitiesResponse, err error)
	Get(ctx context.Context, tableName string, input GetEntityInput) (resp GetEntityResponse, err error)
}
