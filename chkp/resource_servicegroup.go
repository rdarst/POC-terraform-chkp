package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceServiceGroup() *schema.Resource {
        return &schema.Resource{
                Create: resourceServiceGroupCreate,
                Read:   resourceServiceGroupRead,
                Update: resourceServiceGroupUpdate,
                Delete: resourceServiceGroupDelete,

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
                        "members": &schema.Schema{
                                Type:     schema.TypeSet,
                                Optional: true,
                                Elem: &schema.Schema{Type: schema.TypeString},
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                            },

                },
        }
}

func resourceServiceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var group = chkp.ServiceGroup{}
	group.Name = d.Get("name").(string)
  group.Color = d.Get("color").(string)

  // Pull in the list of ServiceGroup Members
  members := d.Get("members").(*schema.Set).List()
  memberlist := make([]string, 0, len(members))
	for _, v := range members {
		val, ok := v.(string)
		if ok && val != "" {
			memberlist = append(memberlist, v.(string))
		}
  }
  group.Members = memberlist

  id, err := client.CreateServiceGroup(group)
  readServiceGroup := chkp.ServiceGroup{}
  json.Unmarshal(id, &readServiceGroup)
	d.SetId(readServiceGroup.Uid)
  d.Set("uid", readServiceGroup.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceServiceGroupRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
  // Call the API to get ServiceGroup info
  id, err := client.ShowServiceGroup(d.Id())
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
	readServiceGroup := chkp.ServiceGroup{}
  json.Unmarshal(id, &readServiceGroup)
	d.SetId(readServiceGroup.Uid)
	d.Set("color", readServiceGroup.Color)
	d.Set("name", readServiceGroup.Name)
  var group chkp.ServiceGroupMembers
  json.Unmarshal(id, &group)

  // Read the list of Current ServiceGroup Members
  i := 0
  mem := make([]string, 0, len(group.Members))
  for _, member := range readServiceGroup.Members {
    mem = append(mem, group.Members[i].Name)
    _ = member
    i += 1
  }

  d.Set("members", mem)

	return nil
}

func resourceServiceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var group = chkp.ServiceGroup{}
	group.Color = d.Get("color").(string)
  if d.HasChange("name") {group.Newname = d.Get("name").(string)}
	group.Uid = d.Get("uid").(string)
  // Pull in the list of ServiceGroup Members
  members := d.Get("members").(*schema.Set).List()
  memberlist := make([]string, 0, len(members))
  for _, v := range members {
    val, ok := v.(string)
    if ok && val != "" {
      memberlist = append(memberlist, v.(string))
    }
  }
  group.Members = memberlist

	id, err := client.SetServiceGroup(group)
  readServiceGroup := chkp.ServiceGroupMembers{}
  if err := json.Unmarshal(id, &readServiceGroup); err != nil {return err}
  //json.Unmarshal(id, &readServiceGroup)
	d.SetId(readServiceGroup.Uid)
  d.Set("uid", readServiceGroup.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceServiceGroupDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  used, err := client.CheckWhereUsed(d.Id())
  if used > 0 {
    client.WaitUntilNotUsed(d.Id())
  }
  result, err := client.DeleteServiceGroup(d.Id())
  _ = result
  if err != nil {
    return err
  }
  return nil
}
