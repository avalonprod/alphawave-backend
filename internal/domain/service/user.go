package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Coke15/AlphaWave-BackEnd/internal/apperrors"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/model"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/repository"
	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/auth/manager"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/codegenerator"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/hash"
	"github.com/Coke15/AlphaWave-BackEnd/pkg/logger"
)

type UserService struct {
	hasher                 *hash.Hasher
	repository             repository.UserRepository
	filesService           FilesServiceI
	JWTManager             *manager.JWTManager
	AccessTokenTTL         time.Duration
	RefreshTokenTTL        time.Duration
	VerificationCodeTTL    time.Duration
	VerificationCodeLength int
	ApiUrl                 string
	codeGenerator          *codegenerator.CodeGenerator
	emailService           *EmailService
	mattermostAdapter      MattermostAdapter
}

func NewUserService(hasher *hash.Hasher, repository repository.UserRepository, filesService FilesServiceI, JWTManager *manager.JWTManager, accessTokenTTL time.Duration, refreshTokenTTL time.Duration, verificationCodeTTL time.Duration, codeGenerator *codegenerator.CodeGenerator, emailService *EmailService, mattermostAdapter MattermostAdapter, verificationCodeLength int, apiUrl string) *UserService {
	return &UserService{
		hasher:                 hasher,
		repository:             repository,
		filesService:           filesService,
		JWTManager:             JWTManager,
		AccessTokenTTL:         accessTokenTTL,
		RefreshTokenTTL:        refreshTokenTTL,
		emailService:           emailService,
		mattermostAdapter:      mattermostAdapter,
		VerificationCodeTTL:    verificationCodeTTL,
		codeGenerator:          codeGenerator,
		VerificationCodeLength: verificationCodeLength,
		ApiUrl:                 apiUrl,
	}
}

func (s *UserService) SignUp(ctx context.Context, input types.UserSignUpDTO) error {
	if err := validateCredentials(input.Email, input.Password); err != nil {
		return err
	}
	if err := validateUserData(input.FirstName, input.LastName, input.JobTitle); err != nil {
		return err
	}
	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}
	verificationCodeHash, err := s.hasher.Hash(input.Email)
	if err != nil {
		return err
	}
	verificationCode := fmt.Sprintf("%s%s", s.codeGenerator.RandomSecret(s.VerificationCodeLength), verificationCodeHash)

	mattermostPassowrdHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	user := model.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		JobTitle:  input.JobTitle,
		Email:     input.Email,
		Password:  passwordHash,
		MattermostData: model.MattermostData{
			Email:    input.Email,
			Password: mattermostPassowrdHash[8:],
		},
		Verification: model.UserVerificationPayload{
			VerificationCode:            verificationCode,
			VerificationCodeExpiresTime: time.Now().Add(s.VerificationCodeTTL),
		},
		RegisteredTime: time.Now(),
		LastVisitTime:  time.Now(),
	}

	isDuplicate, err := s.repository.IsDuplicate(ctx, input.Email)
	if err != nil {
		return err
	}

	if isDuplicate {
		return apperrors.ErrUserAlreadyExists
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return err
	}

	// if err := s.mattermostAdapter.CreateUser(ctx, types.CreateUserMattermostPayloadDTO{
	// 	Email:     input.Email,
	// 	Username:  input.Email[:strings.Index(input.Email, "@")],
	// 	FirstName: input.FirstName,
	// 	LastName:  input.LastName,
	// 	Password:  mattermostPassowrdHash[8:],
	// }); err != nil {
	// 	if err := s.repository.DeleteUserByEmail(ctx, input.Email); err != nil {
	// 		return err
	// 	}

	// 	return err
	// }

	go func() {
		err = s.emailService.SendUserVerificationEmail(VerificationEmailInput{
			Name:  input.FirstName,
			Email: input.Email,
			URL:   s.ApiUrl + "/api/v1/users/verify/" + verificationCode,
		})
		logger.Error(err)
	}()

	return nil
}

