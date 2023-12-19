package service

import (
	"context"
	"io"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/config"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/auth/manager"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/email"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/hash"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/tokengenerator"
)

type MattermostAdapter interface {
	CreateUser(ctx context.Context, input types.CreateUserMattermostPayloadDTO) error
	SignIn(email string, password string) (string, error)
}

type PaymentProvider interface {
	CreateCustomer(name, email, descr string) (*string, error)
	NewCard(customerID string) (secret *string, err error)
	CreateSubscription(customerID, priceID string) (*string, error)
}

type UserServiceI interface {
	SignUp(ctx context.Context, input types.UserSignUpDTO) error
	SignIn(ctx context.Context, input types.UserSignInDTO) (types.AuthPayload, error)
	GetUserById(ctx context.Context, userID string) (types.UserDTO, error)
	ChangePassword(ctx context.Context, userID, newPassword, oldPassword string) error
	UpdateUserInfo(ctx context.Context, userID string, input types.UpdateUserInfoDTO) error
	UpdateUserSettings(ctx context.Context, userID string, input types.UpdateUserSettingsDTO) error
	UploadUserAvatar(ctx context.Context, userID string, reader io.Reader, extension string, fileName string, size int) (string, error)
	UploadUserBanner(ctx context.Context, userID string, reader io.Reader, extension string, fileName string, size int) (string, error)
	ResetPassword(ctx context.Context, email, token, tokenResult, password string) error
	VerifyForgotPasswordToken(ctx context.Context, email, token, tokenResult string) (types.ForgotPasswordPayloadDTO, error)
	ForgotPassword(ctx context.Context, email string) error
	LogOut(ctx context.Context, userID string) error
	ResendVerificationCode(ctx context.Context, email string) error
	RefreshTokens(ctx context.Context, refreshToken string) (types.AuthPayload, error)
	Verify(ctx context.Context, verificationCode string) (types.Tokens, error)
}

type MemberServiceI interface {
	MemberSignUp(ctx context.Context, token string, input types.MemberSignUpDTO) error
	GetMembersByQuery(ctx context.Context, teamID string, query types.GetMembersByQuery) ([]types.MemberDTO, error)
	GetMemberByTeamIdAndUserId(ctx context.Context, teamID string, userID string) (model.Member, error)
	UserInvite(ctx context.Context, teamID string, email string, role string) error
	UpdateRoles(ctx context.Context, memberId, teamId string, roles []string) error
	AcceptInvite(ctx context.Context, token string) (string, error)
}

type TeamsServiceI interface {
	Create(ctx context.Context, userID string, input types.CreateTeamsDTO) (string, error)
	GetTeamByID(ctx context.Context, userID, teamID string) (model.Team, error)
	GetTeamByOwnerId(ctx context.Context, ownerId string) (types.TeamsDTO, error)
	GetTeamsByUser(ctx context.Context, userID string) ([]model.Team, error)
	UpdateTeamSettings(ctx context.Context, teamID string, input types.UpdateTeamSettingsDTO) error
}

type RolesServiceI interface {
	Create(ctx context.Context, teamID string) error
	GetRolesByTeamId(ctx context.Context, teamID string) ([]types.GetRoleDTO, error)
	UpdatePermissions(ctx context.Context, teamID string, input []types.UpdatePermissionsDTO) error
}

type TasksServiceI interface {
	Create(ctx context.Context, userID string, input types.TasksCreateDTO) (string, error)
	UpdatePosition(ctx context.Context, userId string, input []types.UpdatePositionDTO) error
	GetById(ctx context.Context, userID, taskID string) (types.TaskDTO, error)
	GetAll(ctx context.Context, userID string) (types.TasksDTO, error)
	UpdateById(ctx context.Context, userID string, input types.UpdateTaskDTO) (types.TaskDTO, error)
	FinishedTaskById(ctx context.Context, userID, taskID string) error
	DeleteTaskById(ctx context.Context, userID, taskID string) error
	DeleteAll(ctx context.Context, userID string) error
	ClearAll(ctx context.Context, userID string) error
	UndoTask(ctx context.Context, userID string, taskID string) error
}

type AiChatServiceI interface {
	NewMessage(messages []types.Message) (*types.MessageOutput, error)
}

type PackagesServiceI interface {
	// CreateDefaultPackages() error
	GetAll(ctx context.Context) ([]model.Package, error)
	GetById(ctx context.Context, packageId string) (model.Package, error)
}

