package core

import (
	"EtherDrop/tools"
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

func (account *Account) parsingQueryData() {
	value, err := url.ParseQuery(account.queryData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse query data: %s", err))
	}

	if len(value.Get("query_id")) > 0 {
		account.queryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.authDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.hash = value.Get("hash")
	}

	userParam := value.Get("user")

	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse user data: %s", err))
	}

	userId, ok := userData["id"].(float64)
	if !ok {
		tools.Logger("error", "Failed to convert ID to float64")
	}

	account.userId = int(userId)

	username, ok := userData["username"].(string)
	if !ok {
		tools.Logger("error", "Failed to get username from query")
		return
	}

	account.username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get first name from query")
	}

	account.firstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get last name from query")
	}
	account.lastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		tools.Logger("error", "Failed to get language code from query")
	}
	account.languageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		tools.Logger("error", "Failed to get allows write to pm from query")
	}

	account.allowWriteToPm = allowWriteToPm
}

func (c *Client) autoFarming() int {
	defer tools.RecoverPanic()

	var points int
	var isAlreadyUseRef, isWelcomeBonusReceived bool

	token, err := c.getToken()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get token: %v", c.account.username, err))
		return points
	}

	if token != "" {
		c.accessToken = fmt.Sprintf("Bearer %s", token)
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | Token Not Found!", c.account.username))
		return points
	}

	userInfo, err := c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return points
	}

	if userInfo != nil {
		tools.Logger("info", fmt.Sprintf("| %s | Balance: %d | Available Invite: %d | Welcome Bonus Claimed: %v | Allow Ref: %v", c.account.username, int(userInfo["balance"].(float64)), int(userInfo["availableInvites"].(float64)), userInfo["welcomeBonusReceived"].(bool), userInfo["allowRefLink"].(bool)))

		if userInfo["usedRefLink"].(bool) {
			isAlreadyUseRef = true
		}

		if userInfo["welcomeBonusReceived"].(bool) {
			isWelcomeBonusReceived = true
		}
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | User Info Not Found", c.account.username))
		return points
	}

	dailyBonus, err := c.dailyBonus()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get daily bonus: %v", c.account.username, err))
	}

	if dailyBonus != nil {
		if status, exits := dailyBonus["result"].(bool); exits && status {
			tools.Logger("success", fmt.Sprintf("| %s | Successfully Claim Daily Bonus | Bonus: %d | Streak: %d", c.account.username, int(dailyBonus["bonus"].(float64)), int(dailyBonus["streaks"].(float64))))
		}
	}

	subscription, err := c.etherDropsSubscription()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get subscription: %v", c.account.username, err))
	}

	if subscription != nil {
		tools.Logger("info", fmt.Sprintf("| %s | Balance: %d | Course: %d | Available: %d | Claimed: %d", c.account.username, int(subscription["balance"].(float64)), int(subscription["course"].(float64)), int(subscription["available"].(float64)), len(subscription["claimed"].([]interface{}))))
	}

	if !isAlreadyUseRef {
		applyRef, err := c.applyRefCode("")
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to apply ref code: %v", c.account.username, err))
		}

		if applyRef {
			tools.Logger("success", fmt.Sprintf("| %s | Successfully Apply Ref Code", c.account.username))
		}
	}

	if !isWelcomeBonusReceived {
		welcomeBonus, err := c.welcomeBonus()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to claim welcome bonus: %v", c.account.username, err))
		}

		if welcomeBonus != nil {
			if status, exits := welcomeBonus["result"].(bool); exits && status {
				tools.Logger("success", fmt.Sprintf("| %s | Successfully Claim Welcome Bonus | Bonus: %d", c.account.username, int(welcomeBonus["bonus"].(float64))))
			}
		}
	}

	refInfo, err := c.refInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get ref info: %v", c.account.username, err))
	}

	if refInfo != nil {
		if currentLevel, exits := refInfo["currentLvl"].(map[string]interface{}); exits {
			tools.Logger("info", fmt.Sprintf("| %s | Current Level: %d | Invite Limit: %d", c.account.username, int(currentLevel["lvl"].(float64)), int(currentLevel["invites"].(float64))))
		}

		if nextLevel, exits := refInfo["nextLvl"].(map[string]interface{}); exits {
			tools.Logger("info", fmt.Sprintf("| %s | Next Level: %d | Invite Limit: %d | Required Balance: %d | Remaining: %d", c.account.username, int(nextLevel["lvl"].(float64)), int(nextLevel["invites"].(float64)), int(nextLevel["required"].(float64)), (int(nextLevel["required"].(float64))-int(refInfo["balance"].(float64)))))
		}

		if referrals, exits := refInfo["referrals"].(map[string]interface{}); exits {
			tools.Logger("info", fmt.Sprintf("| %s | Ref Code: %s | Total Ref: %d | Total Reward: %d | Available Claim: %d", c.account.username, refInfo["code"].(string), int(referrals["total"].(float64)), int(refInfo["totalReward"].(float64)), int(refInfo["availableToClaim"].(float64))))
		}

		if availableClaim, exits := refInfo["availableToClaim"].(float64); exits && int(availableClaim) > 0 {
			claimRef, err := c.claimRef()
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Failed to claim ref: %v", c.account.username, err))
			}

			if claimRef != nil {
				if updateAvailableClaim, exits := claimRef["availableToClaim"].(float64); exits && updateAvailableClaim == 0 {
					tools.Logger("success", fmt.Sprintf("| %s | Successfully Claim Ref | Amount Claimed: %d", c.account.username, int(availableClaim)))
				} else {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to claim ref...", c.account.username))
				}
			}
		}
	}

	activeTaskList, err := c.activeTaskList()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get active task list: %v", c.account.username, err))
	}

	if activeTaskList != nil {
		for _, task := range activeTaskList {
			taskMap := task.(map[string]interface{})
			if quests, exits := taskMap["quests"].([]interface{}); exits && len(quests) > 0 {
				for _, quest := range quests {
					questMap := quest.(map[string]interface{})
					questName := questMap["name"].(string)
					questId := int(questMap["id"].(float64))

					if claimAllowed, exits := questMap["claimAllowed"].(bool); exits && !claimAllowed {
						verifyTask, err := c.verifyTask(questId)
						if err != nil {
							tools.Logger("error", fmt.Sprintf("| %s | Failed to verify task %s: %v | Sleep 5s Before Verify Next Task...", c.account.username, questName, err))
						}

						if verifyTask == "OK" {
							tools.Logger("success", fmt.Sprintf("| %s | Successfully Verify Task %s | Sleep 5s Before Verify Next Task...", c.account.username, questName))
						} else {
							tools.Logger("error", fmt.Sprintf("| %s | Failed Verify Task %s | Sleep 5s Before Verify Next Task...", c.account.username, questName))
						}

						time.Sleep(5 * time.Second)
					}
				}
			}
		}
	}

	tools.Logger("info", fmt.Sprintf("| %s | Sleep 60s Before Claiming Task...", c.account.username))
	time.Sleep(60 * time.Second)

	updateActiveTaskList, err := c.activeTaskList()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get update active task list: %v", c.account.username, err))
	}

	if updateActiveTaskList != nil {
		for _, task := range updateActiveTaskList {
			taskMap := task.(map[string]interface{})
			if quests, exits := taskMap["quests"].([]interface{}); exits && len(quests) > 0 {
				for _, quest := range quests {
					questMap := quest.(map[string]interface{})
					questName := questMap["name"].(string)
					questId := int(questMap["id"].(float64))

					if claimAllowed, exits := questMap["claimAllowed"].(bool); exits && claimAllowed {
						claimTask, err := c.claimTask(questId)
						if err != nil {
							tools.Logger("error", fmt.Sprintf("| %s | Failed to claim task %s: %v | Sleep 5s Before Claim Next Task...", c.account.username, questName, err))
						}

						if claimTask == "OK" {
							tools.Logger("success", fmt.Sprintf("| %s | Successfully Claim Task %s | Sleep 5s Before Claim Next Task...", c.account.username, questName))
						} else {
							tools.Logger("error", fmt.Sprintf("| %s | Failed Verify Claim %s | Sleep 5s Before Claim Next Task...", c.account.username, questName))
						}

						time.Sleep(5 * time.Second)
					}
				}
			}
		}
	}

	userInfo, err = c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return points
	}

	if userInfo != nil {
		points = int(userInfo["balance"].(float64))
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | User Info Not Found", c.account.username))
		return points
	}

	return points
}

