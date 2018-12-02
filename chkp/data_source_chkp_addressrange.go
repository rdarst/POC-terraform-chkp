package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func dataSourceAddressRange() *schema.Resource {
        return &schema.Resource{

                Read:   dataSourceAddressRangeRead,

                Schema: map[string]*schema.Schema{
							        	"ipv4addressfirst": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
                        "ipv4addresslast": {
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


func dataSourceAddressRangeRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  name := d.Get("name").(string)

	id, err := client.ReadAddressRangeData(name)

	readAddressRange := chkp.AddressRange{}
  json.Unmarshal(id, &readAddressRange)
	d.SetId(readAddressRange.Uid)
  d.Set("uid", readAddressRange.Uid)
	d.Set("name", readAddressRange.Name)
	d.Set("ipv4addressfirst", readAddressRange.Ipv4addressfirst)
  d.Set("ipv4addresslast", readAddressRange.Ipv4addresslast)
	if err != nil {
		return err
	}
	return nil
}
