---
applyTo: "internal/**/*.go"
description: Performance optimization patterns and efficiency guidelines for the Terraform AzureRM provider including Azure API optimization, resource management, and scalability considerations.
---

# Performance Optimization Guide

Performance optimization patterns and efficiency guidelines for the Terraform AzureRM provider including Azure API optimization, resource management, and scalability considerations.

**Quick navigation:** [⚡ Azure API Efficiency](#⚡-azure-api-efficiency-patterns) | [🔄 Resource Management](#🔄-resource-management-optimization) | [📊 Monitoring](#📊-monitoring--observability-patterns) | [🚀 Scalability](#🚀-scalability-patterns)

## ⚡ Azure API Efficiency Patterns

### Batch Operations and Context Management

```go
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 30 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // Batch operations when possible
            var operations []azuretype.Operation

            // Use context with appropriate timeouts
            ctx, cancel := context.WithTimeout(ctx, 25*time.Minute)
            defer cancel()

            // Efficient resource queries with minimal API calls
            existing, err := client.Get(ctx, id)
            if err != nil && !response.WasNotFound(existing.HttpResponse) {
                return fmt.Errorf("checking for existing %s: %+v", id, err)
            }

            return nil
        },
    }
}
```

### Connection Pooling and Client Optimization

```go
// Reuse clients efficiently
client := metadata.Client.ServiceName.ResourceClient

// Connection pooling configuration
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
        IdleConnTimeout:     90 * time.Second,
    },
    Timeout: 30 * time.Second,
}
```

### Parallel Processing Patterns

```go
func processResourcesInParallel(ctx context.Context, resources []Resource) error {
    const maxWorkers = 10
    semaphore := make(chan struct{}, maxWorkers)

    var wg sync.WaitGroup
    errorsChan := make(chan error, len(resources))

    for _, resource := range resources {
        wg.Add(1)
        go func(r Resource) {
            defer wg.Done()
            semaphore <- struct{}{} // Acquire
            defer func() { <-semaphore }() // Release

            if err := processResource(ctx, r); err != nil {
                errorsChan <- err
            }
        }(resource)
    }

    wg.Wait()
    close(errorsChan)

    // Collect errors
    var errors []error
    for err := range errorsChan {
        errors = append(errors, err)
    }

    if len(errors) > 0 {
        return fmt.Errorf("processing failed: %v", errors)
    }

    return nil
}
```

## 🔄 Resource Management Optimization

### Efficient State Management

```go
// Minimize state reads and writes
func (r ServiceResource) Read() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Timeout: 5 * time.Minute,
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // Single API call for all resource data
            resp, err := client.Get(ctx, id)
            if err != nil {
                if response.WasNotFound(resp.HttpResponse) {
                    return metadata.MarkAsGone(id)
                }
                return fmt.Errorf("retrieving %s: %+v", id, err)
            }

            // Batch state updates
            state := buildCompleteState(resp.Model)
            return metadata.Encode(&state)
        },
    }
}
```

### Memory Optimization

```go
// Use object pooling for frequently allocated objects
var resourcePool = sync.Pool{
    New: func() interface{} {
        return &ResourceModel{}
    },
}

func processResource(data interface{}) error {
    model := resourcePool.Get().(*ResourceModel)
    defer func() {
        // Reset and return to pool
        *model = ResourceModel{}
        resourcePool.Put(model)
    }()

    // Process using pooled object
    return nil
}
```

### Caching Strategies

```go
// Implement intelligent caching
type CachedClient struct {
    client azure.Client
    cache  map[string]CacheEntry
    mutex  sync.RWMutex
    ttl    time.Duration
}

type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

func (c *CachedClient) GetWithCache(ctx context.Context, id string) (interface{}, error) {
    c.mutex.RLock()
    if entry, exists := c.cache[id]; exists && time.Now().Before(entry.ExpiresAt) {
        c.mutex.RUnlock()
        return entry.Data, nil
    }
    c.mutex.RUnlock()

    // Cache miss - fetch from Azure
    data, err := c.client.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // Update cache
    c.mutex.Lock()
    c.cache[id] = CacheEntry{
        Data:      data,
        ExpiresAt: time.Now().Add(c.ttl),
    }
    c.mutex.Unlock()

    return data, nil
}
```

## 📊 Monitoring & Observability Patterns

### Structured Logging with Performance Metrics

```go
func (r ServiceResource) Create() sdk.ResourceFunc {
    return sdk.ResourceFunc{
        Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
            // Enhanced logging with context
            logger := metadata.Logger.WithFields(map[string]interface{}{
                "resource_type": r.ResourceType(),
                "operation":     "create",
                "subscription":  metadata.Client.Account.SubscriptionId,
            })

            logger.Infof("Starting resource creation")

            // Performance monitoring
            start := time.Now()
            defer func() {
                logger.WithField("duration", time.Since(start)).Infof("Resource creation completed")
            }()

            // Operation implementation
            return nil
        },
    }
}
```

### Metrics Collection

```go
// Custom metrics for performance tracking
type PerformanceMetrics struct {
    APICallDuration    time.Duration
    StateUpdateTime    time.Duration
    ResourceCount      int
    ErrorRate          float64
}

func collectMetrics(operation string, duration time.Duration, success bool) {
    // Integration with monitoring systems
    labels := map[string]string{
        "operation": operation,
        "status":    map[bool]string{true: "success", false: "error"}[success],
    }

    // Record metrics (pseudocode - integrate with your monitoring system)
    recordDuration("terraform_azure_operation_duration", duration, labels)
    recordCounter("terraform_azure_operations_total", labels)
}
```

## 🚀 Scalability Patterns

### Resource Dependency Optimization

```go
// Optimize resource creation order
func optimizeResourceCreation(resources []Resource) []Resource {
    // Topological sort based on dependencies
    graph := buildDependencyGraph(resources)
    return topologicalSort(graph)
}

// Parallel resource creation where possible
func createResourcesInParallel(ctx context.Context, resources []Resource) error {
    dependencyLevels := groupByDependencyLevel(resources)

    for _, level := range dependencyLevels {
        // Create resources at same dependency level in parallel
        if err := processLevelInParallel(ctx, level); err != nil {
            return err
        }
    }

    return nil
}
```

### Large-Scale Resource Management

```go
// Chunked processing for large resource sets
func processLargeResourceSet(ctx context.Context, resources []Resource) error {
    const chunkSize = 50

    for i := 0; i < len(resources); i += chunkSize {
        end := i + chunkSize
        if end > len(resources) {
            end = len(resources)
        }

        chunk := resources[i:end]

        // Process chunk with error handling
        if err := processChunk(ctx, chunk); err != nil {
            return fmt.Errorf("processing chunk %d-%d: %+v", i, end, err)
        }

        // Rate limiting between chunks
        time.Sleep(100 * time.Millisecond)
    }

    return nil
}
```

### Memory-Efficient Data Processing

```go
// Stream processing for large datasets
func processLargeDataset(ctx context.Context, dataSource DataSource) error {
    reader := dataSource.NewReader()
    defer reader.Close()

    // Process in streaming fashion
    for {
        batch, err := reader.ReadBatch(1000) // Read in batches
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("reading batch: %+v", err)
        }

        if err := processBatch(ctx, batch); err != nil {
            return fmt.Errorf("processing batch: %+v", err)
        }

        // Force garbage collection periodically
        if reader.Position()%10000 == 0 {
            runtime.GC()
        }
    }

    return nil
}
```

## Quick Reference Links

- 🏠 **Home**: [../copilot-instructions.md](../copilot-instructions.md)
- ☁️ **Azure Patterns**: [azure-patterns.instructions.md](./azure-patterns.instructions.md)
- 🏗️ **Implementation Guide**: [implementation-guide.instructions.md](./implementation-guide.instructions.md)
- 🧪 **Testing Guide**: [testing-guidelines.instructions.md](./testing-guidelines.instructions.md)
- 📝 **Documentation Guide**: [documentation-guidelines.instructions.md](./documentation-guidelines.instructions.md)
- ❌ **Error Patterns**: [error-patterns.instructions.md](./error-patterns.instructions.md)
- 🔄 **Migration Guide**: [migration-guide.instructions.md](./migration-guide.instructions.md)
- 🏢 **Provider Guidelines**: [provider-guidelines.instructions.md](./provider-guidelines.instructions.md)
- 📐 **Schema Patterns**: [schema-patterns.instructions.md](./schema-patterns.instructions.md)
- 📋 **Code Clarity**: [code-clarity-enforcement.instructions.md](./code-clarity-enforcement.instructions.md)

---
[⬆️ Back to top](#performance-optimization-guide)