func (c *Client) getRefCode() {
	defer tools.RecoverPanic()

	var refData map[string]interface{}

	token, err := c.getToken()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get token: %v", c.account.username, err))
		return
	}

	if token != "" {
		c.accessToken = fmt.Sprintf("Bearer %s", token)
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | Token Not Found!", c.account.username))
		return
	}

	userInfo, err := c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return
	}

	refInfo, err := c.refInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get ref info: %v", c.account.username, err))
	}

	if refInfo != nil && userInfo != nil {
		if isAllowRefLink, exits := userInfo["allowRefLink"].(bool); exits && isAllowRefLink {
			if referrals, exits := refInfo["referrals"].(map[string]interface{}); exits {
				refCode := refInfo["code"].(string)
				availableInvite := int(userInfo["availableInvites"].(float64))
				tools.Logger("info", fmt.Sprintf("| %s | Ref Code: %s | Total Ref: %d | Total Reward: %d | Available Claim: %d | Available Invite: %d", c.account.username, refCode, int(referrals["total"].(float64)), int(refInfo["totalReward"].(float64)), int(refInfo["availableToClaim"].(float64)), availableInvite))

				ref := map[string]interface{}{
					"refCode":         refCode,
					"totalRef":        int(referrals["total"].(float64)),
					"availableInvite": availableInvite,
				}

				refData = map[string]interface{}{
					c.account.username: ref,
				}
			}

			if currentLevel, exits := refInfo["currentLvl"].(map[string]interface{}); exits {
				tools.Logger("info", fmt.Sprintf("| %s | Current Level: %d | Invite Limit: %d", c.account.username, int(currentLevel["lvl"].(float64)), int(currentLevel["invites"].(float64))))
			}

			if nextLevel, exits := refInfo["nextLvl"].(map[string]interface{}); exits {
				tools.Logger("info", fmt.Sprintf("| %s | Next Level: %d | Invite Limit: %d | Required Balance: %d | Remaining: %d", c.account.username, int(nextLevel["lvl"].(float64)), int(nextLevel["invites"].(float64)), int(nextLevel["required"].(float64)), (int(nextLevel["required"].(float64))-int(refInfo["balance"].(float64)))))
			}
		} else {
			tools.Logger("error", fmt.Sprintf("| %s | Ref Link Not Allowed!, Increase Your Balance With Auto Farming First!!!", c.account.username))
		}
	}

	if len(refData) > 0 {
		err := tools.AppendFileToJson("./ref_code.json", refData)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to save ref code: %v", c.account.username, err))
			return
		}

		tools.Logger("success", fmt.Sprintf("| %s | Successfully Save Ref Code", c.account.username))
	}
}