func (s *UserService) SignIn(ctx context.Context, input types.UserSignInDTO) (types.AuthPayload, error) {

	passwordHash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return types.AuthPayload{}, err
	}

	user, err := s.repository.GetByСredentials(ctx, input.Email, passwordHash)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.AuthPayload{}, err
		}
		return types.AuthPayload{}, err
	}

	if !user.Verification.Verified {
		return types.AuthPayload{}, apperrors.ErrUserNotVerifyed
	}
	tokens, err := s.createSession(ctx, user.ID)
	if err != nil {
		return types.AuthPayload{}, err
	}
	return types.AuthPayload{
		UserId: user.ID,
		UserInfo: types.UserDTO{
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			JobTitle:       user.JobTitle,
			Email:          user.Email,
			Verification:   user.Verification.Verified,
			RegisteredTime: user.RegisteredTime,
			LastVisitTime:  user.LastVisitTime,
			Blocked:        user.Blocked,
			Settings: types.UserSettings{
				UserIcon: types.UserImageDTO{
					Url:              user.Settings.UserIcon.Url,
					LastModifiedTime: user.Settings.UserIcon.LastModifiedTime,
				},
				BannerImage: types.UserImageDTO{
					Url:              user.Settings.BannerImage.Url,
					LastModifiedTime: user.Settings.BannerImage.LastModifiedTime,
				},
				TimeZone:   user.Settings.TimeZone,
				DateFormat: user.Settings.DateFormat,
				TimeFormat: user.Settings.TimeFormat,
			},
		},
		Tokens: tokens,
	}, nil
}

func (s *UserService) LogOut(ctx context.Context, userID string) error {
	return s.repository.RemoveSession(ctx, userID)
}

func (s *UserService) EnableTwoFactorAuth(ctx context.Context, userID string) error {
	// TODO
	return nil
}

func (s *UserService) GetUserById(ctx context.Context, userID string) (types.UserDTO, error) {
	res, err := s.repository.GetUserById(ctx, userID)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.UserDTO{}, err
		}
		return types.UserDTO{}, err
	}

	user := types.UserDTO{
		FirstName:      res.FirstName,
		LastName:       res.LastName,
		JobTitle:       res.JobTitle,
		Email:          res.Email,
		LastVisitTime:  res.LastVisitTime,
		RegisteredTime: res.RegisteredTime,
		Verification:   res.Verification.Verified,
		Blocked:        res.Blocked,
		Settings: types.UserSettings{
			UserIcon: types.UserImageDTO{
				Url:              res.Settings.UserIcon.Url,
				LastModifiedTime: res.Settings.UserIcon.LastModifiedTime,
			},
			BannerImage: types.UserImageDTO{
				Url:              res.Settings.BannerImage.Url,
				LastModifiedTime: res.Settings.BannerImage.LastModifiedTime,
			},
			TimeZone:   res.Settings.TimeZone,
			DateFormat: res.Settings.DateFormat,
			TimeFormat: res.Settings.TimeFormat,
		},
	}

	return user, nil
}

func (s *UserService) Verify(ctx context.Context, verificationCode string) (types.Tokens, error) {
	user, err := s.repository.GetUserByVerificationCode(ctx, verificationCode)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.Tokens{}, err
		}
		return types.Tokens{}, err
	}
	if user.Verification.Verified {
		return types.Tokens{}, apperrors.ErrUserAlreadyVerifyed
	}
	if user.Verification.VerificationCode != verificationCode {
		return types.Tokens{}, apperrors.ErrIncorrectVerificationCode
	}
	if user.Verification.VerificationCodeExpiresTime.UTC().Unix() < time.Now().UTC().Unix() {
		return types.Tokens{}, apperrors.ErrVerificationCodeExpired
	}

	err = s.repository.Verify(ctx, verificationCode)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.Tokens{}, apperrors.ErrUserNotFound
		}
		return types.Tokens{}, err
	}
	return s.createSession(ctx, user.ID)
}

func (s *UserService) ResendVerificationCode(ctx context.Context, email string) error {

	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}

	verificationCodeHash, err := s.hasher.Hash(email)
	if err != nil {
		return err
	}
	verificationCode := fmt.Sprintf("%s%s", s.codeGenerator.RandomSecret(s.VerificationCodeLength), verificationCodeHash)

	verificationPayload := model.UserVerificationPayload{
		VerificationCode:            verificationCode,
		VerificationCodeExpiresTime: time.Now().Add(s.VerificationCodeTTL),
	}
	err = s.repository.ChangeVerificationCode(ctx, email, verificationPayload)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}

	go func() {
		err = s.emailService.SendUserVerificationEmail(VerificationEmailInput{
			Name:  user.FirstName,
			Email: user.Email,
			URL:   s.ApiUrl + "/api/v1/users/verify/" + verificationCode,
		})
		logger.Error(err)
	}()
	return nil
}

