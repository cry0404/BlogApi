package steam

import (
	config "BlogApi/config"
	"fmt"
	
)

var (
	steamID                string
	steamKey               string
	ownedGamesURL          string
	recentlyPlayedGamesURL string
	playerSummariesURL     string
)

func initSteamConfig(cfg *config.Config) {
	steamID = cfg.Steam.SteamID
	steamKey = cfg.Steam.SteamKey
	ownedGamesURL = fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json&include_appinfo=true&include_played_free_games=true", steamKey, steamID)
	recentlyPlayedGamesURL = fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v0001/?key=%s&steamid=%s&format=json&count=10", steamKey, steamID)
	playerSummariesURL = fmt.Sprintf("http://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", steamKey, steamID)
}

func UpdateGames(cfg *config.Config) error {
	initSteamConfig(cfg)
	games, err := GetOwnedGames(cfg)
	//测试
	/*
	for _, game := range games {
		fmt.Println("id 是", game.APPID)
	}*/
	if err != nil {
		return fmt.Errorf("failed to get owned games: %w", err)
	}

	err = DownloadImages(games)

	if err != nil {
		return err
	}

	//err = saveAsJson()
	// 这里可以添加将游戏数据保存到数据库或文件的逻辑
	//fmt.Printf("获取到 %d 个游戏\n", len(games))

	return nil
}


