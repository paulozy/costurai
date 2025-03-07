package configs

import "github.com/spf13/viper"

type config struct {
	WebPort           string `mapstructure:"WEB_PORT"`
	WebHost           string `mapstructure:"WEB_HOST"`
	JWTSecret         string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn      int64  `mapstructure:"JWT_EXPIRES_IN"`
	FirebaseProjectId string `mapstructure:"FIREBASE_PROJECT_ID"`
	TwilioAccountSID  string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioSMSSID      string `mapstructure:"TWILIO_SMS_SID"`
	TwilioAuthToken   string `mapstructure:"TWILIO_AUTH_TOKEN"`
	Env               string `mapstructure:"ENV"`
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
