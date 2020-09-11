/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apis

import (
	core "k8s.io/api/core/v1"
)

const (
	Finalizer = "licenses.appscode.com"

	// Specifies the path where auth is enabled
	AuthPathKey = "licenses.appscode.com/auth-path"

	// required fields:
	// - Secret.Data["token"] - a vault token
	SecretTypeTokenAuth core.SecretType = "licenses.appscode.com/token"

	// required for SecretTypeTokenAut
	TokenAuthTokenKey = "token"

	// required fields:
	// - Secret.Data["access_key_id"] - aws access key id
	// - Secret.Data["secret_access_key"] - aws access secret key
	//
	// optional fields:
	// - Secret.Annotations["licenses.appscode.com/aws.header-value"] - specifies the header value that required if X-Vault-AWS-IAM-Server-ID Header is set
	// - Secret.Annotations["licenses.appscode.com/auth-path"] - Specifies the path where aws auth is enabled
	SecretTypeAWSAuth core.SecretType = "licenses.appscode.com/aws"

	// required for SecretTypeAWSAuth
	AWSAuthAccessKeyIDKey = "access_key_id"
	// required for SecretTypeAWSAuth
	AWSAuthAccessSecretKey = "secret_access_key"
	// optional for SecretTypeAWSAuth
	AWSAuthSecurityTokenKey = "security_token"

	// Specifies the header value that required if X-Vault-AWS-IAM-Server-ID Header is set
	// optional for annotation for  SecretTypeAWSAuth
	AWSHeaderValueKey = "licenses.appscode.com/aws.header-value"

	// required fields:
	// - Secret.Data["sa.json"] - gcp access secret key
	//
	// optional fields:
	// - Secret.Annotations["licenses.appscode.com/auth-path"] - Specifies the path where gcp auth is enabled
	SecretTypeGCPAuth core.SecretType = "licenses.appscode.com/gcp"
	// required for SecretTypeGCPAuth
	GCPAuthSACredentialJson = "sa.json"

	// - Secret.Data["msiToken"] - azure managed service identity (MSI)  jwt token
	//
	// optional fields:
	// - Secret.Annotations["licenses.appscode.com/azure.subscription-id"] - The subscription ID for the machine that generated the MSI token. This information can be obtained through instance metadata.
	// - Secret.Annotations["licenses.appscode.com/azure.resource-group-name"] - The resource group for the machine that generated the MSI token. This information can be obtained through instance metadata.
	// - Secret.Annotations["licenses.appscode.com/azure.vm-name"] - The virtual machine name for the machine that generated the MSI token. This information can be obtained through instance metadata. If vmss_name is provided, this value is ignored.
	// - Secret.Annotations["licenses.appscode.com/azure.vmss-name"] - The virtual machine scale set name for the machine that generated the MSI token. This information can be obtained through instance metadata.
	SecretTypeAzureAuth = "licenses.appscode.com/azure"

	// required for SecretTypeAzureAuth
	AzureMSIToken = "msiToken"

	// optional for SecretTypeAzureAuth
	AzureSubscriptionId    = "licenses.appscode.com/azure.subscription-id"
	AzureResourceGroupName = "licenses.appscode.com/azure.resource-group-name"
	AzureVmName            = "licenses.appscode.com/azure.vm-name"
	AzureVmssName          = "licenses.appscode.com/azure.vmss-name"
)
