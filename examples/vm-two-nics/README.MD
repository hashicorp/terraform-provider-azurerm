# Virtual Machine with two Network Interface Cards

deploy a VM with two NICs on two different subnets. one is the Primary named External and holds the Public IP and NSG. Default Route is set to the Primary NIC.

The Second NIC named Internal and does not hold a Public IP nor NSG.

refer to Azure Documentation regarding routing in Guest OS:
Windows VM:
https://docs.microsoft.com/en-us/azure/virtual-machines/windows/multiple-nics#configure-guest-os-for-multiple-nics

Linux VM:
https://docs.microsoft.com/en-us/azure/virtual-machines/linux/multiple-nics#configure-guest-os-for-multiple-nics
