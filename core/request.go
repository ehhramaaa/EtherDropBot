package core

import (
	"fmt"
)

const apiUrl = "https://api.miniapp.dropstab.com"

// Login
func (c *Client) getToken() (string, error) {
	payload := map[string]string{
		"webAppData": c.account.queryData,
	}

	res, err := c.makeRequest("POST", apiUrl+"/api/auth/login", payload)
	if err != nil {
		return "", err
	}

	if jwt, exits := res["jwt"].(map[string]interface{}); exits {
		if access, exits := jwt["access"].(map[string]interface{}); exits {
			if token, exits := access["token"].(string); exits {
				return token, nil
			} else {
				return "", fmt.Errorf("Token field not found!")
			}
		} else {
			return "", fmt.Errorf("Access field not found!")
		}
	} else {
		return "", fmt.Errorf("JWT field not found!")
	}
}

// User Info
func (c *Client) getUserInfo() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/api/user/current", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Response is nil")
	}
}

// Apply Ref Code
func (c *Client) applyRefCode(refCode string) (bool, error) {
	var payload map[string]string

	if refCode == "" {
		payload = map[string]string{
			"code": "JDLG5",
		}
	} else {
		payload = map[string]string{
			"code": refCode,
		}
	}

	res, err := c.makeRequest("PUT", apiUrl+"/api/user/applyRefLink", payload)
	if err != nil {
		return false, err
	}

	if status, exits := res["usedRefLink"].(bool); exits {
		return status, nil
	} else {
		return false, fmt.Errorf("usedRefLink field not found!")
	}
}

// Welcome Bonus
func (c *Client) welcomeBonus() (map[string]interface{}, error) {
	res, err := c.makeRequest("POST", apiUrl+"/api/bonus/welcomeBonus", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("usedRefLink field not found!")
	}
}

// Daily Bonus
func (c *Client) dailyBonus() (map[string]interface{}, error) {
	res, err := c.makeRequest("POST", apiUrl+"/api/bonus/dailyBonus", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Response is nil")
	}
}

// Ether Drops Subscription
func (c *Client) etherDropsSubscription() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/api/etherDropsSubscription", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Response is nil")
	}
}

// Ref Info
func (c *Client) refInfo() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/api/refLink", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Response is nil")
	}
}

// Claim Ref
func (c *Client) claimRef() (map[string]interface{}, error) {
	res, err := c.makeRequest("POST", apiUrl+"/api/refLink/claim", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Response is nil")
	}
}

// Active Task List
func (c *Client) activeTaskList() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/api/quest/active", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Response is nil")
	}
}

// Verify Task
func (c *Client) verifyTask(taskId int) (string, error) {
	res, err := c.makeRequest("PUT", apiUrl+fmt.Sprintf("/api/quest/%v/verify", taskId), nil)
	if err != nil {
		return "", err
	}

	if status, exits := res["status"].(string); exits {
		return status, nil
	} else {
		return "", fmt.Errorf("Response is nil")
	}
}

// Claim Task
func (c *Client) claimTask(taskId int) (string, error) {
	res, err := c.makeRequest("PUT", apiUrl+fmt.Sprintf("/api/quest/%v/claim", taskId), nil)
	if err != nil {
		return "", err
	}

	if status, exits := res["status"].(string); exits {
		return status, nil
	} else {
		return "", fmt.Errorf("Response is nil")
	}
}

// Completed Task List
func (c *Client) completedTaskList() (map[string]interface{}, error) {
	res, err := c.makeRequest("GET", apiUrl+"/api/quest/completed", nil)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return res, nil
	} else {
		return nil, fmt.Errorf("Response is nil")
	}
}
