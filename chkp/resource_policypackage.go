package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourcePolicyPackage() *schema.Resource {
        return &schema.Resource{
                Create: resourcePolicyPackageCreate,
                Read:   resourcePolicyPackageRead,
                Update: resourcePolicyPackageUpdate,
                Delete: resourcePolicyPackageDelete,

                Schema: map[string]*schema.Schema{
							        	"access": {
                                Type:     schema.TypeBool,
                                Required: true,
                                  },
                        "desktopsecurity": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: false,
                        },
                        "qos": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: false,
                        },
                        "threatprevention": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: false,
                        },
                        "qospolicytype": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Default: "recommended",
                        },
						            "name": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
						            "color": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Default: "black",
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
                },
        }
}

func resourcePolicyPackageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var policypackage = chkp.PolicyPackage{}
	policypackage.Name = d.Get("name").(string)
	policypackage.Access = d.Get("access").(bool)
  policypackage.DesktopSecurity = d.Get("desktopsecurity").(bool)
  policypackage.Qos = d.Get("qos").(bool)
  policypackage.ThreatPrevention = d.Get("threatprevention").(bool)
  policypackage.QosPolicyType = d.Get("qospolicytype").(string)
  policypackage.Color = d.Get("color").(string)
	id, err := client.CreatePolicyPackage(policypackage)
  //Read in from the API Output
  readPolicyPackage := chkp.PolicyPackage{}
  json.Unmarshal(id, &readPolicyPackage)
	d.SetId(readPolicyPackage.Uid)
  d.Set("uid", readPolicyPackage.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourcePolicyPackageRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
	id, err := client.ShowPolicyPackage(d.Id())
  if err != nil {
    status := err.Error()
    if (status == "404") {
          // If the object is not found remove it from state
          d.SetId("")
          return nil
    } else {
      return err
    }
  }
	readPolicyPackage := chkp.PolicyPackage{}
  json.Unmarshal(id, &readPolicyPackage)
	d.SetId(readPolicyPackage.Uid)
	d.Set("color", readPolicyPackage.Color)
	d.Set("name", readPolicyPackage.Name)
	d.Set("access", readPolicyPackage.Access)
  d.Set("desktopsecurity", readPolicyPackage.DesktopSecurity)
  d.Set("qos", readPolicyPackage.Qos)
  d.Set("threatprevention", readPolicyPackage.ThreatPrevention)
  d.Set("qospolicytype", readPolicyPackage.QosPolicyType)
	
	return nil
}

func resourcePolicyPackageUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var policypackage = chkp.PolicyPackage{}
  policypackage.Access = d.Get("access").(bool)
  policypackage.DesktopSecurity = d.Get("desktopsecurity").(bool)
  policypackage.Qos = d.Get("qos").(bool)
  policypackage.ThreatPrevention = d.Get("threatprevention").(bool)
  policypackage.QosPolicyType = d.Get("qospolicytype").(string)
  policypackage.Color = d.Get("color").(string)
	policypackage.Uid = d.Get("uid").(string)
	id, err := client.SetPolicyPackage(policypackage)
  if err != nil {return err}
  //Read in from the API Output
  readPolicyPackage := chkp.PolicyPackage{}
  if err := json.Unmarshal(id, &readPolicyPackage); err != nil {return err}
  d.SetId(readPolicyPackage.Uid)
  d.Set("uid", readPolicyPackage.Uid)

	if err != nil {
		return err
	}
	return nil
}

func resourcePolicyPackageDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeletePolicyPackage(d.Id())
	return nil

}
