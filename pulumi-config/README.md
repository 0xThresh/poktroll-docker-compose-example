# Pulumi Configuration
The code in this folder can be used to create a Poktrolld node on AWS without needing to create it manually, or install all the required dependencies on the instance. 

## Install Pulumi 
Instructions can be found on [Pulumi's website](https://www.pulumi.com/docs/clouds/aws/get-started/begin/). If you're on Mac and have brew, you can run `brew install pulumi/tap/pulumi`.

## Create a Pulumi account
Sign up for a [Pulumi account](https://app.pulumi.com/). You can use a Github account. 

## Configure AWS
I'm assuming you already have an AWS account. Create an Access Key, and export those values into your environment variables if you haven't already set up your `~/.aws/credentials` file. 
```bash
export AWS_ACCESS_KEY_ID="<YOUR_ACCESS_KEY_ID>" 
export AWS_SECRET_ACCESS_KEY="<YOUR_SECRET_ACCESS_KEY>"
```

## Create the Resources
Switch to the `pulumi-config` directory in your terminal and run the commands below - 
```bash
cd pulumi-config
pulumi stack init pokt
pulumi stack select pokt 
pulumi up
```