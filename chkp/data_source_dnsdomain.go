package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func dataSourceDNSDomain() *schema.Resource {
        return &schema.Resource{

                Read:   dataSourceDNSDomainRead,
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
                                Optional: true,
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                            },

                },
        }
}

func dataSourceDNSDomainRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  name := d.Get("name").(string)
  // Call the API to get DNSDomain info
  id, err := client.ShowDNSDomain(name)
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
