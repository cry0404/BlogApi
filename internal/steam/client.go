package steam

import (
	config "BlogApi/config"
	"fmt"
	"time"
	"math/rand"
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
	//这个部分先获取 recently， 然后再更新写入删除， player 信息动态更新
	recentlyGames, err := GetRecentlyPlayedGames(cfg) 
	
	if err != nil {
		return fmt.Errorf("failed to get recently games: %w", err)
	}

	num := len(recentlyGames)
	if num <= 10 {
		recentlyGames = appendRandomGames(recentlyGames, games)
	}else {
		recentlyGames = recentlyGames[:10] //保留前 10
	}

	err = saveAsJson(recentlyGames)
	
	if err != nil {
		return fmt.Errorf("failed to save as json")
	}

	return nil
}

func appendRandomGames(recentlygames, games []Game) []Game {
	needNum := 10 - len(recentlygames)
	


	exist := make(map[int]struct{}, len(recentlygames))
	for _, g := range recentlygames {
		exist[g.APPID] = struct{}{}
	}


	candidates := make([]Game, 0, len(games))
	for _, g := range games {
		if _, ok := exist[g.APPID]; !ok {
			candidates = append(candidates, g)
		}
	}

	if len(candidates) == 0 {
		return recentlygames
	}


	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	if needNum > len(candidates) {
		needNum = len(candidates)
	}

	return append(recentlygames, candidates[:needNum]...)
}