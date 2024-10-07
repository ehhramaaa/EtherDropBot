package core

import (
	"EtherDrop/tools"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

func LaunchBot(selectedTools int) {
	sessionsPath := "sessions"
	proxyPath := "configs/proxy.txt"
	maxThread := config.Int("MAX_THREAD")
	isUseProxy := config.Bool("USE_PROXY")

	if !tools.CheckFileOrFolderExits(sessionsPath) {
		os.MkdirAll(sessionsPath, os.ModeDir)
	}

	sessionList, err := tools.ReadFileInDir(sessionsPath)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed To Read File Directory: %v", err))
	}

	if len(sessionList) <= 0 {
		tools.Logger("error", "No Session Found")
		return
	}

	tools.Logger("info", fmt.Sprintf("%v Session Detected", len(sessionList)))

	var wg sync.WaitGroup
	var semaphore chan struct{}
	var proxyList []string

	if isUseProxy {
		proxyList, err = tools.ReadFileTxt(proxyPath)
		if err != nil {
			tools.Logger("error", fmt.Sprintf("Proxy Data Not Found: %s", err))
		}

		tools.Logger("info", fmt.Sprintf("%v Proxy Detected", len(proxyList)))
	}

	if maxThread > len(sessionList) {
		semaphore = make(chan struct{}, len(sessionList))
	} else {
		semaphore = make(chan struct{}, maxThread)
	}

	switch selectedTools {
	case 1:
		for {
			totalPointsChan := make(chan int, len(sessionList))
			for index, session := range sessionList {
				wg.Add(1)
				account := &Account{
					phone: strings.TrimSuffix(session.Name(), ".json"),
				}

				go account.worker(&wg, &semaphore, &totalPointsChan, index, session, proxyList, selectedTools, "")
			}

			go func() {
				wg.Wait()
				close(totalPointsChan)
			}()

			var totalPoints int

			for points := range totalPointsChan {
				totalPoints += points
			}

			tools.Logger("success", fmt.Sprintf("Total Points All Account: %v", totalPoints))

			randomSleep := tools.RandomNumber(config.Int("RANDOM_SLEEP.MIN"), config.Int("RANDOM_SLEEP.MAX"))

			tools.Logger("info", fmt.Sprintf("Launch Bot Finished | Sleep %vs Before Next Lap...", randomSleep))

			time.Sleep(time.Duration(randomSleep) * time.Second)
		}
	case 2:
		if tools.CheckFileOrFolderExits("./ref_code.json") {
			os.Rename("./ref_code.json", "./ref_code_old.json")
			tools.Logger("info", fmt.Sprintf("Rename ref_code.json to ref_code_old.json"))
		}

		for index, session := range sessionList {
			wg.Add(1)
			account := &Account{
				phone: strings.TrimSuffix(session.Name(), ".json"),
			}

			go account.worker(&wg, &semaphore, nil, index, session, proxyList, selectedTools, "")
		}
		wg.Wait()
	case 3:
		refCode := tools.InputTerminal("Input Your Ref Code: ")
		countAccount, _ := strconv.Atoi(tools.InputTerminal(fmt.Sprintf("How Much Account You Want To Register With Ref Code %s: ", refCode)))

		for index, session := range sessionList {
			wg.Add(1)

			if index != 0 && index%countAccount == 0 {
				refCode = tools.InputTerminal("Input Another Ref Code: ")
				countAccount, _ = strconv.Atoi(tools.InputTerminal(fmt.Sprintf("How Much Account You Want To Register With Ref Code %s: ", refCode)))
			}

			account := &Account{
				phone: strings.TrimSuffix(session.Name(), ".json"),
			}

			go account.worker(&wg, &semaphore, nil, index, session, proxyList, selectedTools, "")
		}
		wg.Wait()
	}
}
