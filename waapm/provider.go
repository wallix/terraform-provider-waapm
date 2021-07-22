package waapm

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider implementation for waapm
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"waapm_account": dataSourceSecret(),
		},
		Schema: map[string]*schema.Schema{
			"waapm_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "path of waapm executable",
			},
		},
		ConfigureFunc: providerConfig,
	}
}

func providerConfig(d *schema.ResourceData) (interface{}, error) {

	waapmPath := d.Get("waapm_path").(string)
	if waapmPath == "" {
		return nil, errors.New("waapm_path not defined")
	}
	return waapmPath, nil
}

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceSecretRead,

		Schema: map[string]*schema.Schema{
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "account of the secret using target syntax",
			},
			"bastion": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "bastion to query",
			},
			"format": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "requested secret format",
			},
			"key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "type of requested secret",
			},
			"modules": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "use modules for fingerprint",
			},
			"forced_modules": &schema.Schema{
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "forced modules for fingerprint",
			},
			"checkin": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "check account in",
			},
			"generations": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2,
				Description: "number of generations to use for fingerprint",
			},
			"directory": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "directory for cred and vault files",
			},
			"application": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "terraform",
				Description: "name of the application",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "value of the secret",
				Sensitive:   true,
			},
		},
	}
}

func dataSourceSecretRead(d *schema.ResourceData, m interface{}) error {

	waapmPath, ok := m.(string)
	if !ok {
		return fmt.Errorf("cannot read waapm_path")
	}

	args := []string{"checkout"}

	var value interface{}

	var bastion string
	value, ok = d.GetOk("bastion")
	if ok {
		bastion, ok = value.(string)
		if !ok {
			return fmt.Errorf("cannot read bastion")
		}
		args = append(args, "-b", bastion)

	}

	var format string
	value, ok = d.GetOk("format")
	if ok {
		format, ok = value.(string)
		if !ok {
			return fmt.Errorf("cannot read format")
		}
		args = append(args, "-f", format)
	}

	var key string
	value, ok = d.GetOk("key")
	if ok {
		key, ok = value.(string)
		if !ok {
			return fmt.Errorf("cannot read key")
		}
		args = append(args, "-k", key)
	}

	var modules string
	value, ok = d.GetOk("modules")
	if ok {
		modules, ok = value.(string)
		if !ok {
			return fmt.Errorf("cannot read modules")
		}
		args = append(args, "-m", modules)
	}

	var forcedModules []string
	value, ok = d.GetOk("forced_modules")
	if ok {
		v, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("cannot read forced modules")
		}
		forcedModules = make([]string, len(v))
		for i, m := range v {
			forcedModules[i], ok = m.(string)
			if !ok {
				return fmt.Errorf("cannot read forced module #%d", i)
			}
			args = append(args, "+"+forcedModules[i])
		}
	}

	var checkin bool
	value, ok = d.GetOk("checkin")
	if ok {
		checkin, ok = value.(bool)
		if !ok {
			return fmt.Errorf("cannot read checkin")
		}
		if checkin == false {
			args = append(args, "-n")
		}
	}

	generations := int(2)
	value, ok = d.GetOk("generations")
	if ok {
		generations, ok = value.(int)
		if !ok || generations < 2 {
			return fmt.Errorf("cannot read generations: %v", value)
		} else if generations > 1 {
			args = append(args, "-g", strconv.FormatUint(uint64(generations), 10))
		}
	}

	var application string
	value, ok = d.GetOk("application")
	if ok {
		application, ok = value.(string)
		if !ok {
			return fmt.Errorf("cannot read application")
		}
		args = append(args, "-a", application)
	}

	value, ok = d.GetOk("account")
	if !ok {
		return fmt.Errorf("account is not set")
	}
	accountName := value.(string)
	args = append(args, accountName)

	log.Printf("[DEBUG] command: %s %v\n", waapmPath, args)

	credential, err := exec.Command(waapmPath, args...).CombinedOutput()
	if err != nil {
		fmt.Println(string(credential))
		return fmt.Errorf("%s", credential)
	}

	d.Set("value", string(credential))
	d.SetId(accountName)

	return nil
}
