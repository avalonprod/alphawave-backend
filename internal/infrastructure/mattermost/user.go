package mattermost

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Coke15/AlphaWave-BackEnd/internal/domain/types"
)

type MattermostAdapter struct {
	apiUrl string
}

func NewMattermostAdapter(apiUrl string) *MattermostAdapter {
	return &MattermostAdapter{
		apiUrl: apiUrl,
	}
}

// Inputs
type CreateUserInput struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type SignInUserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Response
type Response struct {
	Id            string `json:"id"`
	CreateAt      int    `json:"create_at"`
	UpdateAt      int    `json:"update_at"`
	DeleteAt      int    `json:"delete_at"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type errorResponse struct {
	StatusCode int    `json:"status_code"`
	Id         string `json:"id"`
	Message    string `json:"message"`
	// RequestId  string `json:"request_id"`
}

func (a *MattermostAdapter) CreateUser(ctx context.Context, input types.CreateUserMattermostPayloadDTO) error {
	client := &http.Client{}

	inputData := CreateUserInput{
		Email:     input.Email,
		Username:  input.Username,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  input.Password,
	}

	inputBytes, err := json.Marshal(inputData)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.apiUrl+"/users", bytes.NewBuffer(inputBytes))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {

		var errorResp errorResponse

		err = json.Unmarshal(body, &errorResp)

		if err != nil {
			return err
		}

		return fmt.Errorf("mattermost error: %s", errorResp.Message)
	}

	var output Response

	err = json.Unmarshal(body, &output)

	if err != nil {
		return err
	}

	return nil
}

func (a *MattermostAdapter) SignIn(email string, password string) (string, error) {
	client := &http.Client{}

	inputData := SignInUserInput{
		Email:    email,
		Password: password,
	}

	inputBytes, err := json.Marshal(inputData)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", a.apiUrl+"/users/login", bytes.NewBuffer(inputBytes))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusOK {

		var errorResp errorResponse

		err = json.Unmarshal(body, &errorResp)
		if err != nil {
			return "", err
		}
		return "", errors.New(errorResp.Message)
	}

	token := response.Header.Get("token")

	return token, nil
}
