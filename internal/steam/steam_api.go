package steam

import (
	config "BlogApi/config"
	"encoding/json"
	"fmt"
	"net/http"
)

// 获取拥有的游戏列表
func GetOwnedGames(cfg *config.Config) ([]Game, error) {

	resp, err := http.Get(ownedGamesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owned games: %w", err)
	}
	defer resp.Body.Close()

	var apiResponse SteamAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var gameInfos []Game
	for _, game := range apiResponse.Response.Games {
		gameInfos = append(gameInfos, Game{
			Name:            game.Name,
			PlaytimeForever: game.PlaytimeForever,
			RtimeLastPlayed: game.RtimeLastPlayed,
			ImgIconURL:      game.ImgIconURL,
		})
	}

	return gameInfos, nil
}

// 获取最近两周玩的游戏
func GetRecentlyPlayedGames(cfg *config.Config) ([]Game, error) {
	resp, err := http.Get(recentlyPlayedGamesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recently played games: %w", err)
	}
	defer resp.Body.Close()

	var apiResponse SteamAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return apiResponse.Response.Games, nil
}

// 获取用户基本信息
func GetPlayerSummaries(cfg *config.Config) (*Player, error) {
	resp, err := http.Get(playerSummariesURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player summaries: %w", err)
	}
	defer resp.Body.Close()

	var apiResponse PlayerSummariesAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(apiResponse.Response.Players) == 0 {
		return nil, fmt.Errorf("no player found")
	}

	return &apiResponse.Response.Players[0], nil
}
