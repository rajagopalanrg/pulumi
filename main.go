package main

import (
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/compute"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/core"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/network"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/storage"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an Azure Resource Group
		c := config.New(ctx, "")

		resourceGroup, err := core.NewResourceGroup(ctx, "testpulumi", &core.ResourceGroupArgs{
			Location: pulumi.String("WestUS"),
			Name:     pulumi.String("testpulumi"),
			Tags: pulumi.StringMap{
				"Created_By": pulumi.String("32943"),
				"Purpose":    pulumi.String("test pulumi"),
				"Start_Date": pulumi.String("16-Nov-2020"),
				"End_Date":   pulumi.String("16-Jan-2020"),
			},
		})
		if err != nil {
			return err
		}

		// Create an Azure resource (Storage Account)
		account, err := storage.NewAccount(ctx, "storage", &storage.AccountArgs{
			ResourceGroupName:      resourceGroup.Name,
			Name:                   pulumi.String(c.Require("storageName")),
			AccountTier:            pulumi.String("Standard"),
			AccountReplicationType: pulumi.String("LRS"),
			Tags: pulumi.StringMap{
				"Created_By": pulumi.String("32943"),
				"Purpose":    pulumi.String("test pulumi"),
				"Start_Date": pulumi.String("16-Nov-2020"),
				"End_Date":   pulumi.String("16-Jan-2020"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("storage account name is ", account.Name)
		// Export the connection string for the storage account

		virtualNetwork, err := network.NewVirtualNetwork(ctx, "vnet", &network.VirtualNetworkArgs{
			ResourceGroupName: resourceGroup.Name,
			Location:          resourceGroup.Location,
			Name:              pulumi.String(c.Require("vnetName")),
			AddressSpaces:     pulumi.StringArray{pulumi.String("10.0.6.0/24")},
			Tags: pulumi.StringMap{
				"Created_By": pulumi.String("32943"),
				"Purpose":    pulumi.String("test pulumi"),
				"Start_Date": pulumi.String("16-Nov-2020"),
				"End_Date":   pulumi.String("16-Jan-2020"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("vnetName", virtualNetwork.Name)

		subnet, err := network.NewSubnet(ctx, "subnet", &network.SubnetArgs{
			Name:               pulumi.String(c.Require("subnetName")),
			ResourceGroupName:  resourceGroup.Name,
			VirtualNetworkName: virtualNetwork.Name,
			AddressPrefixes: pulumi.StringArray{
				pulumi.String("10.0.6.0/24"),
			},
		})
		if err != nil {
			return err
		}
		networkInterface, err := network.NewNetworkInterface(ctx, "networkinterface", &network.NetworkInterfaceArgs{
			Name:              pulumi.String(c.Require("nicName")),
			Location:          resourceGroup.Location,
			ResourceGroupName: resourceGroup.Name,
			IpConfigurations: network.NetworkInterfaceIpConfigurationArray{
				&network.NetworkInterfaceIpConfigurationArgs{
					Name:                       pulumi.String(c.Require("ipconfig")),
					Primary:                    pulumi.Bool(true),
					PrivateIpAddressAllocation: pulumi.String("Dynamic"),
					SubnetId:                   subnet.ID(),
				},
			},
		})
		if err != nil {
			return err
		}
		/* datadisk := compute.VirtualMachineStorageDataDiskArgs{
			DiskSizeGb: pulumi.Int(60),
			Name:       pulumi.String("datadisk"),
		} */

		virtualMachine, err := compute.NewVirtualMachine(ctx, "virtualmachine", &compute.VirtualMachineArgs{
			Name:              pulumi.String(c.Require("vmName")),
			ResourceGroupName: resourceGroup.Name,
			StorageOsDisk: compute.VirtualMachineStorageOsDiskArgs{
				Name:            pulumi.String("myosdisk1"),
				Caching:         pulumi.String("ReadWrite"),
				CreateOption:    pulumi.String("FromImage"),
				ManagedDiskType: pulumi.String("Standard_LRS"),
			},
			NetworkInterfaceIds: pulumi.StringArray{networkInterface.ID()},
			VmSize:              pulumi.String("Standard_DS1_v2"),
			StorageDataDisks: compute.VirtualMachineStorageDataDiskArray{
				&compute.VirtualMachineStorageDataDiskArgs{
					DiskSizeGb:   pulumi.Int(60),
					Name:         pulumi.String("datadisk"),
					CreateOption: pulumi.String("Empty"),
					Lun:          pulumi.Int(0),
				},
			},
			StorageImageReference: &compute.VirtualMachineStorageImageReferenceArgs{
				Publisher: pulumi.String("Canonical"),
				Offer:     pulumi.String("UbuntuServer"),
				Sku:       pulumi.String("16.04-LTS"),
				Version:   pulumi.String("latest"),
			},
			OsProfile: &compute.VirtualMachineOsProfileArgs{
				ComputerName:  pulumi.String("hostname"),
				AdminUsername: pulumi.String("testadmin"),
				AdminPassword: pulumi.String("Password1234!"),
			},
			OsProfileLinuxConfig: &compute.VirtualMachineOsProfileLinuxConfigArgs{
				DisablePasswordAuthentication: pulumi.Bool(false),
			},
			BootDiagnostics: compute.VirtualMachineBootDiagnosticsArgs{
				Enabled:    pulumi.Bool(true),
				StorageUri: account.PrimaryConnectionString,
			},
			Tags: pulumi.StringMap{
				"Created_By": pulumi.String("32943"),
				"Purpose":    pulumi.String("test pulumi"),
				"Start_Date": pulumi.String("16-Nov-2020"),
				"End_Date":   pulumi.String("16-Jan-2020"),
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("virtual machine name", virtualMachine.Name)
		return nil
	})
}
