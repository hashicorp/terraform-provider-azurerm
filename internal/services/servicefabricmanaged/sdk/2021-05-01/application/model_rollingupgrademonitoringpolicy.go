package application

type RollingUpgradeMonitoringPolicy struct {
	FailureAction             FailureAction `json:"failureAction"`
	HealthCheckRetryTimeout   string        `json:"healthCheckRetryTimeout"`
	HealthCheckStableDuration string        `json:"healthCheckStableDuration"`
	HealthCheckWaitDuration   string        `json:"healthCheckWaitDuration"`
	UpgradeDomainTimeout      string        `json:"upgradeDomainTimeout"`
	UpgradeTimeout            string        `json:"upgradeTimeout"`
}
