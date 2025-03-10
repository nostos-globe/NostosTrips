package config

import (
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
)

type VaultConfig struct {
	Client *vault.Client
}

func InitVault() (*VaultConfig, error) {
	config := vault.DefaultConfig()
	config.Address = getVaultAddress()

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %v", err)
	}

	// Set Vault token from environment variable
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("VAULT_TOKEN environment variable not set")
	}
	client.SetToken(token)

	return &VaultConfig{Client: client}, nil
}

func (vc *VaultConfig) GetSecret(path, key string) (string, error) {
	secret, err := vc.Client.Logical().Read(path)
	if err != nil {
		return "", fmt.Errorf("failed to read secret: %v", err)
	}

	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("secret not found at path: %s", path)
	}

	// In KV v2, secret data is under the "data" key
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid secret format at path: %s", path)
	}

	value, ok := data[key].(string)
	if !ok {
		return "", fmt.Errorf("key not found or not a string: %s", key)
	}

	return value, nil
}

func (vc *VaultConfig) GetVaultURL() string {
	return vc.Client.Address()
}

func getVaultAddress() string {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return "http://localhost:8200"
	}
	return addr
}
