package utils

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Config struct {
	Url      string        `yaml:"URL"`
	Method   string        `yaml:"METHOD"`
	Location string        `yaml:"QUERY"`
	Time     time.Duration `yaml:"TIME"`

	AccessKey string
}

func LoadConfig(path string, fileName string) (config Config, err error) {
	// Load config file formatted in .yaml
	split := strings.Split(fileName, ".")
	viper.AddConfigPath(path)
	viper.SetConfigName(split[0])
	viper.SetConfigType(split[1])
	viper.AutomaticEnv()

	// Load environment variables
	config.AccessKey = os.Getenv("ACCESS_KEY")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}