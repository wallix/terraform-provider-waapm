# terraform-provider-waapm
Terraform provider for WALLIX Application-to-Application Password Manager (WAAPM)

The provider configuration should only define the path of the WAAPM executable:

```terraform
provider "waapm" {
	waapm_path = "c:\\Programs (X86)\\WALLIX\\WAAPM\\waapm.exe"
}
```

In order to retrieve a secret for the bastion and to store it in a variable, define a data source:

```terraform
data "waapm_account" "service_account" {
  account = "service_account@acme.net"
}
```
The complete list for datasource arguement is the following:
* account (string): account of the secret using bastion target syntax
* bastion (string): bastion url to to query
* format (string): requested secret format
* key (string): type of requested secret
* modules (string) :use modules for fingerprint
* forced_modules (list of strings): forced modules for fingerprint
* checkin (bool): check account in
* generations (int): number of generations to use for fingerprint
* directory (string:) directory for cred and vault files
* application (string): name of the application

Only account is mandatory. The secret s return as value.