type FilesServiceI interface {
	Create(ctx context.Context, reader io.Reader, input types.CreateFileDTO) (types.FileDTO, error)
	RenameFile(ctx context.Context, teamId string, fileId string, name string) (string, error)
	CreateFolder(ctx context.Context, input types.CreateFolderDTO) (types.FolderDTO, error)
	GetFilePresignedURL(ctx context.Context, teamID, fileID string) (string, error)
	GetFolderContent(ctx context.Context, teamId, folderID string) (types.FolderContentDTO, error)
	GetFolderRoot(ctx context.Context, teamId string) (types.RootFolderDTO, error)
	CreateRootFolder(ctx context.Context, teamId string) error
	Delete(ctx context.Context, teamId, fileId string) error
	UploadImage(ctx context.Context, reader io.Reader, input types.UploadImageDTO) (types.ImageOutputDTO, error)
	DeleteImage(ctx context.Context, filePath string) error
}

type PaymentServiceI interface {
	CreateCustomer(name, email, descr string) (*string, error)
	CreateNewPaymentMethod(ctx context.Context, teamID string) (*string, error)
}

type SubscriptionServiceI interface {
	Create(ctx context.Context, userID string, packageID string, teamID string) error
}

type ProjectsServiceI interface {
}

type Service struct {
	UserService         UserServiceI
	MemberService       MemberServiceI
	TasksService        TasksServiceI
	RolesService        RolesServiceI
	ProjectsService     ProjectsServiceI
	TeamsService        TeamsServiceI
	AiChatService       AiChatServiceI
	PackagesService     PackagesServiceI
	PaymentService      PaymentServiceI
	SubscriptionService SubscriptionServiceI
	FilesService        FilesServiceI
}

type Deps struct {
	Hasher                 *hash.Hasher
	UserRepository         repository.UserRepository
	MemberRepository       repository.MemberRepository
	TasksRepository        repository.TasksRepository
	ProjectsRepository     repository.ProjectsRepository
	TeamsRepository        repository.TeamsRepository
	RolesRepository        repository.RolesRepository
	PackagesRepository     repository.PackagesRepository
	FilesRepository        repository.FilesRepository
	FolderRepository       repository.FolderRepository
	SubscriptionRepository repository.SubscriptionRepository
	StorageProvider        storageProvider
	StorageEndpointURL     string
	JWTManager             *manager.JWTManager
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeTTL    time.Duration
	Sender                 email.Sender
	EmailConfig            config.EmailConfig
	CodeGenerator          *codegenerator.CodeGenerator
	TokenGenerator         *tokengenerator.TokenGenerator
	PaymentProvider        PaymentProvider
	OpenAI                 openAI
	MattermostAdapter      MattermostAdapter
	VerificationCodeLength int
	ApiUrl                 string
}

func NewService(deps *Deps) *Service {
	packagesService := NewPackagesService(deps.PackagesRepository)
	paymentService := NewPaymentService(deps.PaymentProvider, deps.TeamsRepository)
	filesService := NewFilesService(deps.StorageProvider, deps.FilesRepository, deps.CodeGenerator, deps.FolderRepository, deps.StorageEndpointURL, deps.UserRepository)
	emailService := NewEmailService(deps.Sender, deps.EmailConfig)
	rolesService := NewRolesService(deps.RolesRepository)
	teamsService := NewTeamsService(deps.TeamsRepository, deps.UserRepository, deps.MemberRepository, paymentService, *rolesService, filesService)
	userService := NewUserService(deps.Hasher, deps.UserRepository, filesService, deps.JWTManager, deps.AccessTokenTTL, deps.RefreshTokenTTL, deps.VerificationCodeTTL, deps.CodeGenerator, emailService, deps.MattermostAdapter, deps.VerificationCodeLength, deps.ApiUrl)
	subscriptionService := NewSubscriptionService(userService, deps.TeamsRepository, packagesService, deps.SubscriptionRepository, deps.PaymentProvider)
	return &Service{
		AiChatService:       NewAiChatService(deps.OpenAI),
		UserService:         userService,
		MemberService:       NewMemberService(deps.MemberRepository, deps.UserRepository, deps.TeamsRepository, deps.CodeGenerator, deps.TokenGenerator, teamsService, emailService, userService, deps.ApiUrl),
		TeamsService:        teamsService,
		RolesService:        rolesService,
		TasksService:        NewTasksService(deps.TasksRepository),
		ProjectsService:     NewProjectsService(deps.ProjectsRepository),
		FilesService:        filesService,
		PaymentService:      paymentService,
		SubscriptionService: subscriptionService,
		PackagesService:     packagesService,
	}
}
