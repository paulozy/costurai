package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver                  string `mapstructure:"DB_DRIVER"`
	DBHost                    string `mapstructure:"DB_HOST"`
	DBPort                    string `mapstructure:"DB_PORT"`
	DBUser                    string `mapstructure:"DB_USER"`
	DBPassword                string `mapstructure:"DB_PASSWORD"`
	DBName                    string `mapstructure:"DB_NAME"`
	WebPort                   string `mapstructure:"WEB_PORT"`
	WebHost                   string `mapstructure:"WEB_HOST"`
	JWTSecret                 string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn              int64  `mapstructure:"JWT_EXPIRES_IN"`
	FirebaseProjectId         string `mapstructure:"FIREBASE_PROJECT_ID"`
	TwilioSID                 string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken           string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioSMSServiceSID       string `mapstructure:"TWILIO_SMS_SERVICE_SID"`
	SMSTimeout                int64  `mapstructure:"SMS_TIMEOUT"`
	TwilioChannel             string `mapstructure:"TWILIO_CHANNEL"`
	DBType                    string `mapstructure:"DB_TYPE"`
	PaymentSuccessRedirectURL string `mapstructure:"PAYMENT_SUCCESS_REDIRECT_URL"`
	PaymentCancelRedirectURL  string `mapstructure:"PAYMENT_CANCEL_REDIRECT_URL"`
	StripeSecretKey           string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret       string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	Env                       string `mapstructure:"ENV"`
}

func LoadConfig(env string) (*Config, error) {
	var cfg *Config

	if env == "local" {
		viper.SetConfigName("app_config")
		viper.SetConfigType("env")
		viper.AddConfigPath("../")
		viper.SetConfigFile(".env")

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

	viper.AutomaticEnv()
	err := viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg, nil
}
