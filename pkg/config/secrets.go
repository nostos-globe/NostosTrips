package config

import (
    "log"
    "sync"
)

var (
    secretsManager *SecretsManager
    once          sync.Once
)

type SecretsManager struct {
    vault *VaultConfig
}

func GetSecretsManager() *SecretsManager {
    once.Do(func() {
        vault, err := InitVault()
        if err != nil {
            log.Fatalf("Failed to initialize Vault: %v", err)
        }
        secretsManager = &SecretsManager{vault: vault}
    })
    return secretsManager
}

func (sm *SecretsManager) LoadSecrets() map[string]string {
    log.Printf("Starting to load secrets from Vault")
    log.Printf("Vault Server URL: %s", sm.vault.GetVaultURL())
    secrets := make(map[string]string)
    
    // All secrets are stored in the same path
    path := "kv/data/NostosAuthService"
    log.Printf("Using Vault path: %s", path)
    log.Printf("Full Vault URL: %s/%s", sm.vault.GetVaultURL(), path)
    
    // List of secrets to retrieve
    secretKeys := []string{
        "DB_HOST",
        "DB_USER",
        "DB_PASSWORD",
        "DB_NAME",
        "DB_PORT",
        "JWT_SECRET",
    }
    log.Printf("Attempting to load %d secrets", len(secretKeys))

    for _, key := range secretKeys {
        log.Printf("Fetching secret for key: %s", key)
        value, err := sm.vault.GetSecret(path, key)
        if err != nil {
            log.Printf("Warning: Failed to load secret for %s: %v", key, err)
            continue
        }
        secrets[key] = value
        log.Printf("Successfully loaded secret for: %s", key)
    }

    log.Printf("Completed loading secrets. Total secrets loaded: %d", len(secrets))
    return secrets
}