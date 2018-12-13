package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceHost() *schema.Resource {
        return &schema.Resource{
                Create: resourceHostCreate,
                Read:   resourceHostRead,
                Update: resourceHostUpdate,
                Delete: resourceHostDelete,

                Schema: map[string]*schema.Schema{
							        	"ipv4address": {
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

func resourceHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var host = chkp.Host{}
	host.Name = d.Get("name").(string)
	host.Ipv4address = d.Get("ipv4address").(string)
  host.Color = d.Get("color").(string)
  hostnat := d.Get("natsettings").(*schema.Set).List()
  if len(hostnat) != 0 {
  hostnatconfig, ok := hostnat[0].(map[string]interface{})
  if ok {
    natconfig := chkp.NatSettings{}
    if v := hostnatconfig["autorule"]; v != nil {
      autorule := v.(bool)
      natconfig.Autorule = autorule
    }
    if v := hostnatconfig["ipaddress"]; v != nil {
      ipaddress := v.(string)
      natconfig.Ipaddress = ipaddress
    }
    if v := hostnatconfig["method"]; v != nil {
      method := v.(string)
      natconfig.Method = method
    }
    if v := hostnatconfig["installon"]; v != nil {
      installon := v.(string)
      natconfig.Installon = installon
    }
    if v := hostnatconfig["hidebehind"]; v != nil {
      hidebehind := v.(string)
      natconfig.Hidebehind = hidebehind
    }
    host.NatSettings = natconfig
  }
}
	id, err := client.CreateHost(host)
  //Read in from the API Output
  readHost := chkp.Host{}
  json.Unmarshal(id, &readHost)
	d.SetId(readHost.Uid)
  d.Set("uid", readHost.Uid)

	if err != nil {
		return err
	}
	return nil
}

func resourceHostRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
	id, err := client.ShowHost(d.Id())

	readHost := chkp.Host{}
  json.Unmarshal(id, &readHost)
	d.SetId(readHost.Uid)
	d.Set("color", readHost.Color)
	d.Set("name", readHost.Name)
	d.Set("ipv4address", readHost.Ipv4address)
  d.Set("natsettings", flattenHostSettings(readHost.NatSettings))
	if err != nil {
		return err
	}
	return nil
}

func resourceHostUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var host = chkp.Host{}
	host.Ipv4address = d.Get("ipv4address").(string)
	host.Color = d.Get("color").(string)
  hostnat := d.Get("natsettings").(*schema.Set).List()
  if len(hostnat) != 0 {
  hostnatconfig, ok := hostnat[0].(map[string]interface{})
  if ok {
    natconfig := chkp.NatSettings{}
    if v := hostnatconfig["autorule"]; v != nil {
      autorule := v.(bool)
      natconfig.Autorule = autorule
    }
    if v := hostnatconfig["ipaddress"]; v != nil {
      ipaddress := v.(string)
      natconfig.Ipaddress = ipaddress
    }
    if v := hostnatconfig["method"]; v != nil {
      method := v.(string)
      natconfig.Method = method
    }
    if v := hostnatconfig["installon"]; v != nil {
      installon := v.(string)
      natconfig.Installon = installon
    }
    if v := hostnatconfig["hidebehind"]; v != nil {
      hidebehind := v.(string)
      natconfig.Hidebehind = hidebehind
    }
    host.NatSettings = natconfig
  }
}

  //Update name if it has changed
	if d.HasChange("name") {host.Newname = d.Get("name").(string)}
	host.Uid = d.Get("uid").(string)
	id, err := client.SetHost(host)
  if err != nil {return err}
  //Read in from the API Output
  readHost := chkp.Host{}
  if err := json.Unmarshal(id, &readHost); err != nil {return err}
  d.SetId(readHost.Uid)
  d.Set("uid", readHost.Uid)

	if err != nil {
		return err
	}
	return nil
}

func resourceHostDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteHost(d.Id())
	return nil
}

func flattenHostSettings(host chkp.NatSettings) []interface{} {
	result := make(map[string]interface{})


		result["installon"] = host.Installon
    result["autorule"] = host.Autorule
    result["hidebehind"] = host.Hidebehind
    result["ipaddress"] = host.Ipaddress
    result["method"] = host.Method
    if result["autorule"] == false { return nil }
	return []interface{}{result}

}
