package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceDNSDomain() *schema.Resource {
        return &schema.Resource{
                Create: resourceDNSDomainCreate,
                Read:   resourceDNSDomainRead,
                Update: resourceDNSDomainUpdate,
                Delete: resourceDNSDomainDelete,

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
                        "issubdomain": {
                                Type:     schema.TypeBool,
                                Required: true,
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                            },

                },
        }
}

func resourceDNSDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var dnsdomain = chkp.DNSDomain{}
	dnsdomain.Name = d.Get("name").(string)
  dnsdomain.Color = d.Get("color").(string)
  dnsdomain.Issubdomain = d.Get("issubdomain").(bool)
  id, err := client.CreateDNSDomain(dnsdomain)
  readDNSDomain := chkp.DNSDomain{}
  json.Unmarshal(id, &readDNSDomain)
	d.SetId(readDNSDomain.Uid)
  d.Set("uid", readDNSDomain.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  // Call the API to get DNSDomain info
  id, err := client.ShowDNSDomain(d.Id())
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
	readDNSDomain := chkp.DNSDomain{}
  json.Unmarshal(id, &readDNSDomain)
	d.SetId(readDNSDomain.Uid)
	d.Set("color", readDNSDomain.Color)
	d.Set("name", readDNSDomain.Name)
  d.Set("issubdomain", readDNSDomain.Issubdomain)

	return nil
}

func resourceDNSDomainUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var dnsdomain = chkp.DNSDomain{}
	dnsdomain.Color = d.Get("color").(string)
  if d.HasChange("name") {dnsdomain.Newname = d.Get("name").(string)}
	dnsdomain.Name = d.Get("name").(string)
  dnsdomain.Issubdomain = d.Get("issubdomain").(bool)

	id, err := client.SetDNSDomain(dnsdomain)

  readDNSDomain := chkp.DNSDomain{}
  if err := json.Unmarshal(id, &readDNSDomain); err != nil {return err}
  //json.Unmarshal(id, &readDNSDomain)
	d.SetId(readDNSDomain.Uid)
  d.Set("uid", readDNSDomain.Uid)
  d.Set("issubdomain", readDNSDomain.Issubdomain)
	if err != nil {
		return err
	}
	return nil
}

func resourceDNSDomainDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  used, err := client.CheckWhereUsed(d.Id())
  if used > 0 {
    client.WaitUntilNotUsed(d.Id())
  }
	result, err := client.DeleteDNSDomain(d.Id())
  _ = result
  if err != nil {
		return err
	}
  return nil
}
