package core

import "fmt"

// Get Profile
func (c *Client) getProfile() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/profile"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Make request error: %v", err)
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	} else {
		return nil, fmt.Errorf("Data field not exits!")
	}
}

// Get Balance
func (c *Client) getBalance() (float64, error) {
	url := "https://elb.seeddao.org/api/v1/profile/balance"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	if balance, exits := res["data"].(float64); exits && balance > 0 {
		return balance, nil
	} else {
		return 0, fmt.Errorf("Data Field Not Exits!")
	}
}

// Check Guild
func (c *Client) checkGuild() (string, error) {
	url := "https://elb.seeddao.org/api/v1/guild/member/detail"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Make request error: %v", err)
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data["guild_id"].(string), nil
	} else {
		return "", fmt.Errorf("Data field not exits!")
	}
}

// Join Guild
func (c *Client) joinGuild() {
	url := "https://elb.seeddao.org/api/v1/guild/join"

	payload := map[string]string{
		"guild_id": "9e02254f-d921-43d3-839f-903706dedeb5",
	}

	c.makeRequest("POST", url, payload)
}

// Leave Guild
func (c *Client) leaveGuild() {
	url := "https://elb.seeddao.org/api/v1/guild/leave"

	payload := map[string]string{
		"guild_id": "9e02254f-d921-43d3-839f-903706dedeb5",
	}

	c.makeRequest("POST", url, payload)
}

// Get Guild Info
func (c *Client) getGuildInfo() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/guild/detail?guild_id=9e02254f-d921-43d3-839f-903706dedeb5&sort_by=total_hunted"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Bird Inventory
func (c *Client) getBirdInventory() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/bird/me?page=1"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Worm Inventory
func (c *Client) getWormInventory() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/worms/me?page=1"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Egg Inventory
func (c *Client) getEggInventory() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/egg/me?page=1"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Check Login Bonus
func (c *Client) checkLoginBonus() ([]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/login-bonuses"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].([]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Claim Login Bonus
func (c *Client) claimLoginBonus() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/login-bonuses"

	res, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Claim Farming Seed
func (c *Client) claimFarmingSeed() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/seed/claim"

	res, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Claim First Egg
func (c *Client) claimFirstEgg() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/give-first-egg"

	res, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Hatch Egg
func (c *Client) hatchEgg(eggId string) (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/egg-hatch/complete"

	payload := map[string]interface{}{
		"egg_id": eggId,
	}

	res, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Bird Status
func (c *Client) getBirdStatus() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/bird/is-leader"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Worm Status
func (c *Client) getWormStatus() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/worms"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Catch Worm
func (c *Client) catchWorm() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/worms/catch"

	res, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Feed Bird
func (c *Client) feedBird(birdId, wormId string) (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/bird-feed"

	wormIds := []string{wormId}
	payload := map[string]interface{}{
		"bird_id":  birdId,
		"worm_ids": wormIds,
	}

	res, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Bird Happiness
func (c *Client) birdHappiness(birdId string, happinessRate int) (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/bird-happiness"

	payload := map[string]interface{}{
		"bird_id":        birdId,
		"happiness_rate": happinessRate * 100,
	}

	res, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Claim Bird Hunt
func (c *Client) claimBirdHunt(birdId string) (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/bird-hunt/complete"

	payload := map[string]string{
		"bird_id": birdId,
	}

	res, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Start Bird Hunt
func (c *Client) startBirdHunt(birdId string, taskLevel int) (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/bird-hunt/start"

	payload := map[string]interface{}{
		"bird_id":    birdId,
		"task_level": taskLevel,
	}

	res, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Tasks
func (c *Client) getTasks() ([]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/tasks/progresses"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].([]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Start Task
func (c *Client) startTask(taskId string) (string, error) {
	url := fmt.Sprintf("https://elb.seeddao.org/api/v1/tasks/%s", taskId)

	res, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	if data, exits := res["data"].(string); exits && data != "" {
		return data, nil
	}

	return "", fmt.Errorf("Data field not exits!")
}

func (c *Client) claimTask(taskId string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://elb.seeddao.org/api/v1/tasks/notification/%s", taskId)

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Error in makeRequest: %v", err)
	}

	if res == nil {
		return nil, fmt.Errorf("Response is nil")
	}

	if data, exists := res["data"]; exists {
		// Pastikan data bukan nil
		if dataMap, ok := data.(map[string]interface{}); ok {
			return dataMap, nil
		} else {
			return nil, fmt.Errorf("'data' exists but is not a map[string]interface{}")
		}
	}

	return nil, fmt.Errorf("Data field does not exist!")
}

// Get Holy Water Task
func (c *Client) getTaskHolyWater() ([]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/upgrades/tasks/progresses"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].([]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Friends Info
func (c *Client) getFriendsInfo() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/profile/recent-referees"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Start Holy Water Task
func (c *Client) startTaskHolyWater(taskId string) (string, error) {
	url := fmt.Sprintf("https://elb.seeddao.org/api/v1/upgrades/tasks/%s", taskId)

	res, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	if data, exits := res["data"].(string); exits && data != "" {
		return data, nil
	}

	return "", fmt.Errorf("Data field not exits!")
}

// Get Streak Reward
func (c *Client) getStreakReward() ([]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/streak-reward"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].([]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Claim Streak Reward
func (c *Client) claimStreakReward(streakId string) ([]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/streak-reward"

	streakIds := []string{streakId}
	payload := map[string]interface{}{
		"streak_reward_ids": streakIds,
	}

	res, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].([]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Get Spin Egg Ticket
func (c *Client) getSpinEggTicket() ([]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/spin-ticket"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].([]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

// Claim Spin Egg
func (c *Client) claimSpinEgg(ticketId string) (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/spin-reward"

	payload := map[string]string{
		"ticket_id": ticketId,
	}

	res, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

func (c *Client) getSettings() (map[string]interface{}, error) {
	url := "https://elb.seeddao.org/api/v1/settings"

	res, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if data, exits := res["data"].(map[string]interface{}); exits && data != nil {
		return data, nil
	}

	return nil, fmt.Errorf("Data field not exits!")
}

func (c *Client) upgradeSpeedSeed() (string, error) {
	url := "https://elb.seeddao.org/api/v1/seed/mining-speed/upgrade"

	_, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	return "{}", nil
}

func (c *Client) upgradeStorageSeed() (string, error) {
	url := "https://elb.seeddao.org/api/v1/seed/storage-size/upgrade"

	_, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	return "{}", nil
}

func (c *Client) upgradeHolyWater() (string, error) {
	url := "https://elb.seeddao.org/api/v1/upgrades/holy-water"

	_, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return "", err
	}

	return "{}", nil
}
