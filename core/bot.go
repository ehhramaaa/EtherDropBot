package core

import (
	"EtherDrop/tools"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
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

func (account *Account) worker(wg *sync.WaitGroup, semaphore *chan struct{}, totalPointsChan *chan int, index int, session fs.DirEntry, proxyList []string, selectedTools int, refCode string) {
	defer wg.Done()
	*semaphore <- struct{}{}

	var points int
	var proxy string

	tools.Logger("info", fmt.Sprintf("| %s | Starting Bot...", account.phone))

	setDns(&net.Dialer{})

	client := Client{
		account: *account,
	}

	var queryData string
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	var querySuccess bool

	for i := 0; i < 3; i++ {
		browser := initializeBrowser()

		defer browser.MustClose()

		client.browser = browser

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		go func(ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Panic recovered while getting query data: %v", account.phone, r))
					errChan <- fmt.Errorf("panic: %v", r)
				}
			}()

			select {
			case <-ctx.Done():
				tools.Logger("warning", fmt.Sprintf("| %s | Context cancelled, stopping query attempt.", account.phone))
				return
			default:
				query, err := client.getQueryData(session)
				if err != nil {
					errChan <- err
					return
				}
				resultChan <- query
			}
		}(ctx)

		select {
		case <-ctx.Done():
			tools.Logger("error", fmt.Sprintf("| %s | Timeout during getQueryData | Try to get query data again...", account.phone))
			browser.MustClose()

			time.Sleep(3 * time.Second)

			continue

		case result := <-resultChan:
			if result != "" {
				queryData = result
				querySuccess = true
			} else {
				continue
			}

		case err := <-errChan:
			tools.Logger("error", fmt.Sprintf("| %s | Error while getting query data: %v", account.phone, err))
			browser.MustClose()

			time.Sleep(3 * time.Second)

			continue
		}

		if querySuccess {
			tools.Logger("info", fmt.Sprintf("| %s | Get Query Data Successfully...", account.phone))
			break
		}

		if i == 2 {
			tools.Logger("error", fmt.Sprintf("| %s | Failed get query data after 3 attempts!", account.phone))
			break
		}
	}

	if queryData != "" {
		account.queryData = queryData
	} else {
		return
	}

	account.parsingQueryData()

	if len(proxyList) > 0 {
		proxy = proxyList[index%len(proxyList)]
	}

	client.account = *account
	client.proxy = proxy
	client.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}

	if len(client.proxy) > 0 {
		err := client.setProxy()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to set proxy: %v", account.username, err))
		} else {
			tools.Logger("success", fmt.Sprintf("| %s | Proxy Successfully Set...", account.username))
		}
	}

	infoIp, err := client.checkIp()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to check ip: %v", err))
	}

	if infoIp != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Ip: %s | City: %s | Country: %s | Provider: %s", account.username, infoIp["ip"].(string), infoIp["city"].(string), infoIp["country"].(string), infoIp["org"].(string)))
	}

	switch selectedTools {
	case 1:
		points = client.autoFarming()
		*totalPointsChan <- points
	case 2:
		client.getRefCode()
	case 3:
		client.autoRegisterWithRef(refCode)
	}

	<-*semaphore
}

func (c *Client) getQueryData(session fs.DirEntry) (string, error) {
	defer c.browser.MustClose()

	// Set Local Storage
	sessionsPath := "sessions"

	page := c.browser.MustPage()

	account, err := tools.ReadFileJson(filepath.Join(sessionsPath, session.Name()))
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to read file %s: %v", c.account.phone, session.Name(), err))
	}

	// Membuka halaman kosong terlebih dahulu
	c.navigate(page, "https://web.telegram.org/k/")

	page.MustWaitLoad()
	page.MustWaitNavigation()

	time.Sleep(2 * time.Second)

	// Evaluasi JavaScript untuk menyimpan data ke localStorage
	switch v := account.(type) {
	case []map[string]interface{}:
		// Jika data adalah array of maps
		for _, acc := range v {
			for key, value := range acc {
				page.Eval(fmt.Sprintf(`localStorage.setItem('%s', '%s');`, key, value))
			}
		}
	case map[string]interface{}:
		// Jika data adalah single map
		for key, value := range v {
			page.Eval(fmt.Sprintf(`localStorage.setItem('%s', '%s');`, key, value))
		}
	default:
		tools.Logger("error", fmt.Sprintf("| %s | Failed to Evaluate Local Storage: Unknown Data Type", c.account.phone))
	}

	tools.Logger("success", fmt.Sprintf("| %s | Local storage successfully set | Check Login Status...", c.account.phone))

	page.MustReload()
	page.MustWaitLoad()
	page.MustWaitNavigation()

	time.Sleep(5 * time.Second)

	isSessionExpired := c.checkElement(page, "#auth-pages > div > div.tabs-container.auth-pages__container > div.tabs-tab.page-signQR.active > div > div.input-wrapper > button")

	if isSessionExpired {
		tools.Logger("error", fmt.Sprintf("| %s | Session Expired Or Account Banned, Please Check Your Account...", c.account.phone))

		return "", fmt.Errorf("session expired or account banned")
	}

	tools.Logger("success", fmt.Sprintf("| %s | Login successfully | Sleep 3s Before Navigate...", c.account.phone))

	time.Sleep(3 * time.Second)

	tools.Logger("info", fmt.Sprintf("| %s | Navigating Telegram...", c.account.phone))

	// Search Bot
	c.searchBot(page, "fomo")

	time.Sleep(2 * time.Second)

	// Click Launch App
	c.clickElement(page, "div.new-message-bot-commands")

	c.popupLaunchBot(page)

	time.Sleep(2 * time.Second)

	isIframe := c.checkElement(page, ".payment-verification")

	if !isIframe {
		return "", fmt.Errorf("Failed To Launch Bot: Iframe Not Detected")
	}

	iframe := page.MustElement(".payment-verification")

	iframePage := iframe.MustFrame()

	tools.Logger("info", fmt.Sprintf("| %s | Process Get Query Data...", c.account.phone))

	res, err := iframePage.Evaluate(rod.Eval(`() => {
			let initParams = sessionStorage.getItem("__telegram__initParams");
			if (initParams) {
				let parsedParams = JSON.parse(initParams);
				return parsedParams.tgWebAppData;
			}
		
			initParams = sessionStorage.getItem("telegram-apps/launch-params");
			if (initParams) {
				let parsedParams = JSON.parse(initParams);
				return parsedParams;
			}
		
			return null;
		}`))

	if err != nil {
		return "", err
	}

	var queryData string

	if strings.Contains(res.Value.String(), "tgWebAppData=") {
		queryParamsString, err := tools.GetTextAfterKey(res.Value.String(), "tgWebAppData=")
		if err != nil {
			return "", err
		}

		queryData = queryParamsString
	} else {
		if res.Type == proto.RuntimeRemoteObjectTypeString {
			queryData = res.Value.String()
		} else {
			return "", fmt.Errorf("Get Query Data Failed...")
		}
	}

	return queryData, nil
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
		tools.Logger("info", fmt.Sprintf("| %s | ID: %d |  Balance: %d | Welcome Bonus Claimed: %v", c.account.username, int(userInfo["id"].(float64)), int(userInfo["balance"].(float64)), userInfo["welcomeBonusReceived"].(bool)))

		if len(userInfo["usedRefLinkCode"].(string)) > 0 {
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

		if userInfo["usedRefLinkCode"].(bool) {
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