func (c *Client) autoRegisterWithRef(refCode string) {
	defer tools.RecoverPanic()
	var isAlreadyUseRef, isWelcomeBonusReceived bool

	token, err := c.getToken()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get token: %v", c.account.username, err))
		return
	}

	if token != "" {
		c.accessToken = fmt.Sprintf("Bearer %s", token)
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | Token Not Found!", c.account.username))
		return
	}

	userInfo, err := c.getUserInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get user info: %v", c.account.username, err))
		return
	}

	if userInfo != nil {
		tools.Logger("info", fmt.Sprintf("| %s | Balance: %d | Available Invite: %d | Welcome Bonus Claimed: %v | Allow Ref: %v", c.account.username, int(userInfo["balance"].(float64)), int(userInfo["availableInvites"].(float64)), userInfo["welcomeBonusReceived"].(bool), userInfo["allowRefLink"].(bool)))

		if userInfo["usedRefLink"].(bool) {
			isAlreadyUseRef = true
		}

		if userInfo["welcomeBonusReceived"].(bool) {
			isWelcomeBonusReceived = true
		}
	} else {
		tools.Logger("error", fmt.Sprintf("| %s | User Info Not Found", c.account.username))
		return
	}

	if !isAlreadyUseRef {
		applyRef, err := c.applyRefCode(refCode)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to apply ref code %s : %v", c.account.username, refCode, err))
		}

		if applyRef {
			tools.Logger("success", fmt.Sprintf("| %s | Successfully Apply Ref Code: %s", c.account.username, refCode))
		} else {
			tools.Logger("error", fmt.Sprintf("| %s | Failed Apply Ref Code: %s", c.account.username, refCode))
		}
	}

	if !isWelcomeBonusReceived {
		welcomeBonus, err := c.welcomeBonus()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to claim welcome bonus: %v", c.account.username, err))
		}

		if welcomeBonus != nil {
			if status, exits := welcomeBonus["result"].(bool); exits && status {
				tools.Logger("success", fmt.Sprintf("| %s | Successfully Claim Welcome Bonus | Bonus: %d", c.account.username, int(welcomeBonus["bonus"].(float64))))
			} else {
				tools.Logger("error", fmt.Sprintf("| %s | Failed Claim Welcome Bonus", c.account.username))
			}
		}
	}
}
