package steam

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

//这里的展示逻辑应该是
var configName string = "./config/steam.ndjson"
func saveAsJson(games []Game) error {
	//每周随机覆盖写入
	file, err := os.Create(configName)
	if err != nil {
		return fmt.Errorf("创建覆盖写入时发生错误")
	}
	defer file.Close()
	bw := bufio.NewWriter(file)
	enc := json.NewEncoder(bw) 
	enc.SetIndent("", " ")
	for _ , game := range games {
		gameinfo := convertGameInfo(game)
		
		if err := enc.Encode(gameinfo); err != nil {
			return fmt.Errorf("编码json时发生错误")
		}
	}	
	if err := bw.Flush(); err != nil {
		return fmt.Errorf("刷新失败: %w", err)
	}
	

	return nil
}


func convertGameInfo(game Game) *GameInfo {

	return &GameInfo{
		Name: game.Name,
		ImagePath: filepath.Join("public", "steam", strconv.Itoa(game.APPID) + ".webp"),
		PlaytimeForever: game.PlaytimeForever,
		RtimeLastPlayed: game.RtimeLastPlayed,
		Comment: game.Comment,
	}
}