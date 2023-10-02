package config

import (
	"time"

	"github.com/spf13/viper"
)

// for testing purpose please change config path into absolute path
func Reader() (config EnvConfig, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// will replace with existing env
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

type Schema struct {
	HTTPServer HTTPServer
	Provider   Provider
	Storage    Storage
	App        App
}

type App struct {
	JWTSecret string
}

type PSQL struct {
	DBName   string
	User     string
	Password string
	Host     string
	Port     int
}

type Storage struct {
	PSQL PSQL
}

type Provider struct {
	MerchantCodeLfi       string
	SecretLfi             string
	TokenAuthLfiUrl       string
	PaymentRequestLfiUrl  string
	PaymentRedirectLfiUrl string
	ApiVersion            string
	CallbackUrl           string
	ReturnUrl             string
	MerchantIdPrismalink  string
	ClientKeyPrismalink   string
	RequestPrismalinkURL  string
}

type HTTPServer struct {
	ListenAddress   string
	Port            string
	GracefulTimeout time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

func BindConfig(env EnvConfig) Schema {
	return Schema{
		HTTPServer: HTTPServer{
			ListenAddress:   "0.0.0.0",
			Port:            env.ServerPort,
			GracefulTimeout: time.Second * 30,
			ReadTimeout:     time.Minute * 3,
			WriteTimeout:    time.Minute * 5,
			IdleTimeout:     time.Minute * 15,
		},
		Provider: Provider{
			MerchantCodeLfi:       env.ProviderMerchantCodeLfi,
			SecretLfi:             env.ProviderMerchantSecretLfi,
			TokenAuthLfiUrl:       env.ProviderTokenAuthLfiUrl,
			PaymentRequestLfiUrl:  env.ProviderPaymentRequestLfiUrl,
			PaymentRedirectLfiUrl: env.ProviderPaymentRedirectLFIUrl,
			ApiVersion:            env.ProviderApiVersionLfi,
			CallbackUrl:           env.ProviderCallbackUrl,
			ReturnUrl:             env.ProviderReturnUrl,
			MerchantIdPrismalink:  env.ProviderMerchantIdPrismalink,
			ClientKeyPrismalink:   env.ProviderClientKeyPrismalink,
			RequestPrismalinkURL:  env.ProviderRequestPrismalinkURL,
		},
		Storage: Storage{
			PSQL: PSQL{
				DBName:   env.StorageDatabaseName,
				User:     env.StorageDatabaseUsername,
				Password: env.StorageDatabasePassword,
				Host:     env.StorageDatabaseHost,
				Port:     5432,
			},
		},
		App: App{
			JWTSecret: env.AppJWTSecret,
		},
	}
}

type EnvConfig struct {
	ServerPort                    string `mapstructure:"CONFIG_SERVER_PORT"`
	StorageDatabaseName           string `mapstructure:"CONFIG_STORAGE_DB_NAME"`
	StorageDatabaseUsername       string `mapstructure:"CONFIG_STORAGE_USERNAME"`
	StorageDatabasePassword       string `mapstructure:"CONFIG_STORAGE_PASSWORD"`
	StorageDatabaseHost           string `mapstructure:"CONFIG_STORAGE_HOST"`
	AppJWTSecret                  string `mapstructure:"CONFIG_APP_JWT_SECRET"`
	ProviderMerchantCodeLfi       string `mapstructure:"CONFIG_PROVIDER_MERCHANT_CODE_LFI"`
	ProviderMerchantSecretLfi     string `mapstructure:"CONFIG_PROVIDER_SECRET_LFI"`
	ProviderTokenAuthLfiUrl       string `mapstructure:"CONFIG_PROVIDER_TOKEN_AUTH_LFI_URL"`
	ProviderPaymentRequestLfiUrl  string `mapstructure:"CONFIG_PROVIDER_PAYMENT_REQUEST_LFI_URL"`
	ProviderApiVersionLfi         string `mapstructure:"CONFIG_PROVIDER_API_VERSION_LFI"`
	ProviderCallbackUrl           string `mapstructure:"CONFIG_CALLBACK_URL"`
	ProviderReturnUrl             string `mapstructure:"CONFIG_RETURN_URL"`
	ProviderPaymentRedirectLFIUrl string `mapstructure:"CONFIG_PROVIDER_PAYMENT_REDIRECT_LFI_URL"`
	ProviderMerchantIdPrismalink  string `mapstructure:"CONFIG_MERCHANT_ID_PRISMALINK"`
	ProviderClientKeyPrismalink   string `mapstructure:"CONFIG_CLIENT_KEY_PRISMALINK"`
	ProviderRequestPrismalinkURL  string `mapstructure:"CONFIG_REQUEST_PRISMALINK_URL"`
}
