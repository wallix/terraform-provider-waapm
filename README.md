# terraform-provider-waapm
Terraform provider for WALLIX Application-to-Application Password Manager (WAAPM)

## Configuration

The provider configuration should only define the path of the WAAPM executable:

```terraform
provider "waapm" {
	waapm_path = "c:\\Programs (X86)\\WALLIX\\WAAPM\\waapm.exe"
}
```

In order to retrieve a secret for the bastion and store it in a variable, define a datasource:

```terraform
data "waapm_account" "service_account" {
  account = "service_account@acme.net"
}
```
The complete list for datasource arguments is as follows:
* account (string): account of the secret using bastion target name syntax. This argument is mandatory.
* bastion (string): bastion url to request
* format (string): requested secret format
* key (string): type of requested secret
* modules (string): use modules for fingerprint
* forced_modules (list of strings): forced modules for fingerprint
* checkin (bool): check account in
* generations (int): number of generations to use for fingerprint
* directory (string): directory for cred and vault files
* application (string): name of the application

Only account is mandatory. The secret is returned as value.

## Installation

Install WAAPM providers by placing the plugin executables in the user plugins directory.
The user plugins directory is in one of the following locations, depending on the host operating system:

| Operating system  | User plugins directory          |
|-------------------|---------------------------------|
|   Windows         | `%APPDATA%\terraform.d\plugins` |
| All other systems | `~/.terraform.d/plugins`        |

Once a plugin is installed, terraform init can initialize it normally.
