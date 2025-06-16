package configs

import (
	"fmt"
	"os"

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
	FirebaseProjectId         string `mapstructure:"FIREBASE_PROJECT_ID"` // <--- O CAMPO EM QUESTÃO
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
	cfg := &Config{} // Inicialize cfg para ser um ponteiro para uma nova Config

	// --- ATIVAR MODO DE DEPURACAO DO VIPER (MUITO VERBOSO, REMOVER PARA PRODUCAO) ---
	viper.Debug() // <<<--- ADICIONE ESTA LINHA
	// viper.SetTypeByDefaultValue(true) // Pode ajudar com valores default, mas Debug é mais útil agora

	if env == "local" {
		viper.SetConfigName("app_config")
		viper.SetConfigType("env")
		viper.AddConfigPath("../")
		viper.SetConfigFile(".env") // Note: SetConfigFile overrides SetConfigName and AddConfigPath for the primary config file

		err := viper.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("erro ao ler app_config.env local: %w", err)
		}

		err = viper.Unmarshal(cfg) // Passe o ponteiro inicializado
		if err != nil {
			return nil, fmt.Errorf("erro ao decodificar config local: %w", err)
		}

		return cfg, nil
	}

	// Caminho para "development", "production" (Cloud Run)
	// --- AQUI É O PONTO CRÍTICO ---
	viper.AutomaticEnv() // Diz ao Viper para ler variáveis de ambiente do SO

	// --- LOG DE DEPURACAO CRUCIAL 1: Verificando o que os.LookupEnv vê ---
	firebaseProjectIDFromOS, foundOS := os.LookupEnv("FIREBASE_PROJECT_ID")
	fmt.Printf("--- Debug OS Env Check ---\n")
	fmt.Printf("os.LookupEnv(\"FIREBASE_PROJECT_ID\"): '%s', Found: %t\n", firebaseProjectIDFromOS, foundOS)
	fmt.Printf("--- Fim Debug OS Env Check ---\n")

	err := viper.Unmarshal(cfg) // <<<--- AQUI O VIPER DEVERIA PREENCHER cfg COM AS VARIAVEIS DE AMBIENTE
	if err != nil {
		// Se este erro for acionado, significa que a estrutura (cfg) ou o mapeamento tem um problema,
		// não necessariamente que a variável está vazia, mas sim um erro de tipo ou de mapeamento.
		return nil, fmt.Errorf("falha ao decodificar configurações do ambiente: %w", err)
	}

	// --- LOG DE DEPURACAO CRUCIAL 2: Verificando o que Viper unmarshaled ---
	fmt.Printf("--- Debug Config Load (via Viper) ---\n")
	fmt.Printf("FirebaseProjectId (lido por Viper): '%s'\n", cfg.FirebaseProjectId)
	fmt.Printf("WebPort (lido por Viper): '%s'\n", cfg.WebPort)
	fmt.Printf("Env (lido por Viper): '%s'\n", cfg.Env)
	fmt.Printf("--- Fim Debug Config Load (via Viper) ---\n")

	// --- LOG DE DEPURACAO 3: DUMP COMPLETO DA STRUCT ---
	// Isso vai mostrar o conteúdo de TODOS os campos da struct 'cfg' após o unmarshal.
	fmt.Printf("--- Full Config Struct Dump ---\n")
	fmt.Printf("%+v\n", cfg) // O '+v' mostra o nome do campo e o valor
	fmt.Printf("--- End Full Config Struct Dump ---\n")

	if cfg.FirebaseProjectId == "" {
		// Isso é o que você está vendo. Se o Debug do Viper não mostrar o motivo,
		// então consideraremos alternativas.
		return nil, fmt.Errorf("FirebaseProjectId está vazio APÓS carregar a configuração (Viper não preencheu)")
	}

	return cfg, nil
}
