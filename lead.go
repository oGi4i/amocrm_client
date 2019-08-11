package amocrm

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	NoTasksLeadFilter         = 1
	InProgressTasksLeadFilter = 2

	ActiveLeadsLeadFilter = 1
)

var (
	leadArrayFields = []string{"tags", "custom_fields"}
)

func (c *ClientInfo) AddLead(lead *LeadAdd) (int, error) {
	if lead.Name == "" {
		return 0, errors.New("name is empty")
	}
	if lead.StatusID == "" {
		return 0, errors.New("statusID is empty")
	}

	url := c.Url + apiUrls["leads"]
	resp, err := c.DoPost(url, &AddLeadRequest{Add: []*LeadAdd{lead}})
	if err != nil {
		return 0, err
	}

	return c.GetResponseID(resp)
}

func (c *ClientInfo) UpdateLead(lead *LeadUpdate) (int, error) {
	if lead.ID == "" {
		return 0, errors.New("ID is empty")
	}
	if lead.UpdatedAt == "" {
		return 0, errors.New("updatedAt is empty")
	}

	url := c.Url + apiUrls["leads"]
	resp, err := c.DoPost(url, &UpdateLeadRequest{Update: []*LeadUpdate{lead}})
	if err != nil {
		return 0, err
	}

	return c.GetResponseID(resp)
}

func (c *ClientInfo) GetLead(reqParams *LeadRequestParams) ([]*Lead, error) {
	addValues := make(map[string]string)
	leads := new(GetLeadResponse)
	var err error

	if len(reqParams.ID) > 0 {
		addValues["id"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.ID)), ","), "[]")
	} else {
		if reqParams.LimitRows != 0 {
			addValues["limit_rows"] = strconv.Itoa(reqParams.LimitRows)
			if reqParams.LimitOffset != 0 {
				addValues["limit_offset"] = strconv.Itoa(reqParams.LimitOffset)
			}
		}
		if reqParams.ResponsibleUserID != 0 {
			addValues["responsible_user_id"] = strconv.Itoa(reqParams.ResponsibleUserID)
		}
		if reqParams.Query != "" {
			addValues["query"] = reqParams.Query
		}
		if len(reqParams.Status) > 0 {
			addValues["status"] = strings.Trim(strings.Join(strings.Fields(fmt.Sprint(reqParams.Status)), ","), "[]")
		}
	}

	url := c.Url + apiUrls["leads"]
	body, err := c.DoGet(url, addValues)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, leads)
	if err != nil {
		// fix bad json serialization, where nil array is serialized as nil object
		stringBody := string(body)
		for _, s := range leadArrayFields {
			stringBody = strings.ReplaceAll(stringBody, "\""+s+"\":{}", "\""+s+"\":[]")
		}
		err = json.Unmarshal([]byte(stringBody), leads)
		if err != nil {
			return nil, err
		}
	}

	if leads.Response != nil {
		return nil, leads.Response
	}

	return leads.Embedded.Items, err
}
