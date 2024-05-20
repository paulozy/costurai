package configs

import "github.com/spf13/viper"

type config struct {
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	WebPort      string `mapstructure:"WEB_PORT"`
	WebHost      string `mapstructure:"WEB_HOST"`
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn int64  `mapstructure:"JWT_EXPIRES_IN"`
}

func LoadConfig(path string) (*config, error) {
	var cfg *config

	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}
