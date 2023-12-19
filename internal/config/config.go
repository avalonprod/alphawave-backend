package config

import (
	"os"
	"time"

	// "github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	defaultHTTPPort           = "8000"
	defaultHTTPWriteTimeout   = 100 * time.Second
	defaultHTTPReadTimeout    = 100 * time.Second
	defaultHTTPMaxHeaderBytes = 1

	defaultAccessTokenTTL  = 10 * time.Minute
	defaultRefreshTokenTTL = 24 * time.Hour * 30
)

type Config struct {
	FrontEndUrl string
	HTTP        HTTPConfig
	Limiter     LimiterConfig
	SMTP        SMTPConfig
	MongoDB     MongoConfig
	Auth        AuthConfig
	Email       EmailConfig
	OpenAI      OpenAIConfig
	Mattermost  MattermostConfig
	MinIO       MinioConfig
}

type (
	HTTPConfig struct {
		Host           string
		Port           string        `mapstructure:"port"`
		ReadTimeout    time.Duration `mapstructure:"readTimeout"`
		WriteTimeout   time.Duration `mapstructure:"writeTimeout"`
		MaxHeaderBytes int           `mapstructure:"maxHeaderBytes"`
	}

	LimiterConfig struct {
		RPS   int           `mapstructure:"rps"`
		BURST int           `mapstructure:"burst"`
		TTL   time.Duration `mapstructure:"ttl"`
	}

	MongoConfig struct {
		Url      string
		Username string
		Password string
		DBName   string
	}
	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string
		VerificationCodeTTL    time.Duration `mapstructure:"verificationCodeTTL"`
		VerificationCodeLength int           `mapstructure:"verificationCodeLength"`
	}

	OpenAIConfig struct {
		Url   string
		Token string
	}

	MinioConfig struct {
		Endpoint        string
		AccessKeyID     string
		SecretAccessKey string
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `mapstructure:"accessTokenTTL"`
		RefreshTokenTTL time.Duration `mapstructure:"refreshTokenTTL"`
		SigningKey      string
	}

	MattermostConfig struct {
		ApiUrl string `mapstructure:"apiUrl"`
	}

	EmailSubjects struct {
		Verification   string `mapstructure:"verification"`
		ForgotPassword string `mapstructure:"forgotPassword"`
		UserInvite     string `mapstructure:"userInvite"`
	}

	EmailTemplates struct {
		Verification   string
		ForgotPassword string
		UserInvite     string
	}

	EmailConfig struct {
		Templates EmailTemplates
		Subjects  EmailSubjects
	}

	SMTPConfig struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		From     string `mapstructure:"from"`
		Password string
	}
)

func Init(configDir string) (*Config, error) {
	// err := godotenv.Load()

	// if err != nil {
	// 	return nil, err
	// }
	var cfg Config
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}
	SetDefault()
	setFromEnv(&cfg)
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func unmarshal(cfg *Config) error {

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("limiter", &cfg.Limiter); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("smtp", &cfg.SMTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth.JWT); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth.verificationCodeTTL", &cfg.Auth.VerificationCodeTTL); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth.verificationCodeLength", &cfg.Auth.VerificationCodeLength); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("email.templates", &cfg.Email.Templates); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("email.subjects", &cfg.Email.Subjects); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("mattermost.apiUrl", &cfg.Mattermost.ApiUrl); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(configsDir string) error {
	viper.AddConfigPath(configsDir)
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func setFromEnv(cfg *Config) {
	cfg.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.MongoDB.Url = os.Getenv("MONGODB_URL")
	cfg.MongoDB.Username = os.Getenv("MONGODB_USERNAME")
	cfg.MongoDB.Password = os.Getenv("MONGODB_PASSWORD")
	cfg.MongoDB.DBName = os.Getenv("MONGODB_NAME")
	cfg.FrontEndUrl = os.Getenv("FRONT_END_URL")
	cfg.Auth.PasswordSalt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.JWT.SigningKey = os.Getenv("JWT_SIGNING_KEY")

	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.OpenAI.Url = os.Getenv("OPEN_AI_API_URL")
	cfg.OpenAI.Token = os.Getenv("OPEN_AI_TOKEN")

	cfg.MinIO.Endpoint = os.Getenv("MINIO_ENDPOINT")
	cfg.MinIO.AccessKeyID = os.Getenv("MINIO_ACCESS_KEY_ID")
	cfg.MinIO.SecretAccessKey = os.Getenv("MINIO_SECRECT_ACCESS_KEY")
}

func SetDefault() {
	viper.SetDefault("http.port", defaultHTTPPort)
	viper.SetDefault("http.maxHeaderBytes", defaultHTTPMaxHeaderBytes)
	viper.SetDefault("http.writeTimeout", defaultHTTPWriteTimeout)
	viper.SetDefault("http.readTimeout", defaultHTTPReadTimeout)
	viper.SetDefault("auth.accessTokenTTL", defaultAccessTokenTTL)
	viper.SetDefault("auth.refreshTokenTTL", defaultRefreshTokenTTL)
}
