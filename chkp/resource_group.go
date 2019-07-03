package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)


func resourceGroup() *schema.Resource {
        return &schema.Resource{
                Create: resourceGroupCreate,
                Read:   resourceGroupRead,
                Update: resourceGroupUpdate,
                Delete: resourceGroupDelete,

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

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var group = chkp.Group{}
	group.Name = d.Get("name").(string)
  group.Color = d.Get("color").(string)

  // Pull in the list of Group Members
  members := d.Get("members").(*schema.Set).List()
  memberlist := make([]string, 0, len(members))
	for _, v := range members {
		val, ok := v.(string)
		if ok && val != "" {
			memberlist = append(memberlist, v.(string))
		}
  }
  group.Members = memberlist

  id, err := client.CreateGroup(group)
  readGroup := chkp.Group{}
  json.Unmarshal(id, &readGroup)
	d.SetId(readGroup.Uid)
  d.Set("uid", readGroup.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  // Call the API to get Group info
  id, err := client.ShowGroup(d.Id())
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
	readGroup := chkp.Group{}
  json.Unmarshal(id, &readGroup)
	d.SetId(readGroup.Uid)
	d.Set("color", readGroup.Color)
	d.Set("name", readGroup.Name)
  var group chkp.GroupMembers
  json.Unmarshal(id, &group)

  // Read the list of Current Group Members
  i := 0
  mem := make([]string, 0, len(group.Members))
  for _, member := range readGroup.Members {
    mem = append(mem, group.Members[i].Name)
    _ = member
    i += 1
  }

  d.Set("members", mem)

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var group = chkp.Group{}
	group.Color = d.Get("color").(string)
  if d.HasChange("name") {group.Newname = d.Get("name").(string)}
	group.Uid = d.Get("uid").(string)
  // Pull in the list of Group Members
  members := d.Get("members").(*schema.Set).List()
  memberlist := make([]string, 0, len(members))
  for _, v := range members {
    val, ok := v.(string)
    if ok && val != "" {
      memberlist = append(memberlist, v.(string))
    }
  }
  group.Members = memberlist

	id, err := client.SetGroup(group)
  readGroup := chkp.GroupMembers{}
  if err := json.Unmarshal(id, &readGroup); err != nil {return err}
  //json.Unmarshal(id, &readGroup)
	d.SetId(readGroup.Uid)
  d.Set("uid", readGroup.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  used, err := client.CheckWhereUsed(d.Id())
  if used > 0 {
    client.WaitUntilNotUsed(d.Id())
  }
	result, err := client.DeleteGroup(d.Id())
  _ = result
  if err != nil {
		return err
	}
  return nil
}
