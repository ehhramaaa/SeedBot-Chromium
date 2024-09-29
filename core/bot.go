package core

import (
	"SeedBot/tools"
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
	"github.com/gookit/config/v2"
)

func (account *Account) parsingQueryData() {
	value, err := url.ParseQuery(account.QueryData)
	if err != nil {
		tools.Logger("error", fmt.Sprintf("Failed to parse query data: %s", err))
	}

	if len(value.Get("query_id")) > 0 {
		account.QueryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.AuthDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.Hash = value.Get("hash")
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

	account.UserId = int(userId)

	username, ok := userData["username"].(string)
	if !ok {
		tools.Logger("error", "Failed to get username from query")
		return
	}

	account.Username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get first name from query")
	}

	account.FirstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		tools.Logger("error", "Failed to get last name from query")
	}
	account.LastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		tools.Logger("error", "Failed to get language code from query")
	}
	account.LanguageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		tools.Logger("error", "Failed to get allows write to pm from query")
	}

	account.AllowWriteToPm = allowWriteToPm
}

func (account *Account) worker(wg *sync.WaitGroup, semaphore *chan struct{}, totalPointsChan *chan float64, index int, session fs.DirEntry, proxyList []string) {
	defer wg.Done()
	*semaphore <- struct{}{}
	defer func() {
		<-*semaphore
	}()

	var points float64

	tools.Logger("info", fmt.Sprintf("| %s | Starting Bot...", account.Phone))

	setDns(&net.Dialer{})

	client := Client{
		Account: *account,
	}

	var queryData string
	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	var querySuccess bool

	for i := 0; i < 3; i++ {
		browser := initializeBrowser()

		defer browser.MustClose()

		client.Browser = browser

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		go func(ctx context.Context) {
			defer func() {
				if r := recover(); r != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Panic recovered while getting query data: %v", account.Phone, r))
					errChan <- fmt.Errorf("panic: %v", r)
				}
			}()

			select {
			case <-ctx.Done():
				tools.Logger("warning", fmt.Sprintf("| %s | Context cancelled, stopping query attempt.", account.Phone))
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
			tools.Logger("error", fmt.Sprintf("| %s | Timeout during getQueryData | Try to get query data again...", account.Phone))
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
			tools.Logger("error", fmt.Sprintf("| %s | Error while getting query data: %v", account.Phone, err))
			browser.MustClose()

			time.Sleep(3 * time.Second)

			continue
		}

		if querySuccess {
			tools.Logger("info", fmt.Sprintf("| %s | Get Query Data Successfully...", account.Phone))
			break
		}

		if i == 2 {
			tools.Logger("error", fmt.Sprintf("| %s | Failed get query data after 3 attempts!", account.Phone))
			break
		}
	}

	if queryData != "" {
		account.QueryData = queryData
	} else {
		return
	}

	account.parsingQueryData()

	if len(proxyList) > 0 {
		client.Proxy = proxyList[index%len(proxyList)]
	}

	client.Account = *account

	client.HttpClient = &http.Client{
		Timeout: 15,
	}

	if len(client.Proxy) > 0 {
		err := client.setProxy()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to set proxy: %v", account.Username, err))
		} else {
			tools.Logger("success", fmt.Sprintf("| %s | Proxy Successfully Set...", account.Username))
		}
	}

	infoIp, err := client.checkIp()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check ip: %v", account.Username, err))
	}

	if infoIp != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Ip: %s | City: %s | Country: %s | Provider: %s", account.Username, infoIp["ip"].(string), infoIp["city"].(string), infoIp["country"].(string), infoIp["org"].(string)))
	}

	points = client.autoCompleteTask(queryData)

	defer func() {
		*totalPointsChan <- points
	}()
}

func (c *Client) checkIp() (map[string]interface{}, error) {
	result, err := c.makeRequest("GET", fmt.Sprintf("https://ipinfo.io/json?token=%v", config.String("IPINFO_TOKEN")), nil)
	if err != nil {
		return nil, fmt.Errorf("error: %v", err)
	}

	return result, nil
}