func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (types.AuthPayload, error) {
	user, err := s.repository.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return types.AuthPayload{}, err
		}
		return types.AuthPayload{}, err
	}
	if user.Blocked {
		return types.AuthPayload{}, apperrors.ErrUserBlocked
	}
	tokens, err := s.createSession(ctx, user.ID)
	if err != nil {
		return types.AuthPayload{}, err
	}
	return types.AuthPayload{
		UserId: user.ID,
		UserInfo: types.UserDTO{
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			JobTitle:       user.JobTitle,
			Email:          user.Email,
			Verification:   user.Verification.Verified,
			RegisteredTime: user.RegisteredTime,
			LastVisitTime:  user.LastVisitTime,
			Blocked:        user.Blocked,
			Settings: types.UserSettings{
				UserIcon: types.UserImageDTO{
					Url:              user.Settings.UserIcon.Url,
					LastModifiedTime: user.Settings.UserIcon.LastModifiedTime,
				},
				BannerImage: types.UserImageDTO{
					Url:              user.Settings.BannerImage.Url,
					LastModifiedTime: user.Settings.BannerImage.LastModifiedTime,
				},
				TimeZone:   user.Settings.TimeZone,
				DateFormat: user.Settings.DateFormat,
				TimeFormat: user.Settings.TimeFormat,
			},
		},
		Tokens: tokens,
	}, nil
}

func (s *UserService) createSession(ctx context.Context, userID string) (types.Tokens, error) {

	accessToken, err := s.JWTManager.NewJWT(userID, s.AccessTokenTTL)
	if err != nil {
		return types.Tokens{}, err
	}
	refreshToken, err := s.JWTManager.NewRefreshToken()
	if err != nil {
		return types.Tokens{}, err
	}
	tokens := model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	session := model.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresTime:  time.Now().Add(s.RefreshTokenTTL),
	}
	// user, err := s.repository.GetUserById(ctx, userID)
	// if err != nil {
	// 	if errors.Is(err, apperrors.ErrUserNotFound) {
	// 		return types.Tokens{}, apperrors.ErrUserNotFound
	// 	}
	// 	return types.Tokens{}, err
	// }

	// token, err := s.mattermostAdapter.SignIn(user.MattermostData.Email, user.MattermostData.Password)
	token := ""

	if err != nil {
		return types.Tokens{}, err
	}

	err = s.repository.SetSession(ctx, userID, session, time.Now())

	return types.Tokens{
		AccessToken:     tokens.AccessToken,
		RefreshToken:    tokens.RefreshToken,
		MattermostToken: token,
	}, err
}

func (s *UserService) ChangePassword(ctx context.Context, userID, newPassword, oldPassword string) error {

	passwordHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	oldPasswordHash, err := s.hasher.Hash(oldPassword)
	if err != nil {
		return err
	}

	err = s.repository.ChangePassword(ctx, userID, passwordHash, oldPasswordHash)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return err
		}
		return err
	}

	return nil
}

func (s *UserService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	if user.Blocked {
		return apperrors.ErrUserBlocked
	}
	tokenHash, err := s.hasher.Hash(user.Email)
	if err != nil {
		return err
	}

	result := fmt.Sprintf("%s.%s", s.codeGenerator.RandomSecret(30), tokenHash)

	tokenExpiresTime := time.Now().Add(time.Hour * 1)

	err = s.repository.SetForgotPassword(ctx, user.Email, model.ForgotPasswordPayload{
		Token:            tokenHash,
		ResultToken:      result,
		TokenExpiresTime: tokenExpiresTime,
	})

	if err != nil {
		return err
	}
	if err := s.emailService.SendUserForgotPassword(ForgotPasswordInput{
		Email:            user.Email,
		TokenExpiresTime: 1,
		URL:              s.ApiUrl + fmt.Sprintf("/api/v1/users/forgot-password-verify?email=%s&token=%s&result=%s", user.Email, tokenHash, result),
	}); err != nil {
		return err
	}

	return nil
}

