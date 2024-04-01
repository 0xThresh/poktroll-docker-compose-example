package main

import (
	awsEc2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	awsx "github.com/pulumi/pulumi-awsx/sdk/v2/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	// "github.com/joho/godotenv"
	// "fmt"
	// "log"
	"os"
)

func main() {
	// // Attempt to load .env file and exit if it doesn't exist
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// Assign .env values
	// dns := os.Getenv("DNS_HOSTNAME")
	// fmt.Println(dns)
	// pocket_account_pw := os.Getenv("POCKET_ACCOUNT_PASSWORD")
	// fmt.Println(pocket_account_pw)

	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create a new VPC
		vpc, err := awsx.NewVpc(ctx, "vpc", &awsx.VpcArgs{
			CidrBlock: pulumi.StringRef("172.16.0.0/16"),
		})
		if err != nil {
			return err
		}

		// Export VPC properties for use in EC2
		ctx.Export("vpcId", vpc.VpcId)
		ctx.Export("privateSubnetIds", vpc.PrivateSubnetIds)
		ctx.Export("publicSubnetIds", vpc.PublicSubnetIds)

		// Create required security group
		sg, err := awsEc2.NewSecurityGroup(ctx, "pokt-node-sg", &awsEc2.SecurityGroupArgs{
			Description: pulumi.String("Allow HTTP(S), SSH, and Pocket node traffic"),
			VpcId:       vpc.VpcId.ToStringPtrOutput(),
			Ingress: awsEc2.SecurityGroupIngressArray{
				&awsEc2.SecurityGroupIngressArgs{
					FromPort:       pulumi.Int(22),
					ToPort:         pulumi.Int(22),
					Protocol:       pulumi.String("tcp"),
					Description:    pulumi.String("HTTP"),
					CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
				},
				&awsEc2.SecurityGroupIngressArgs{
					FromPort:       pulumi.Int(80),
					ToPort:         pulumi.Int(80),
					Protocol:       pulumi.String("tcp"),
					Description:    pulumi.String("HTTP"),
					CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
				},
				&awsEc2.SecurityGroupIngressArgs{
					FromPort:       pulumi.Int(443),
					ToPort:         pulumi.Int(443),
					Protocol:       pulumi.String("tcp"),
					Description:    pulumi.String("HTTPS"),
					CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
				},
				&awsEc2.SecurityGroupIngressArgs{
					FromPort:       pulumi.Int(8081),
					ToPort:         pulumi.Int(8081),
					Protocol:       pulumi.String("tcp"),
					Description:    pulumi.String("Pocket HTTP API"),
					CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
				},
				&awsEc2.SecurityGroupIngressArgs{
					FromPort:       pulumi.Int(26656),
					ToPort:         pulumi.Int(26656),
					Protocol:       pulumi.String("tcp"),
					Description:    pulumi.String("Pocket RPC API"),
					CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
				}},
			Egress: awsEc2.SecurityGroupEgressArray{
				&awsEc2.SecurityGroupEgressArgs{
					FromPort:       pulumi.Int(0),
					ToPort:         pulumi.Int(0),
					Protocol:       pulumi.String("-1"),
					Description:    pulumi.String("Allow all outbound"),
					CidrBlocks:     pulumi.StringArray{pulumi.String("0.0.0.0/0")},
					Ipv6CidrBlocks: pulumi.StringArray{pulumi.String("::/0")},
				},
			},
		})
		if err != nil {
			return err
		}

		// Build EC2 instance in public subnets
		// Read in user data script
		userdata, err := os.ReadFile("poktrolld-userdata.sh")
		if err != nil {
			return err
		}

		ami, err := awsEc2.LookupAmi(ctx, &awsEc2.LookupAmiArgs{
			Filters: []awsEc2.GetAmiFilter{
				{
					Name:   "name",
					Values: []string{"ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"},
				},
				{
					Name:   "virtualization-type",
					Values: []string{"hvm"},
				},
			},
			Owners:     []string{"099720109477"}, // Canonical
			MostRecent: pulumi.BoolRef(true),
			ExecutableUsers: []string{
				"all",
			},
		})
		if err != nil {
			return err
		}

		// Create EC2 instance
		instance, err := awsEc2.NewInstance(ctx, "pokt-gw", &awsEc2.InstanceArgs{
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			SubnetId:            vpc.PublicSubnetIds.Index(pulumi.Int(0)),
			InstanceType:        pulumi.String("m6a.2xlarge"),
			UserData:            pulumi.String(string(userdata)),
			Ami:                 pulumi.String(ami.Id),
			// TODO: Create the required role in Pulumi
			IamInstanceProfile: pulumi.String("ssm"),
		})
		if err != nil {
			return err
		}

		// Export instance properties
		ctx.Export("instanceId", instance.ID())
		ctx.Export("instancePublicIp", instance.PublicIp)

		return nil
	})
}
