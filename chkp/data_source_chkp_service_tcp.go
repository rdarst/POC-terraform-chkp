package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func dataSourceServiceTcp() *schema.Resource {
        return &schema.Resource{

                Read:   dataSourceServiceTcpRead,

                Schema: map[string]*schema.Schema{
							        	"port": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
						            "name": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
						            "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
                },
        }
}

func dataSourceServiceTcpRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  name := d.Get("name").(string)

	id, err := client.ReadServiceTcpData(name)

	readServiceTcp := chkp.ServiceTcp{}
  json.Unmarshal(id, &readServiceTcp)
	d.SetId(readServiceTcp.Uid)
  d.Set("uid", readServiceTcp.Uid)
	d.Set("name", readServiceTcp.Name)
	d.Set("port", readServiceTcp.Port)
	if err != nil {
		return err
	}
	return nil
}
