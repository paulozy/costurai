package configs

import (
	"fmt"
	"os"
	"strconv" // Necessário para converter string para int64

	"github.com/spf13/viper" // Mantenha para o caso 'local' se ainda o usar
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
	JWTExpiresIn              int64  `mapstructure:"JWT_EXPIRES_IN"` // int64
	FirebaseProjectId         string `mapstructure:"FIREBASE_PROJECT_ID"`
	TwilioSID                 string `mapstructure:"TWILIO_ACCOUNT_SID"`
	TwilioAuthToken           string `mapstructure:"TWILIO_AUTH_TOKEN"`
	TwilioSMSServiceSID       string `mapstructure:"TWILIO_SMS_SERVICE_SID"`
	SMSTimeout                int64  `mapstructure:"SMS_TIMEOUT"` // int64
	TwilioChannel             string `mapstructure:"TWILIO_CHANNEL"`
	DBType                    string `mapstructure:"DB_TYPE"`
	PaymentSuccessRedirectURL string `mapstructure:"PAYMENT_SUCCESS_REDIRECT_URL"`
	PaymentCancelRedirectURL  string `mapstructure:"PAYMENT_CANCEL_REDIRECT_URL"`
	StripeSecretKey           string `mapstructure:"STRIPE_SECRET_KEY"`
	StripeWebhookSecret       string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
	Env                       string `mapstructure:"ENV"`
}

var Env string

func LoadConfig(env string) (*Config, error) {
	Env = env
	cfg := &Config{} // Inicialize a struct

	if env == "local" {
		// Para o ambiente local, podemos continuar usando Viper com o arquivo .env
		// Ou você pode migrar para os.Getenv aqui também para consistência se preferir.
		fmt.Println("Carregando configurações do ambiente LOCAL (via arquivo .env)")
		viper.SetConfigName("app_config")
		viper.SetConfigType("env")
		viper.AddConfigPath("../")
		viper.SetConfigFile(".env") // Isso define o caminho absoluto ou relativo para o .env

		err := viper.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("erro ao ler app_config.env local: %w", err)
		}

		err = viper.Unmarshal(cfg)
		if err != nil {
			return nil, fmt.Errorf("erro ao decodificar config local: %w", err)
		}

	} else { // Caminho para ambientes de produção/desenvolvimento (Cloud Run)
		fmt.Println("Carregando configurações do ambiente Cloud Run (diretamente de variáveis de ambiente do SO)")

		// Preencha a struct Config diretamente de os.Getenv
		cfg.DBDriver = os.Getenv("DB_DRIVER")
		cfg.DBHost = os.Getenv("DB_HOST")
		cfg.DBPort = os.Getenv("DB_PORT")
		cfg.DBUser = os.Getenv("DB_USER")
		cfg.DBPassword = os.Getenv("DB_PASSWORD")
		cfg.DBName = os.Getenv("DB_NAME")
		cfg.WebPort = os.Getenv("WEB_PORT")
		cfg.WebHost = os.Getenv("WEB_HOST")
		cfg.JWTSecret = os.Getenv("JWT_SECRET")

		// JWTExpiresIn (int64)
		jwtExpiresInStr := os.Getenv("JWT_EXPIRES_IN")
		if jwtExpiresInStr != "" {
			val, err := strconv.ParseInt(jwtExpiresInStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("erro ao converter JWT_EXPIRES_IN: %w", err)
			}
			cfg.JWTExpiresIn = val
		}

		cfg.FirebaseProjectId = os.Getenv("FIREBASE_PROJECT_ID")

		cfg.TwilioSID = os.Getenv("TWILIO_ACCOUNT_SID")
		cfg.TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
		cfg.TwilioSMSServiceSID = os.Getenv("TWILIO_SMS_SERVICE_SID")

		// SMSTimeout (int64)
		smsTimeoutStr := os.Getenv("SMS_TIMEOUT")
		if smsTimeoutStr != "" {
			val, err := strconv.ParseInt(smsTimeoutStr, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("erro ao converter SMS_TIMEOUT: %w", err)
			}
			cfg.SMSTimeout = val
		}

		cfg.TwilioChannel = os.Getenv("TWILIO_CHANNEL")
		cfg.DBType = os.Getenv("DB_TYPE")
		cfg.PaymentSuccessRedirectURL = os.Getenv("PAYMENT_SUCCESS_REDIRECT_URL")
		cfg.PaymentCancelRedirectURL = os.Getenv("PAYMENT_CANCEL_REDIRECT_URL")
		cfg.StripeSecretKey = os.Getenv("STRIPE_SECRET_KEY")
		cfg.StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
		cfg.Env = os.Getenv("ENV") // O valor "development" ou "production"

		// Manter os logs de depuração para os.LookupEnv para confirmação final
		firebaseProjectIDFromOS, foundOS := os.LookupEnv("FIREBASE_PROJECT_ID")
		fmt.Printf("--- Debug OS Env Check ---\n")
		fmt.Printf("os.LookupEnv(\"FIREBASE_PROJECT_ID\"): '%s', Found: %t\n", firebaseProjectIDFromOS, foundOS)
		fmt.Printf("--- Fim Debug OS Env Check ---\n")

		// Log de depuração final da struct preenchida por os.Getenv
		fmt.Printf("--- Full Config Struct Dump (via os.Getenv) ---\n")
		fmt.Printf("%+v\n", cfg)
		fmt.Printf("--- End Full Config Struct Dump (via os.Getenv) ---\n")
	}

	// Esta verificação ainda é válida, pois garantirá que o valor foi carregado.
	if cfg.FirebaseProjectId == "" {
		return nil, fmt.Errorf("FirebaseProjectId está vazio APÓS carregar a configuração (provavelmente via os.Getenv)")
	}

	return cfg, nil
}
