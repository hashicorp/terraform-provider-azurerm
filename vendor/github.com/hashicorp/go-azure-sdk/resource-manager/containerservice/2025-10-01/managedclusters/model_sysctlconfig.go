package managedclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SysctlConfig struct {
	FsAioMaxNr                     *int64  `json:"fsAioMaxNr,omitempty"`
	FsFileMax                      *int64  `json:"fsFileMax,omitempty"`
	FsInotifyMaxUserWatches        *int64  `json:"fsInotifyMaxUserWatches,omitempty"`
	FsNrOpen                       *int64  `json:"fsNrOpen,omitempty"`
	KernelThreadsMax               *int64  `json:"kernelThreadsMax,omitempty"`
	NetCoreNetdevMaxBacklog        *int64  `json:"netCoreNetdevMaxBacklog,omitempty"`
	NetCoreOptmemMax               *int64  `json:"netCoreOptmemMax,omitempty"`
	NetCoreRmemDefault             *int64  `json:"netCoreRmemDefault,omitempty"`
	NetCoreRmemMax                 *int64  `json:"netCoreRmemMax,omitempty"`
	NetCoreSomaxconn               *int64  `json:"netCoreSomaxconn,omitempty"`
	NetCoreWmemDefault             *int64  `json:"netCoreWmemDefault,omitempty"`
	NetCoreWmemMax                 *int64  `json:"netCoreWmemMax,omitempty"`
	NetIPv4IPLocalPortRange        *string `json:"netIpv4IpLocalPortRange,omitempty"`
	NetIPv4NeighDefaultGcThresh1   *int64  `json:"netIpv4NeighDefaultGcThresh1,omitempty"`
	NetIPv4NeighDefaultGcThresh2   *int64  `json:"netIpv4NeighDefaultGcThresh2,omitempty"`
	NetIPv4NeighDefaultGcThresh3   *int64  `json:"netIpv4NeighDefaultGcThresh3,omitempty"`
	NetIPv4TcpFinTimeout           *int64  `json:"netIpv4TcpFinTimeout,omitempty"`
	NetIPv4TcpKeepaliveProbes      *int64  `json:"netIpv4TcpKeepaliveProbes,omitempty"`
	NetIPv4TcpKeepaliveTime        *int64  `json:"netIpv4TcpKeepaliveTime,omitempty"`
	NetIPv4TcpMaxSynBacklog        *int64  `json:"netIpv4TcpMaxSynBacklog,omitempty"`
	NetIPv4TcpMaxTwBuckets         *int64  `json:"netIpv4TcpMaxTwBuckets,omitempty"`
	NetIPv4TcpTwReuse              *bool   `json:"netIpv4TcpTwReuse,omitempty"`
	NetIPv4TcpkeepaliveIntvl       *int64  `json:"netIpv4TcpkeepaliveIntvl,omitempty"`
	NetNetfilterNfConntrackBuckets *int64  `json:"netNetfilterNfConntrackBuckets,omitempty"`
	NetNetfilterNfConntrackMax     *int64  `json:"netNetfilterNfConntrackMax,omitempty"`
	VMMaxMapCount                  *int64  `json:"vmMaxMapCount,omitempty"`
	VMSwappiness                   *int64  `json:"vmSwappiness,omitempty"`
	VMVfsCachePressure             *int64  `json:"vmVfsCachePressure,omitempty"`
}
