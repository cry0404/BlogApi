package steam

type SteamAPIResponse struct {
	Response OwnedGamesResponse `json:"response"`
}

type OwnedGamesResponse struct {
	GameCount int    `json:"game_count"`
	Games     []Game `json:"games"`
}

//完整的图标构建 
//https://media.steampowered.com/steamcommunity/public/images/apps/{appid}/{img_icon_url}.jpg
type Game struct {
	APPID 			int		`json:"appid"`
	Name            string  `json:"name"`
	PlaytimeForever int     `json:"playtime_forever"`
	RtimeLastPlayed int64   `json:"rtime_last_played"`
	ImgIconURL      string  `json:"img_icon_url"`
	Comment  		string    //为以后的评测保留
}

type GameInfo struct {

	Name            string  `json:"name"`
	ImagePath		string  `json:"image_path"` //这里直接根据 appid 转换得到
	PlaytimeForever int     `json:"playtime_forever"`
	RtimeLastPlayed int64   `json:"rtime_last_played"`
	Comment  		string  `json:"comment"`  //为以后的评测保留
}

type PlayerSummariesAPIResponse struct {
	Response PlayerSummariesResponse `json:"response"`
}

type PlayerSummariesResponse struct {
	Players []Player `json:"players"`
}

type Player struct {
	SteamID      string `json:"steamid"`
	ProfileState int    `json:"profilestate"`
	PersonaName  string `json:"personaname"`

	ProfileURL string `json:"profileurl"`
	Avatar     string `json:"avatar"`

	LastLogoff int64 `json:"lastlogoff"`

	TimeCreated int64 `json:"timecreated"`
}
