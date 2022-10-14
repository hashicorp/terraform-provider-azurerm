package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/manicminer/hamilton/odata"
)

// SynchronizationJobClient performs operations on SynchronizationJobs.
type SynchronizationJobClient struct {
	BaseClient Client
}

// NewSynchronizationJobClient returns a new SynchronizationJobClient
func NewSynchronizationJobClient(tenantId string) *SynchronizationJobClient {
	return &SynchronizationJobClient{
		BaseClient: NewClient(VersionBeta, tenantId),
	}
}

// Api calls give UnknownError on consistency errors
func ServicePrincipalDoesNotExistConsistency(resp *http.Response, o *odata.OData) bool {
	return o != nil && o.Error != nil && resp.StatusCode == http.StatusUnauthorized
}

// Api call give StatusConflict on consistency errors
func ConflictConsistencyFailureFunc(resp *http.Response, o *odata.OData) bool {
	return o != nil && o.Error != nil && resp.StatusCode == http.StatusConflict
}

// List returns a list of SynchronizationJobs
func (c *SynchronizationJobClient) List(ctx context.Context, servicePrincipalId string) (*[]SynchronizationJob, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes: []int{http.StatusOK},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs", servicePrincipalId),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var data struct {
		SynchronizationJobs []SynchronizationJob `json:"value"`
	}
	if err := json.Unmarshal(respBody, &data); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &data.SynchronizationJobs, status, nil
}

// Get retrieves a SynchronizationJob
func (c *SynchronizationJobClient) Get(ctx context.Context, id string, servicePrincipalId string) (*SynchronizationJob, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes:       []int{http.StatusOK},
		ConsistencyFailureFunc: ServicePrincipalDoesNotExistConsistency,
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s", servicePrincipalId, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var synchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &synchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &synchronizationJob, status, nil
}

// GetSecrets retrieves a SynchronizationSecret
func (c *SynchronizationJobClient) GetSecrets(ctx context.Context, servicePrincipalId string) (*SynchronizationSecret, int, error) {
	resp, status, _, err := c.BaseClient.Get(ctx, GetHttpRequestInput{
		ValidStatusCodes:       []int{http.StatusOK},
		ConsistencyFailureFunc: ServicePrincipalDoesNotExistConsistency,
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/secrets", servicePrincipalId),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Get(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var synchronizationSecret SynchronizationSecret
	if err := json.Unmarshal(respBody, &synchronizationSecret); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &synchronizationSecret, status, nil
}

// Adds a SynchronizationSecrets.
func (c *SynchronizationJobClient) SetSecrets(ctx context.Context, synchronizationSecret SynchronizationSecret, servicePrincipalId string) (int, error) {
	var status int

	body, err := json.Marshal(synchronizationSecret)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}
	resp, status, _, err := c.BaseClient.Put(ctx, PutHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: ServicePrincipalDoesNotExistConsistency,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/secrets", servicePrincipalId),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("SynchronizationJobClient.BaseClient.Put(): %v", err)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	return status, nil
}

// Creates a SynchronizationJob.
func (c *SynchronizationJobClient) Create(ctx context.Context, synchronizationJob SynchronizationJob, servicePrincipalId string) (*SynchronizationJob, int, error) {
	var status int

	body, err := json.Marshal(synchronizationJob)
	if err != nil {
		return nil, status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:                   body,
		ConsistencyFailureFunc: ServicePrincipalDoesNotExistConsistency,
		ValidStatusCodes:       []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs", servicePrincipalId),
			HasTenantId: true,
		},
	})
	if err != nil {
		return nil, status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	var newSynchronizationJob SynchronizationJob
	if err := json.Unmarshal(respBody, &newSynchronizationJob); err != nil {
		return nil, status, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return &newSynchronizationJob, status, nil
}

// Starts a SynchronizationJob.
func (c *SynchronizationJobClient) Start(ctx context.Context, id string, servicePrincipalId string) (int, error) {
	var status int
	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		ConsistencyFailureFunc: ConflictConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/start", servicePrincipalId, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return status, fmt.Errorf("io.ReadAll(): %v", err)
	}
	return status, nil
}

// Delete
func (c *SynchronizationJobClient) Delete(ctx context.Context, id string, servicePrincipalId string) (int, error) {
	_, status, _, err := c.BaseClient.Delete(ctx, DeleteHttpRequestInput{
		ConsistencyFailureFunc: ConflictConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/", servicePrincipalId, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("SynchronizationJobClient.BaseClient.Delete(): %v", err)
	}

	return status, nil
}

// Pause
func (c *SynchronizationJobClient) Pause(ctx context.Context, id string, servicePrincipalId string) (int, error) {
	var status int
	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		ConsistencyFailureFunc: ConflictConsistencyFailureFunc,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/pause", servicePrincipalId, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	return status, nil
}

func (c *SynchronizationJobClient) Restart(ctx context.Context, id string, synchronizationJobRestartCriteria SynchronizationJobRestartCriteria, servicePrincipalId string) (int, error) {
	var status int

	body, err := json.Marshal(synchronizationJobRestartCriteria)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		ConsistencyFailureFunc: ConflictConsistencyFailureFunc,
		Body:                   body,
		ValidStatusCodes:       []int{http.StatusNoContent},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/restart", servicePrincipalId, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	return status, nil
}

// Provision on demand
func (c *SynchronizationJobClient) ProvisionOnDemand(ctx context.Context, id string, synchronizationJobProvisionOnDemand *SynchronizationJobProvisionOnDemand, servicePrincipalId string) (int, error) {
	var status int

	body, err := json.Marshal(synchronizationJobProvisionOnDemand)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/provisionOnDemand", servicePrincipalId, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	return status, nil
}

// Validate credentials
func (c *SynchronizationJobClient) ValidateCredentials(ctx context.Context, id string, synchronizationJobValidateCredentials *SynchronizationJobValidateCredentials, servicePrincipalId string) (int, error) {
	var status int
	body, err := json.Marshal(synchronizationJobValidateCredentials)
	if err != nil {
		return status, fmt.Errorf("json.Marshal(): %v", err)
	}

	resp, status, _, err := c.BaseClient.Post(ctx, PostHttpRequestInput{
		Body:             body,
		ValidStatusCodes: []int{http.StatusCreated},
		Uri: Uri{
			Entity:      fmt.Sprintf("/servicePrincipals/%s/synchronization/jobs/%s/validateCredentials", servicePrincipalId, id),
			HasTenantId: true,
		},
	})
	if err != nil {
		return status, fmt.Errorf("SynchronizationJobClient.BaseClient.Post(): %v", err)
	}

	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return status, fmt.Errorf("io.ReadAll(): %v", err)
	}

	return status, nil
}
