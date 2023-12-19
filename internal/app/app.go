package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/config"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/service"
	openai "github.com/Coke15/AlphaWave-BackEnd/internal/infrastructure/ai/openAI"
	"github.com/Coke15/AlphaWave-BackEnd/internal/infrastructure/mattermost"
	httpRoutes "github.com/Coke15/AlphaWave-BackEnd/internal/interface/api/http"
	"github.com/Coke15/AlphaWave-BackEnd/internal/repository"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/auth/manager"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/db/mongodb"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/email/smtp"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/hash"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/paymants"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/storage"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/tokengenerator"
)

const configDir = "configs"

func Run() {

	cfg, err := config.Init(configDir)

	if err != nil {
		panic("error parse config")
	}

	// -----
	hasher := hash.NewHasher(cfg.Auth.PasswordSalt)

	mongoClient, err := mongodb.NewConnection(cfg.MongoDB.Url, cfg.MongoDB.Username, cfg.MongoDB.Password)
	if err != nil {
		logger.Errorf("failed to create new mongo client. err: %v", err)
		panic("failed to create new mongo client")
	}

	JWTManager, err := manager.NewJWTManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	emailSender, err := smtp.NewSMTPSender(cfg.SMTP.From, cfg.SMTP.Password, cfg.SMTP.Host, cfg.SMTP.Port)
	if err != nil {
		logger.Error(err)
		return
	}
	codeGenerator := codegenerator.NewCodeGenerator()

	openAI := openai.NewOpenAiAPI(cfg.OpenAI.Token, cfg.OpenAI.Url)
	fmt.Printf("endpoint: %s, acesskeyid: %s, secretaccesskey: %s", cfg.MinIO.Endpoint, cfg.MinIO.AccessKeyID, cfg.MinIO.SecretAccessKey)
	storageProvider, err := storage.NewClient(cfg.MinIO.Endpoint, cfg.MinIO.AccessKeyID, cfg.MinIO.SecretAccessKey)

	if err != nil {
		logger.Error(err)
		return
	}

	paymentProvider := paymants.NewPaymentProvider("sk_test_51NnlKlH75mUJKHqVvdKp7fZOTPu6QqoXr4Sc5YxXKdbY6H4QY6O9dwwEc9VAMiT3CrcMZoNTPWk2whrArX5Phz4z00k5N8TkN9")

	// paymentProvider.Paymant(paymants.PaymantPayload{
	// 	Amount:   2000,
	// 	Currency: "usd",
	// })

	// id, err := paymentProvider.CreateCustomer("Roman", "abramenkoroman22@gmail.com", "New Client")
	// secret, err := paymentProvider.NewCard("cus_OkkHewbxDaNh2T")

	mattermostAdapter := mattermost.NewMattermostAdapter(cfg.Mattermost.ApiUrl)

	tokenGenerator := tokengenerator.NewTokenGenerator()
	// -----

	mongodb := mongoClient.Database(cfg.MongoDB.DBName)
	repository := repository.NewRepository(mongodb)
	service := service.NewService(&service.Deps{
		UserRepository:         repository.User,
		TasksRepository:        repository.Tasks,
		ProjectsRepository:     repository.Projects,
		TeamsRepository:        repository.Teams,
		RolesRepository:        repository.Roles,
		MemberRepository:       repository.Members,
		PackagesRepository:     repository.Packages,
		FilesRepository:        repository.Files,
		FolderRepository:       repository.Folder,
		SubscriptionRepository: repository.Subscription,
		StorageProvider:        storageProvider,
		StorageEndpointURL:     cfg.MinIO.Endpoint,
		Hasher:                 hasher,
		JWTManager:             JWTManager,
		AccessTokenTTL:         cfg.Auth.JWT.AccessTokenTTL,
		RefreshTokenTTL:        cfg.Auth.JWT.RefreshTokenTTL,
		VerificationCodeTTL:    cfg.Auth.VerificationCodeTTL,
		Sender:                 emailSender,
		MattermostAdapter:      mattermostAdapter,
		PaymentProvider:        paymentProvider,
		EmailConfig:            cfg.Email,
		CodeGenerator:          codeGenerator,
		TokenGenerator:         tokenGenerator,
		VerificationCodeLength: cfg.Auth.VerificationCodeLength,
		ApiUrl:                 cfg.HTTP.Host,
		OpenAI:                 openAI,
	})
	handler := httpRoutes.NewHandler(service, JWTManager, cfg.Auth.JWT.RefreshTokenTTL, cfg.FrontEndUrl)

	srv := NewServer(cfg, handler.InitRoutes(cfg))
	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()
	if err := srv.Shotdown(ctx); err != nil {
		logger.Errorf("failed to stop server: %x", err)
	}
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Errorf("error disconnect to mongoClient. err: %v", err)
	}
}

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.HTTP.Port,
			Handler:        handler,
			MaxHeaderBytes: cfg.HTTP.MaxHeaderBytes << 20,
			ReadTimeout:    cfg.HTTP.ReadTimeout,
			WriteTimeout:   cfg.HTTP.WriteTimeout,
		},
	}
}

func (s *Server) Run() error {
	port := strings.Replace(s.httpServer.Addr, ":", "", 1)

	logger.Infof("Server has ben started on port: %s", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shotdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