func (c *Client) getQueryData(session fs.DirEntry) (string, error) {
	defer c.Browser.MustClose()

	// Set Local Storage
	sessionsPath := "sessions"

	page := c.Browser.MustPage()

	account, err := tools.ReadFileJson(filepath.Join(sessionsPath, session.Name()))
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to read file %s: %v", c.Account.Phone, session.Name(), err))
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
		tools.Logger("error", fmt.Sprintf("| %s | Failed to Evaluate Local Storage: Unknown Data Type", c.Account.Phone))
	}

	tools.Logger("success", fmt.Sprintf("| %s | Local storage successfully set | Check Login Status...", c.Account.Phone))

	page.MustReload()
	page.MustWaitLoad()
	page.MustWaitNavigation()

	time.Sleep(5 * time.Second)

	isSessionExpired := c.checkElement(page, "#auth-pages > div > div.tabs-container.auth-pages__container > div.tabs-tab.page-signQR.active > div > div.input-wrapper > button")

	if isSessionExpired {
		tools.Logger("error", fmt.Sprintf("| %s | Session Expired Or Account Banned, Please Check Your Account...", c.Account.Phone))

		return "", fmt.Errorf("session expired or account banned")
	}

	tools.Logger("success", fmt.Sprintf("| %s | Login successfully | Sleep 3s Before Navigate...", c.Account.Phone))

	time.Sleep(3 * time.Second)

	tools.Logger("info", fmt.Sprintf("| %s | Navigating Telegram...", c.Account.Phone))

	// Search Bot
	c.searchBot(page, "seed_coin_bot")

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

	tools.Logger("info", fmt.Sprintf("| %s | Process Get Query Data...", c.Account.Phone))

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

