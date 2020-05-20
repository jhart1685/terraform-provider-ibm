package ibm

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/bluemix-go/api/container/registryv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIBMICCRNamespace() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMICCRNamespaceCreate,
		Read:     resourceIBMICCRNoOp,
		Delete:   resourceIBMICCRNoOp,
		Update:   resourceIBMICCRNoOp,
		Exists:   resourceIBMICCRNamespaceExists,
		Importer: &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Namespace name",
			},
			"account": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Account ID",
			},
		},
	}
}

func resourceIBMICCRNamespaceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	regClient, err := meta.(ClientSession).RegistryAPI()
	if err != nil {
		return false, err
	}
	namespace := d.Get("name").(string)
	account := d.Get("account").(string)
	head := registryv1.NamespaceTargetHeader{
		AccountID: account,
	}
	nspaces, err := regClient.Namespaces().GetNamespaces(head)
	if err != nil {
		return false, fmt.Errorf("Error getting namespaces: %v", err)
	}
	for _, ns := range nspaces {
		if ns == namespace {
			return true, nil
		}
	}
	return false, nil
}

func resourceIBMICCRNoOp(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIBMICCRNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	regClient, err := meta.(ClientSession).RegistryAPI()
	if err != nil {
		return err
	}
	namespace := d.Get("name").(string)
	account := d.Get("account").(string)
	head := registryv1.NamespaceTargetHeader{
		AccountID: account,
	}
	_, err = regClient.Namespaces().AddNamespace(namespace, head)

	if err != nil {
		return fmt.Errorf("Error creating namespace %s, err: %v", namespace, err)
	}
	return nil
}
