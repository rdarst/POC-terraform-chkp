package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceSecurityZone() *schema.Resource {
        return &schema.Resource{
                Create: resourceSecurityZoneCreate,
                Read:   resourceSecurityZoneRead,
                Update: resourceSecurityZoneUpdate,
                Delete: resourceSecurityZoneDelete,

                Schema: map[string]*schema.Schema{

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

func resourceSecurityZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var securityzone = chkp.SecurityZone{}
	securityzone.Name = d.Get("name").(string)
  securityzone.Color = d.Get("color").(string)

  id, err := client.CreateSecurityZone(securityzone)
  readSecurityZone := chkp.SecurityZone{}
  json.Unmarshal(id, &readSecurityZone)
	d.SetId(readSecurityZone.Uid)
  d.Set("uid", readSecurityZone.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceSecurityZoneRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  // Call the API to get SecurityZone info
  id, err := client.ShowSecurityZone(d.Id())
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
	readSecurityZone := chkp.SecurityZone{}
  json.Unmarshal(id, &readSecurityZone)
	d.SetId(readSecurityZone.Uid)
	d.Set("color", readSecurityZone.Color)
	d.Set("name", readSecurityZone.Name)
	return nil
}

func resourceSecurityZoneUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var securityzone = chkp.SecurityZone{}
	securityzone.Color = d.Get("color").(string)
  if d.HasChange("name") {securityzone.Newname = d.Get("name").(string)}
	securityzone.Uid = d.Get("uid").(string)

	id, err := client.SetSecurityZone(securityzone)

  readSecurityZone := chkp.SecurityZone{}
  if err := json.Unmarshal(id, &readSecurityZone); err != nil {return err}
  //json.Unmarshal(id, &readSecurityZone)
	d.SetId(readSecurityZone.Uid)
  d.Set("uid", readSecurityZone.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceSecurityZoneDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  used, err := client.CheckWhereUsed(d.Id())
  if used > 0 {
    client.WaitUntilNotUsed(d.Id())
  }
	result, err := client.DeleteSecurityZone(d.Id())
  _ = result
  if err != nil {
		return err
	}
  return nil
}