func (c *Client) autoCompleteTask(query string) float64 {
	var isClaimFirstEgg, isClaimFarmingSeed, isWalletConnected bool
	var speedSeedLevel, storageSeedLevel, holyWaterLevel int

	isAutoHatchEgg := config.Bool("AUTO_HATCH_EGG")
	isAutoFeedBird := config.Bool("AUTO_FEED_BIRD")
	isAutoBirdHunt := config.Bool("AUTO_BIRD_HUNT")
	isAutoPlaySpinEgg := config.Bool("AUTO_PLAY_SPIN_EGG")
	isAutoUpgradeSpeed := config.Bool("AUTO_UPGRADE.SPEED")
	isAutoUpgradeStorage := config.Bool("AUTO_UPGRADE.STORAGE")
	isAutoUpgradeHolyWater := config.Bool("AUTO_UPGRADE.HOLY_WATER")

	maxSpeedLevel := config.Int("AUTO_UPGRADE.MAX_LEVEL.SPEED")
	maxStorageLevel := config.Int("AUTO_UPGRADE.MAX_LEVEL.STORAGE")
	maxHolyWaterLevel := config.Int("AUTO_UPGRADE.MAX_LEVEL.HOLY_WATER")

	claimFarmingSeedAfter := tools.RandomNumber(config.Int("CLAIM_FARMING_SEED_AFTER.MIN"), config.Int("CLAIM_FARMING_SEED_AFTER.MAX"))

	c.AccessToken = query

	profileInfo, err := c.getProfile()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed Get Profile: %v", c.Account.Username, err))
		return 0
	}

	if !profileInfo["give_first_egg"].(bool) {
		isClaimFirstEgg = true
	}

	lastClaimFarming, err := time.Parse(time.RFC3339Nano, profileInfo["last_claim"].(string))
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Check Last Claim Time Seed Failed | Failed to parse time: %v", c.Account.Username, err))
	}

	if (lastClaimFarming.Unix() + int64(claimFarmingSeedAfter)) < time.Now().Unix() {
		isClaimFarmingSeed = true
	}

	if upgradesInfo, exits := profileInfo["upgrades"].([]interface{}); exits && len(upgradesInfo) > 0 {
		for _, info := range upgradesInfo {
			infoMap := info.(map[string]interface{})
			if infoMap["upgrade_type"] == "mining-speed" {
				speedSeedLevel = int(infoMap["upgrade_level"].(float64))
			}
			if infoMap["upgrade_type"] == "storage-size" {
				storageSeedLevel = int(infoMap["upgrade_level"].(float64))
			}
			if infoMap["upgrade_type"] == "holy-water" {
				holyWaterLevel = int(infoMap["upgrade_level"].(float64))
			}
		}
	}

	if walletConnected, exits := profileInfo["wallet_connected"].(string); exits && len(walletConnected) > 0 {
		isWalletConnected = true
	} else {
		isWalletConnected = false
	}

	tools.Logger("success", fmt.Sprintf("| %s | Claimed First Egg: %v | Speed Seed Level: %v | Storage Seed Level: %v | Holy Water Level: %v | Wallet Connected: %v", c.Account.Username, isClaimFirstEgg, speedSeedLevel, storageSeedLevel, holyWaterLevel, isWalletConnected))

	balance, err := c.getBalance()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get balance: %v", c.Account.Username, err))
	}

	if balance != 0 {
		tools.Logger("success", fmt.Sprintf("| %s | Balance: %.9f", c.Account.Username, (balance/1e9)))
	}

	guildId, err := c.checkGuild()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check guild: %v", c.Account.Username, err))
		c.joinGuild()
	}

	if guildId != "" && guildId != "9e02254f-d921-43d3-839f-903706dedeb5" {
		c.leaveGuild()
		time.Sleep(3 * time.Second)
		c.joinGuild()
		time.Sleep(3 * time.Second)
	}

	guildInfo, err := c.getGuildInfo()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get guild info: %v", c.Account.Username, err))
	}

	if guildInfo != nil {
		tools.Logger("success", fmt.Sprintf("| %s | %s | Members: %v | Hunted: %.9f | Reward: %.9f | Rank: %v", c.Account.Username, guildInfo["name"].(string), int(guildInfo["number_member"].(float64)), (guildInfo["hunted"].(float64)/1e9), (guildInfo["reward"].(float64)/1e9), int(guildInfo["rank_index"].(float64))))
	}

	birdInventory, err := c.getBirdInventory()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get bird inventory: %v", c.Account.Username, err))
	}

	wormInventory, err := c.getWormInventory()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get worm inventory: %v", c.Account.Username, err))
	}

	eggInventory, err := c.getEggInventory()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get egg inventory: %v", c.Account.Username, err))
	}

	if birdInventory != nil && wormInventory != nil && eggInventory != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Bird Inventory: %v | Worm Inventory: %v | Egg Inventory: %v", c.Account.Username, int(birdInventory["total"].(float64)), int(wormInventory["total"].(float64)), int(eggInventory["total"].(float64))))
	}

	loginBonus, err := c.checkLoginBonus()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check login bonus: %v", c.Account.Username, err))
	}

	if loginBonus != nil && len(loginBonus) > 0 {
		for _, detail := range loginBonus {
			detailMap := detail.(map[string]interface{})
			parsedTime, err := time.Parse(time.RFC3339Nano, detailMap["timestamp"].(string))
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Claim Login Bonus Failed | Failed to parse time: %v", c.Account.Username, err))
				continue
			}

			if parsedTime.Format("2006/01/02") != time.Now().Format("2006/01/02") {
				claimLoginBonus, err := c.claimLoginBonus()
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to claim login bonus: %v", c.Account.Username, err))
				}

				if claimLoginBonus != nil {
					timestamp, err := time.Parse(time.RFC3339Nano, claimLoginBonus["timestamp"].(string))
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Claim Login Bonus Failed | Failed to parse time: %v", c.Account.Username, err))
						continue
					}

					if timestamp.Format("2006/01/02") == time.Now().Format("2006/01/02") {
						tools.Logger("success", fmt.Sprintf("| %s | Claim Login Bonus Successfully | Amount: %.9f", c.Account.Username, (claimLoginBonus["amount"].(float64)/1e9)))
					}
				}
			}
		}
	}

	if isClaimFarmingSeed {
		claimSeed, err := c.claimFarmingSeed()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to claim farming seed: %v", c.Account.Username, err))
		}

		if claimSeed != nil {
			tools.Logger("success", fmt.Sprintf("| %s | Claim Seed Successfully | Amount: %.9f", c.Account.Username, (claimSeed["amount"].(float64)/1e9)))
		}
	}

	if isClaimFirstEgg {
		claimFirstEgg, err := c.claimFirstEgg()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to claim first egg: %v", c.Account.Username, err))
		}

		if claimFirstEgg != nil {
			tools.Logger("success", fmt.Sprintf("| %s | Claim First Egg Successfully | Egg Type: %s | Status: %s", c.Account.Username, claimFirstEgg["type"].(string), claimFirstEgg["status"].(string)))
		}
	}

	if isAutoHatchEgg {
		eggInventory, err = c.getEggInventory()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to check egg inventory: %v", c.Account.Username, err))
		}

		if eggInventory != nil {
			if int(eggInventory["total"].(float64)) > 0 {
				egg := eggInventory["items"].([]interface{})
				for _, item := range egg {
					itemMap := item.(map[string]interface{})
					hatchEgg, err := c.hatchEgg(itemMap["id"].(string))
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Failed to hatch egg: %v", c.Account.Username, err))
					}

					if hatchEgg != nil {
						tools.Logger("success", fmt.Sprintf("| %s | Hatch Egg Successfully | Bird Type: %s", c.Account.Username, hatchEgg["type"].(string)))
					}
				}
			}
		}
	}

	birdStatus, err := c.getBirdStatus()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check bird status: %v", c.Account.Username, err))
	}

	if birdStatus != nil {
		tools.Logger("success", fmt.Sprintf("| %s | Bird Status: %s | Type: %s | Energy Level: %v | Energy Max: %v | Happiness: %v | Task Level: %v", c.Account.Username, birdStatus["status"].(string), birdStatus["type"].(string), int(birdStatus["energy_level"].(float64)/1e9), int(birdStatus["energy_max"].(float64)/1e9), int(birdStatus["happiness_level"].(float64)/1e9), int(birdStatus["task_level"].(float64)/1e9)))
	}

	wormStatus, err := c.getWormStatus()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to check worm status: %v", c.Account.Username, err))
	}

	if wormStatus != nil {
		if !wormStatus["is_caught"].(bool) {
			catchWorm, err := c.catchWorm()
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Failed to catch worm: %v", c.Account.Username, err))
			}

			if catchWorm != nil {
				tools.Logger("success", fmt.Sprintf("| %s | Catch Worm Successfully | Worm Type: %s | Catch Status: %s", c.Account.Username, catchWorm["type"].(string), catchWorm["status"].(string)))
			} else {
				tools.Logger("warning", fmt.Sprintf("| %s | Catch Worm Failed...", c.Account.Username))
			}
		} else {
			tools.Logger("info", fmt.Sprintf("| %s | Catch Worm After: %v", c.Account.Username, wormStatus["ended_at"]))
		}
	}

	if isAutoFeedBird {
		birdInventory, err = c.getBirdInventory()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to get bird inventory: %v", c.Account.Username, err))
		}

		wormInventory, err = c.getWormInventory()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to get worm inventory: %v", c.Account.Username, err))
		}

		if birdInventory != nil && wormInventory != nil {
			if int(birdInventory["total"].(float64)) > 0 {
				if int(wormInventory["total"].(float64)) > 0 {
					birds := birdInventory["items"].([]interface{})
					for _, bird := range birds {
						birdMap := bird.(map[string]interface{})
						if birdMap["is_leader"].(bool) {
							worms := wormInventory["items"].([]interface{})
							for _, worm := range worms {
								wormMap := worm.(map[string]interface{})
								feedBird, err := c.feedBird(birdMap["id"].(string), wormMap["id"].(string))
								if err != nil {
									tools.Logger("error", fmt.Sprintf("| %s | Failed to feed bird: %v", c.Account.Username, err))
								}

								if feedBird != nil {
									if feedBird != nil {
										tools.Logger("success", fmt.Sprintf("| %s | Feed Bird Successfully | Current Energy: %v | Max Energy: %v", c.Account.Username, int(feedBird["energy_level"].(float64)/1e9), int(feedBird["energy_max"].(float64)/1e9)))
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if isAutoBirdHunt {
		birdStatus, err = c.getBirdStatus()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to check bird status: %v", c.Account.Username, err))
		}

		if birdStatus != nil {
			if (int(birdStatus["happiness_level"].(float64)) / 100) < 100 {
				birdHappiness, err := c.birdHappiness(birdStatus["id"].(string), 100)
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Bird Happiness Failed: %v", c.Account.Username, err))
				}

				if birdHappiness != nil {
					tools.Logger("success", fmt.Sprintf("| %s | Bird Happiness Successfully | Current Happiness: %v", c.Account.Username, (int(birdHappiness["happiness_level"].(float64))/100)))
				}

			}

			huntEnd, err := time.Parse(time.RFC3339, birdStatus["hunt_end_at"].(string))
			if err != nil {
				tools.Logger("error", fmt.Sprintf("| %s | Bird Hunt Failed | Failed to parse time: %v", c.Account.Username, err))
			}

			if huntEnd.Unix() < time.Now().Unix() {
				if birdStatus["hunt_end_at"].(string) != "0001-01-01T00:00:00Z" {
					claimBirdHunt, err := c.claimBirdHunt(birdStatus["id"].(string))
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Claim Bird Hunt Failed: %v", c.Account.Username, err))
					}

					if claimBirdHunt != nil {
						tools.Logger("success", fmt.Sprintf("| %s | Claim Bird Hunt Successfully | Sleep 3s Before Start Hunt", c.Account.Username))
					}

					time.Sleep(3 * time.Second)
				}

				if birdStatus["energy_level"].(float64) > 0 {
					startBirdHunt, err := c.startBirdHunt(birdStatus["id"].(string), int(birdStatus["task_level"].(float64)))
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Start Bird Hunt Failed: %v", c.Account.Username, err))
					}

					if startBirdHunt != nil {
						huntEnd, err := time.Parse(time.RFC3339, startBirdHunt["hunt_end_at"].(string))
						if err != nil {
							tools.Logger("success", fmt.Sprintf("| %s | Start Bird Hunt Successfully | Claim Reward After: %v", c.Account.Username, startBirdHunt["hunt_end_at"].(string)))
						} else {
							tools.Logger("success", fmt.Sprintf("| %s | Start Bird Hunt Successfully | Claim Reward After: %vs", c.Account.Username, (huntEnd.Unix()-time.Now().Unix())))
						}
					}
				} else {
					tools.Logger("info", fmt.Sprintf("| %s | Failed Start Bird Hunt | Energy Bird Not Enough", c.Account.Username))
				}
			} else {
				tools.Logger("info", fmt.Sprintf("| %s | Bird Still Hunting | Claim After %vs", c.Account.Username, (huntEnd.Unix()-time.Now().Unix())))
			}
		}
	}

	streakReward, err := c.getStreakReward()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get streak reward: %v", c.Account.Username, err))
	}

	if streakReward != nil {
		if len(streakReward) > 0 {
			for _, streak := range streakReward {
				streakMap := streak.(map[string]interface{})
				claimStreakReward, err := c.claimStreakReward(streakMap["id"].(string))
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to claim streak reward: %v", c.Account.Username, err))
				}

				for _, status := range claimStreakReward {
					statusMap := status.(map[string]interface{})
					if statusMap["status"].(string) == "received" {
						tools.Logger("success", fmt.Sprintf("| %s | Claim Streak Reward Successfully | Streak Reward: %s", c.Account.Username, statusMap["status"].(string)))
					}
				}
			}
		}
	}

	if isAutoPlaySpinEgg {
		spinEgg, err := c.getSpinEggTicket()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to get spin egg ticket: %v", c.Account.Username, err))
		}
		if spinEgg != nil {
			if len(spinEgg) > 0 {
				for _, ticket := range spinEgg {
					ticketMap := ticket.(map[string]interface{})
					playSpinEgg, err := c.claimSpinEgg(ticketMap["id"].(string))
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Failed to claim spin egg: %v", c.Account.Username, err))
					}

					if playSpinEgg != nil {
						tools.Logger("success", fmt.Sprintf("| %s | Play Spin Egg Successfully | Reward Status: %s | Type: %s", c.Account.Username, playSpinEgg["status"].(string), playSpinEgg["type"].(string)))
					}
				}
			}
		}
	}

	mainTask, err := c.getTasks()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get tasks: %v", c.Account.Username, err))
	}

	if mainTask != nil {
		for _, task := range mainTask {
			if task != nil {
				taskMap := task.(map[string]interface{})

				if taskUser, exits := taskMap["task_user"].(map[string]interface{}); exits && taskUser != nil {
					if !taskUser["completed"].(bool) {

						taskId, err := c.startTask(taskMap["id"].(string))
						if err != nil {
							tools.Logger("error", fmt.Sprintf("| %s | Failed to start task: %v", c.Account.Username, err))
						}

						if taskId != "" {
							tools.Logger("success", fmt.Sprintf("| %s | Start Task %s Successfully | Sleep 5s Before Claim Task...", c.Account.Username, taskMap["name"].(string)))

							time.Sleep(5 * time.Second)

							claimTask, err := c.claimTask(taskId)
							if err != nil {
								tools.Logger("error", fmt.Sprintf("| %s | Failed to claim task: %v", c.Account.Username, err))
							}

							if claimTask != nil {
								if status, exits := claimTask["data"].(map[string]interface{}); exits {
									if status["completed"].(bool) {
										tools.Logger("success", fmt.Sprintf("| %s | Claim Task %s Successfully | Sleep 5s Before Next Task...", c.Account.Username, taskMap["name"].(string)))
									}
								} else {
									tools.Logger("error", fmt.Sprintf("| %s | Claim Task %s Failed | Status: %s | You Can Try Manual | Sleep 5s Before Next Task...", c.Account.Username, taskMap["name"].(string), claimTask["error"].(string)))
								}
							}
						}
					}
				} else {
					c.startTask(taskMap["id"].(string))
				}

				time.Sleep(5 * time.Second)
			}
		}
	}

	if isAutoUpgradeSpeed || isAutoUpgradeStorage || isAutoUpgradeHolyWater {
		settings, err := c.getSettings()
		if err != nil {
			tools.Logger("error", fmt.Sprintf("| %s | Failed to get settings: %v", c.Account.Username, err))
		}

		if settings != nil {
			speedSeedCosts := settings["mining-speed-costs"].([]interface{})
			storageSeedCosts := settings["mining-speed-costs"].([]interface{})

			if isAutoUpgradeSpeed {
				var currentBalance float64

				balance, err = c.getBalance()
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to get balance: %v", c.Account.Username, err))
				}

				if balance != 0 {
					currentBalance = balance / 1e9
				}

				if speedSeedLevel != maxSpeedLevel {
					if speedSeedLevel < 1 {
						speedSeedLevel = speedSeedLevel + 1
					}

					if currentBalance > (float64(speedSeedCosts[(speedSeedLevel-1)].(float64)) / 1e9) {
						upgradeSpeed, err := c.upgradeSpeedSeed()
						if err != nil {
							tools.Logger("error", fmt.Sprintf("| %s | Failed to upgrade speed seed: %v", c.Account.Username, err))
						}

						if upgradeSpeed == "{}" {
							tools.Logger("success", fmt.Sprintf("| %s | Upgrade Speed Seed Successfully | Current Speed Level: %v", c.Account.Username, (speedSeedLevel+1)))
						} else {
							tools.Logger("error", fmt.Sprintf("| %s | Upgrade Speed Seed Failed | Current Speed Level: %v", c.Account.Username, (speedSeedLevel)))
						}
					} else {
						tools.Logger("error", fmt.Sprintf("| %s | Upgrade Speed Seed Failed | Not Enough Balance...", c.Account.Username))
					}
				}
			}

			if isAutoUpgradeStorage {
				var currentBalance float64
				balance, err = c.getBalance()
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to get balance: %v", c.Account.Username, err))
				}

				if balance != 0 {
					currentBalance = balance / 1e9
				}

				if storageSeedLevel != maxStorageLevel {
					if storageSeedLevel < 1 {
						storageSeedLevel = storageSeedLevel + 1
					}

					if currentBalance > (float64(storageSeedCosts[(storageSeedLevel-1)].(float64)) / 1e9) {
						upgradeStorage, err := c.upgradeStorageSeed()
						if err != nil {
							tools.Logger("error", fmt.Sprintf("| %s | Failed to upgrade storage seed: %v", c.Account.Username, err))
						}

						if upgradeStorage == "{}" {
							tools.Logger("success", fmt.Sprintf("| %s | Upgrade Storage Seed Successfully | Current Speed Level: %v", c.Account.Username, (storageSeedLevel+1)))
						} else {
							tools.Logger("error", fmt.Sprintf("| %s | Upgrade Storage Seed Failed | Current Speed Level: %v", c.Account.Username, (storageSeedLevel)))
						}
					} else {
						tools.Logger("error", fmt.Sprintf("| %s | Upgrade Storage Seed Failed | Not Enough Balance...", c.Account.Username))
					}
				}
			}

			if isAutoUpgradeHolyWater {
				holyWaterTask, err := c.getTaskHolyWater()
				if err != nil {
					tools.Logger("error", fmt.Sprintf("| %s | Failed to get holy water task: %v", c.Account.Username, err))
				}

				if holyWaterTask != nil {
					for _, task := range holyWaterTask {
						if taskMap, exits := task.(map[string]interface{}); exits && taskMap != nil {
							var isCompletingReferTask bool
							friendsInfo, err := c.getFriendsInfo()
							if err != nil {
								tools.Logger("error", fmt.Sprintf("| %s | Failed to get friends info: %v", c.Account.Username, err))
							}

							if friendsInfo != nil {
								if len(friendsInfo["referees"].([]interface{})) > 0 {
									isCompletingReferTask = true
								}
							}

							if !isCompletingReferTask && taskMap["type"].(string) == "refer" {
								continue
							}

							taskId, err := c.startTaskHolyWater(taskMap["id"].(string))
							if err != nil {
								tools.Logger("error", fmt.Sprintf("| %s | Failed to start holy water task: %v", c.Account.Username, err))
							}

							if taskId != "" {
								tools.Logger("success", fmt.Sprintf("| %s | Start Holy Water Task %s Successfully | Sleep 5s Before Claim Task...", c.Account.Username, taskMap["name"].(string)))

								time.Sleep(5 * time.Second)

								claimTask, err := c.claimTask(taskId)
								if err != nil {
									tools.Logger("error", fmt.Sprintf("| %s | Failed to claim holy water task: %v", c.Account.Username, err))
								}

								if claimTask != nil {
									if id, exists := claimTask["id"]; exists {
										if idStr, ok := id.(string); ok && idStr == taskId {
											tools.Logger("success", fmt.Sprintf("| %s | Claim Task %s Successfully | Sleep 5s Before Next Task...", c.Account.Username, taskMap["name"].(string)))
										} else {
											tools.Logger("error", fmt.Sprintf("| %s | Task ID mismatch or not a string", c.Account.Username))
										}
									} else {
										tools.Logger("error", fmt.Sprintf("| %s | 'id' key does not exist in claimTask", c.Account.Username))
									}
								}
							}

							time.Sleep(15 * time.Second)
						}
					}
				}

				if holyWaterLevel != maxHolyWaterLevel {
					upgradeHolyWater, err := c.upgradeHolyWater()
					if err != nil {
						tools.Logger("error", fmt.Sprintf("| %s | Failed to upgrade holy water: %v", c.Account.Username, err))
					}

					if upgradeHolyWater == "{}" {
						tools.Logger("success", fmt.Sprintf("| %s | Upgrade Holy Water Successfully | Current Holy Water Level: %v", c.Account.Username, (holyWaterLevel+1)))
					} else {
						tools.Logger("error", fmt.Sprintf("| %s | Upgrade Holy Water Failed | Current Holy Water Level: %v", c.Account.Username, (holyWaterLevel)))
					}
				} else {
					tools.Logger("error", fmt.Sprintf("| %s | Upgrade Holy Water Failed | Not Enough Balance...", c.Account.Username))
				}
			}
		}
	}

	balance, err = c.getBalance()
	if err != nil {
		tools.Logger("error", fmt.Sprintf("| %s | Failed to get balance: %v", c.Account.Username, err))
	}

	if balance != 0 {
		return balance
	} else {
		return 0
	}
}
