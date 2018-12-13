package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceNetwork() *schema.Resource {
        return &schema.Resource{
                Create: resourceNetworkCreate,
                Read:   resourceNetworkRead,
                Update: resourceNetworkUpdate,
                Delete: resourceNetworkDelete,

                Schema: map[string]*schema.Schema{
								    "subnet4": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                    "masklength4": {
                                    Type:     schema.TypeInt,
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

func resourceNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var network = chkp.Network{}
	network.Name = d.Get("name").(string)
	network.Subnet4 = d.Get("subnet4").(string)
  network.Masklength4 = d.Get("masklength4").(int)
  network.Color = d.Get("color").(string)
  networknat := d.Get("natsettings").(*schema.Set).List()
  natconfig := chkp.NatSettings{}
  if len(networknat) != 0 {
  networknatconfig := networknat[0].(map[string]interface{})
  if v := networknatconfig["autorule"]; v != nil {
    autorule := v.(bool)
    natconfig.Autorule = autorule
  }
  if v := networknatconfig["ipaddress"]; v != nil {
    ipaddress := v.(string)
    natconfig.Ipaddress = ipaddress
  }
  if v := networknatconfig["method"]; v != nil {
    method := v.(string)
    natconfig.Method = method
  }
  if v := networknatconfig["installon"]; v != nil {
    installon := v.(string)
    natconfig.Installon = installon
  }
  if v := networknatconfig["hidebehind"]; v != nil {
    hidebehind := v.(string)
    natconfig.Hidebehind = hidebehind
  }
network.NatSettings = natconfig
}
	id, err := client.CreateNetwork(network)
  readNetwork := chkp.Network{}
  json.Unmarshal(id, &readNetwork)
	d.SetId(readNetwork.Uid)
  d.Set("uid", readNetwork.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceNetworkRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	id, err := client.ShowNetwork(d.Id())

	var readnetwork chkp.Network
  json.Unmarshal(id, &readnetwork)
	d.SetId(readnetwork.Uid)
	d.Set("color", readnetwork.Color)
	d.Set("name", readnetwork.Name)
	d.Set("subnet4", readnetwork.Subnet4)
  d.Set("masklength4", readnetwork.Masklength4)
  d.Set("natsettings", flattenNatSettings(readnetwork.NatSettings))
	if err != nil {
		return err
	}
	return nil
}

func resourceNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var network = chkp.Network{}
	network.Subnet4 = d.Get("subnet4").(string)
  network.Masklength4 = d.Get("masklength4").(int)
	network.Color = d.Get("color").(string)
  networknat := d.Get("natsettings").(*schema.Set).List()
  natconfig := chkp.NatSettings{}
  if len(networknat) != 0 {
  networknatconfig := networknat[0].(map[string]interface{})

  if v := networknatconfig["autorule"]; v != nil {
    autorule := v.(bool)
    natconfig.Autorule = autorule
  }
  if v := networknatconfig["ipaddress"]; v != nil {
    ipaddress := v.(string)
    natconfig.Ipaddress = ipaddress
  }
  if v := networknatconfig["method"]; v != nil {
    method := v.(string)
    natconfig.Method = method
  }
  if v := networknatconfig["installon"]; v != nil {
    installon := v.(string)
    natconfig.Installon = installon
  }
  if v := networknatconfig["hidebehind"]; v != nil {
    hidebehind := v.(string)
    natconfig.Hidebehind = hidebehind
  }

  network.NatSettings = natconfig
}

  if d.HasChange("name") {network.Newname = d.Get("name").(string)}
	//network.Newname = d.Get("name").(string)
	network.Uid = d.Get("uid").(string)

	id, err := client.SetNetwork(network)
  readNetwork := chkp.Network{}
  if err := json.Unmarshal(id, &readNetwork); err != nil {return err}
  //json.Unmarshal(id, &readNetwork)
	d.SetId(readNetwork.Uid)
  d.Set("uid", readNetwork.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceNetworkDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteNetwork(d.Id())
	return nil

}

func flattenNatSettings(nat chkp.NatSettings) []interface{} {
	result := make(map[string]interface{})

		result["installon"] = nat.Installon
    result["autorule"] = nat.Autorule
    result["hidebehind"] = nat.Hidebehind
    result["ipaddress"] = nat.Ipaddress
    result["method"] = nat.Method
    if result["autorule"] == false { return nil }
	return []interface{}{result}

}
