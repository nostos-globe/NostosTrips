package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthClient struct {
	BaseURL string
}

type TokenResponse struct {
	UserID uint   `json:"user_id"`
	Valid  bool   `json:"valid"`
	Error  string `json:"error,omitempty"`
}

func (c *AuthClient) ValidateToken(token string) (*TokenResponse, error) {
	print("" + c.BaseURL)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/validate", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

type UserResponse struct {
	Message string `json:"message"`
	User    struct {
		UserID              uint      `json:"user_id"`
		Email               string    `json:"email"`
		FailedLoginAttempts int       `json:"failed_login_attempts"`
		AccountLocked       bool      `json:"account_locked"`
		RegistrationDate    time.Time `json:"registration_date"`
	} `json:"user"`
}

func (c *AuthClient) GetUserID(token string) (uint, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/profile", c.BaseURL), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Cookie", fmt.Sprintf("auth_token=%s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get profile: status code %d", resp.StatusCode)
	}

	var userResponse UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return 0, err
	}

	return userResponse.User.UserID, nil
}