func (s *UserService) VerifyForgotPasswordToken(ctx context.Context, email, token, tokenResult string) (types.ForgotPasswordPayloadDTO, error) {

	user, err := s.repository.GetByForgotPasswordToken(ctx, token, tokenResult)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {

			return types.ForgotPasswordPayloadDTO{}, apperrors.ErrUserNotFound
		}
		return types.ForgotPasswordPayloadDTO{}, err
	}

	if user.Email != email {

		return types.ForgotPasswordPayloadDTO{}, apperrors.ErrUserNotFound
	}

	return types.ForgotPasswordPayloadDTO{
		Email:       user.Email,
		Token:       token,
		ResultToken: tokenResult,
	}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, email, token, tokenResult, password string) error {
	user, err := s.repository.GetByForgotPasswordToken(ctx, token, tokenResult)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	if user.Email != email {
		return apperrors.ErrUserNotFound
	}

	passwordHash, err := s.hasher.Hash(password)
	if err != nil {
		return err
	}

	err = s.repository.ResetPassword(ctx, token, user.Email, passwordHash)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserInfo(ctx context.Context, userID string, input types.UpdateUserInfoDTO) error {
	if err := s.repository.UpdateUserInfo(ctx, userID, model.UpdateUserInfoInput{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		JobTitle:  input.JobTitle,
		Email:     input.Email,
	}); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *UserService) UpdateUserSettings(ctx context.Context, userID string, input types.UpdateUserSettingsDTO) error {
	if err := s.repository.UpdateUserSettings(ctx, userID, model.UpdateUserSettingsInput{
		UserIcon: &model.UserImage{
			Url:              input.UserIcon.Url,
			Path:             input.UserIcon.Path,
			LastModifiedTime: time.Now(),
		},
		BannerImage: &model.UserImage{
			Url:              input.BannerImage.Url,
			Path:             input.BannerImage.Path,
			LastModifiedTime: time.Now(),
		},
		TimeZone:   input.TimeZone,
		DateFormat: input.DateFormat,
		TimeFormat: input.TimeFormat,
	}); err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (s *UserService) UploadUserAvatar(ctx context.Context, userID string, reader io.Reader, extension string, fileName string, size int) (string, error) {
	user, err := s.repository.GetUserById(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.Settings.UserIcon.SetUp {
		err := s.filesService.DeleteImage(ctx, user.Settings.UserIcon.Path)
		if err != nil {
			return "", err
		}
	}
	image, err := s.filesService.UploadImage(ctx, reader, types.UploadImageDTO{
		FileName:  fileName,
		Extension: extension,
		Size:      size,
	})
	if err != nil {
		return "", err
	}
	if err = s.repository.UpdateUserSettings(ctx, userID, model.UpdateUserSettingsInput{
		UserIcon: &model.UserImage{
			Url:              image.Url,
			Path:             image.Path,
			LastModifiedTime: time.Now(),
			SetUp:            true,
		},
	}); err != nil {
		return "", err
	}
	return image.Url, nil
}

func (s *UserService) UploadUserBanner(ctx context.Context, userID string, reader io.Reader, extension string, fileName string, size int) (string, error) {
	user, err := s.repository.GetUserById(ctx, userID)
	if err != nil {
		return "", err
	}

	if user.Settings.BannerImage.SetUp {
		err := s.filesService.DeleteImage(ctx, user.Settings.BannerImage.Path)
		if err != nil {
			return "", err
		}
	}
	image, err := s.filesService.UploadImage(ctx, reader, types.UploadImageDTO{
		FileName:  fileName,
		Extension: extension,
		Size:      size,
	})
	if err != nil {
		return "", err
	}
	if err = s.repository.UpdateUserSettings(ctx, userID, model.UpdateUserSettingsInput{
		BannerImage: &model.UserImage{
			Url:              image.Url,
			Path:             image.Path,
			LastModifiedTime: time.Now(),
			SetUp:            true,
		},
	}); err != nil {
		return "", err
	}
	return image.Url, nil
}
