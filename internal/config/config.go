package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/realalphabet/protonmail-client/internal/constants"
)

type Config struct {
	TOTPSecretKey string `json:"totpSecretKey"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(constants.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read configuration file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unable to parse configuration file: %w", err)
	}

	return &config, nil
}

func GenerateTOTPCode(secretKey string) (string, error) {
	return totp.GenerateCode(secretKey, time.Now())
}
