package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	FeiShu   FeiShuConfig
	Bilibili BiliBiliConfig
	Steam    SteamConfig
}

type FeiShuConfig struct {
	FeiShuAppID     string `mapstructure:"app_id"`
	FeiShuAppSecret string `mapstructure:"app_secret"`
	BookTableID     string `mapstructure:"book_table_id"`
	MovieTableID    string `mapstructure:"movie_table_id"`
	AnimeTableID    string `mapstructure:"anime_table_id"`
	FeiShuAppToken  string `mapstructure:"app_token"`
}

type BiliBiliConfig struct {
}

type SteamConfig struct {
	SteamID  string `mapstructure:"steam_id"`
	SteamKey string `mapstructure:"steam_key"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetEnvPrefix("BLOGAPI")
	viper.AutomaticEnv()

	viper.BindEnv("feishu.app_id", "BLOGAPI_FEISHU_APP_ID")
	viper.BindEnv("feishu.app_secret", "BLOGAPI_FEISHU_APP_SECRET")
	viper.BindEnv("feishu.book_table_id", "BLOGAPI_FEISHU_BOOK_TABLE_ID")
	viper.BindEnv("feishu.movie_table_id", "BLOGAPI_FEISHU_MOVIE_TABLE_ID")
	viper.BindEnv("feishu.anime_table_id", "BLOGAPI_FEISHU_ANIME_TABLE_ID")
	viper.BindEnv("feishu.app_token", "BLOGAPI_FEISHU_APP_TOKEN")
	viper.BindEnv("steam.steam_id", "BLOGAPI_STEAM_ID")
	viper.BindEnv("steam.steam_key", "BLOGAPI_STEAM_key")

	viper.SetDefault("feishu.app_id", "")
	viper.SetDefault("feishu.app_secret", "")
	viper.SetDefault("feishu.book_table_id", "")
	viper.SetDefault("feishu.movie_table_id", "")
	viper.SetDefault("feishu.anime_table_id", "")
	viper.SetDefault("feishu.app_token", "")
	viper.BindEnv("steam.steam_id", "")
	viper.BindEnv("steam.steam_key", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("配置文件未找到，使用默认值和环境变量")
		} else {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &config, nil
}
