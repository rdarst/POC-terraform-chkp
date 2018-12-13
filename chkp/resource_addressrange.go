package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceAddressRange() *schema.Resource {
        return &schema.Resource{
                Create: resourceAddressRangeCreate,
                Read:   resourceAddressRangeRead,
                Update: resourceAddressRangeUpdate,
                Delete: resourceAddressRangeDelete,

                Schema: map[string]*schema.Schema{
							        	"ipv4addressfirst": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "ipv4addresslast": {
                                Type:     schema.TypeString,
                                Required: true,
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
                        "natsettings": {
                                Type:     schema.TypeSet,
                                Optional: true,
                                Elem: &schema.Resource{
                                  Schema: map[string]*schema.Schema{
                                      "hidebehind": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "ipaddress": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "autorule": {
                                            Type:     schema.TypeBool,
                                            Optional: true,

                                      },
                                      "installon": {
                                            Type:     schema.TypeString,
                                            Optional: true,

                                      },
                                      "method": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                          },
                                        },
                                      },
                                    },
                },
        }
}

func resourceAddressRangeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var addressrange = chkp.AddressRange{}
	addressrange.Name = d.Get("name").(string)
	addressrange.Ipv4addressfirst = d.Get("ipv4addressfirst").(string)
  addressrange.Ipv4addresslast = d.Get("ipv4addresslast").(string)
  addressrange.Color = d.Get("color").(string)
  addressrangenat := d.Get("natsettings").(*schema.Set).List()
  if len(addressrangenat) != 0 {
  addressrangenatconfig, ok := addressrangenat[0].(map[string]interface{})
  if ok {
    natconfig := chkp.NatSettings{}
    if v := addressrangenatconfig["autorule"]; v != nil {
      autorule := v.(bool)
      natconfig.Autorule = autorule
    }
    if v := addressrangenatconfig["ipaddress"]; v != nil {
      ipaddress := v.(string)
      natconfig.Ipaddress = ipaddress
    }
    if v := addressrangenatconfig["method"]; v != nil {
      method := v.(string)
      natconfig.Method = method
    }
    if v := addressrangenatconfig["installon"]; v != nil {
      installon := v.(string)
      natconfig.Installon = installon
    }
    if v := addressrangenatconfig["hidebehind"]; v != nil {
      hidebehind := v.(string)
      natconfig.Hidebehind = hidebehind
    }
    addressrange.NatSettings = natconfig
  }
}
	id, err := client.CreateAddressRange(addressrange)
  //Read in from the API Output
  readAddressRange := chkp.AddressRange{}
  json.Unmarshal(id, &readAddressRange)
	d.SetId(readAddressRange.Uid)
  d.Set("uid", readAddressRange.Uid)

	if err != nil {
		return err
	}
	return nil
}

func resourceAddressRangeRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
	id, err := client.ShowAddressRange(d.Id())

	readAddressRange := chkp.AddressRange{}
  json.Unmarshal(id, &readAddressRange)
	d.SetId(readAddressRange.Uid)
	d.Set("color", readAddressRange.Color)
	d.Set("name", readAddressRange.Name)
	d.Set("ipv4addressfirst", readAddressRange.Ipv4addressfirst)
  d.Set("ipv4addresslast", readAddressRange.Ipv4addresslast)
  d.Set("natsettings", flattenAddressRangeSettings(readAddressRange.NatSettings))
	if err != nil {
		return err
	}
	return nil
}

func resourceAddressRangeUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var addressrange = chkp.AddressRange{}
  addressrange.Ipv4addressfirst = d.Get("ipv4addressfirst").(string)
  addressrange.Ipv4addresslast = d.Get("ipv4addresslast").(string)
	addressrange.Color = d.Get("color").(string)
  addressrangenat := d.Get("natsettings").(*schema.Set).List()
  if len(addressrangenat) != 0 {
  addressrangenatconfig, ok := addressrangenat[0].(map[string]interface{})
  if ok {
    natconfig := chkp.NatSettings{}
    if v := addressrangenatconfig["autorule"]; v != nil {
      autorule := v.(bool)
      natconfig.Autorule = autorule
    }
    if v := addressrangenatconfig["ipaddress"]; v != nil {
      ipaddress := v.(string)
      natconfig.Ipaddress = ipaddress
    }
    if v := addressrangenatconfig["method"]; v != nil {
      method := v.(string)
      natconfig.Method = method
    }
    if v := addressrangenatconfig["installon"]; v != nil {
      installon := v.(string)
      natconfig.Installon = installon
    }
    if v := addressrangenatconfig["hidebehind"]; v != nil {
      hidebehind := v.(string)
      natconfig.Hidebehind = hidebehind
    }
    addressrange.NatSettings = natconfig
  }
}

  //Update name if it has changed
	if d.HasChange("name") {addressrange.Newname = d.Get("name").(string)}
	addressrange.Uid = d.Get("uid").(string)
	id, err := client.SetAddressRange(addressrange)
  if err != nil {return err}
  //Read in from the API Output
  readAddressRange := chkp.AddressRange{}
  if err := json.Unmarshal(id, &readAddressRange); err != nil {return err}
  d.SetId(readAddressRange.Uid)
  d.Set("uid", readAddressRange.Uid)

	if err != nil {
		return err
	}
	return nil
}

func resourceAddressRangeDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteAddressRange(d.Id())
	return nil
}

func flattenAddressRangeSettings(addressrange chkp.NatSettings) []interface{} {
	result := make(map[string]interface{})

		result["installon"] = addressrange.Installon
    result["autorule"] = addressrange.Autorule
    result["hidebehind"] = addressrange.Hidebehind
    result["ipaddress"] = addressrange.Ipaddress
    result["method"] = addressrange.Method
    if result["autorule"] == false { return nil }
	return []interface{}{result}

}
