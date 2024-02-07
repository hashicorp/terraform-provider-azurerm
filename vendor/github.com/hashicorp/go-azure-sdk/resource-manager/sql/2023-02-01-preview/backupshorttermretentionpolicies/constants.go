package backupshorttermretentionpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DiffBackupIntervalInHours int64

const (
	DiffBackupIntervalInHoursOneTwo  DiffBackupIntervalInHours = 12
	DiffBackupIntervalInHoursTwoFour DiffBackupIntervalInHours = 24
)

func PossibleValuesForDiffBackupIntervalInHours() []int64 {
	return []int64{
		int64(DiffBackupIntervalInHoursOneTwo),
		int64(DiffBackupIntervalInHoursTwoFour),
	}
}
