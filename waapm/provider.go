package waapm

import (
	"log"
	"errors"
	"os/exec"
	
        "github.com/hashicorp/terraform/helper/schema"
        "github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
        return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"waapm_account": dataSourceSecret(),
		},
		Schema: map[string]*schema.Schema{
			"waapm_exe": {
				Type: schema.TypeString,
				Optional: true,
				Description: "waapm executable",
			},
		},
		ConfigureFunc: providerConfig,
        }
}

func providerConfig(d *schema.ResourceData) (interface{}, error) {

	waapm_exe := d.Get("waapm_exe").(string)
	if waapm_exe == "" {
		return nil, errors.New("waapm_exe not defined")
	}
	return waapm_exe, nil
}

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
	
		Read: dataSourceSecretRead,

                Schema: map[string]*schema.Schema{
                        "account": &schema.Schema{
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "bastion": &schema.Schema{
			        Type:     schema.TypeString,
			        Optional: true,
                        },
                        "format": &schema.Schema{
			        Type:     schema.TypeString,
			        Optional: true,
                        },
                        "key": &schema.Schema{
			        Type:     schema.TypeString,
			        Optional: true,
                        },
                        "modules": &schema.Schema{
			        Type:     schema.TypeString,
			        Optional: true,
                        },
                        "checkin": &schema.Schema{
			        Type:     schema.TypeBool,
			        Optional: true,
			        Default: true,
                        },
                        "generations": &schema.Schema{
			        Type:     schema.TypeInt,
			        Optional: true,
			        Default: 1,
                        },
                        "shared": &schema.Schema{
			        Type:     schema.TypeBool,
			        Optional: true,
			        Default: false,
                        },
                        "directory": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
                        },
                        "value": &schema.Schema{
			        Type:     schema.TypeString,
				Computed:    true,
				Description: "value of the secret",
				Sensitive: true,
                        },
                },
        }
}

func dataSourceSecretRead(d *schema.ResourceData, m interface{}) error {

	waapm_exe := m.(string)
	log.Printf("[DEBUG] dataSourceSecretRead: waapm_exe = %s\n", waapm_exe)

	account_name := d.Get("account").(string)
	format := d.Get("type").(string)
	account_name := d.Get("account").(string)
	
	log.Printf("[DEBUG] Getting %s for account %s\n", account_type, account_name)
	
	credential, err := exec.Command(waapm_exe, "checkout", "-g", "2", account_name).CombinedOutput()
	log.Printf("[DEBUG] credential=%s\n", credential)
	if err != nil {
		log.Fatal(credential)
	}
	
	d.Set("value", string(credential))
	d.SetId(account_name)

        return nil
}