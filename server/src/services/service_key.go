package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const PublicServiceKey = "sk_etheria_public_a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"

type ServiceKey struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	KeyHash   string    `json:"-"`
	Scope     string    `json:"scope"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	LastUsed  time.Time `json:"last_used,omitempty"`
}

type ServiceKeyService struct {
	keys map[string]*ServiceKey
}

func NewServiceKeyService() *ServiceKeyService {
	s := &ServiceKeyService{
		keys: make(map[string]*ServiceKey),
	}
	s.generateDefaultKeys()
	s.registerPublicKey()
	s.printKeys()
	return s
}

func (s *ServiceKeyService) registerPublicKey() {
	serviceKey := &ServiceKey{
		ID:        "svc_public",
		Name:      "Public Frontend Key",
		Key:       PublicServiceKey,
		KeyHash:   hashServiceKey(PublicServiceKey),
		Scope:     "read,write",
		Active:    true,
		CreatedAt: time.Now(),
	}
	s.keys[PublicServiceKey] = serviceKey
}

func (s *ServiceKeyService) printKeys() {
	fmt.Println("\n\033[1;36m═══════════════════════════════════════════════════════════\033[0m")
	fmt.Println("\033[1;36m           ETHERIA - SERVICE KEYS\033[0m")
	fmt.Println("\033[1;36m═══════════════════════════════════════════════════════════\033[0m")
	for _, key := range s.keys {
		scope := key.Scope
		if key.ID == "svc_public" {
			scope = "\033[1;33mPUBLIC\033[0m (Frontend)"
		}
		fmt.Printf("\033[1;32m[%s]\033[0m \033[1;37m%s\033[0m\n", key.ID, key.Key)
		fmt.Printf("       Name: %s | Scope: %s\n", key.Name, scope)
	}
	fmt.Println("\033[1;36m═══════════════════════════════════════════════════════════\033[0m\n")
}

func (s *ServiceKeyService) generateDefaultKeys() {
	defaultKeys := []struct {
		name  string
		scope string
	}{
		{"Production Frontend", "read,write"},
		{"Development", "read,write"},
		{"Admin Panel", "admin"},
	}

	for i, dk := range defaultKeys {
		envKey := ""
		switch i {
		case 0:
			envKey = "SERVICE_KEY_PROD"
		case 1:
			envKey = "SERVICE_KEY_DEV"
		case 2:
			envKey = "SERVICE_KEY_ADMIN"
		}

		var key string
		if envKey != "" {
			key = generateServiceKey(envKey)
		} else {
			key = generateServiceKey(dk.name)
		}

		serviceKey := &ServiceKey{
			ID:        fmt.Sprintf("svc_%d", i+1),
			Name:      dk.name,
			Key:       key,
			KeyHash:   hashServiceKey(key),
			Scope:     dk.scope,
			Active:    true,
			CreatedAt: time.Now(),
		}
		s.keys[key] = serviceKey
	}
}

func (s *ServiceKeyService) GenerateKey(name string, scope string) (*ServiceKey, error) {
	key := generateServiceKey(name)
	hash := hashServiceKey(key)

	for s.keys[key] != nil {
		key = generateServiceKey(name + "_" + randomString(4))
		hash = hashServiceKey(key)
	}

	serviceKey := &ServiceKey{
		ID:        fmt.Sprintf("svc_%s", randomString(8)),
		Name:      name,
		Key:       key,
		KeyHash:   hash,
		Scope:     scope,
		Active:    true,
		CreatedAt: time.Now(),
	}

	s.keys[key] = serviceKey
	return serviceKey, nil
}

func (s *ServiceKeyService) ValidateKey(key string) (*ServiceKey, error) {
	if key == "" {
		return nil, fmt.Errorf("service key is required")
	}

	if !strings.HasPrefix(key, "sk_") {
		return nil, fmt.Errorf("invalid service key format")
	}

	serviceKey, exists := s.keys[key]
	if !exists {
		return nil, fmt.Errorf("service key not found")
	}

	if !serviceKey.Active {
		return nil, fmt.Errorf("service key is inactive")
	}

	if !serviceKey.ExpiresAt.IsZero() && time.Now().After(serviceKey.ExpiresAt) {
		return nil, fmt.Errorf("service key has expired")
	}

	serviceKey.LastUsed = time.Now()
	return serviceKey, nil
}

func (s *ServiceKeyService) RevokeKey(key string) error {
	if !strings.HasPrefix(key, "sk_") {
		return fmt.Errorf("invalid service key format")
	}

	serviceKey, exists := s.keys[key]
	if !exists {
		return fmt.Errorf("service key not found")
	}

	serviceKey.Active = false
	return nil
}

func (s *ServiceKeyService) ListKeys() []*ServiceKey {
	keys := make([]*ServiceKey, 0, len(s.keys))
	for _, key := range s.keys {
		keys = append(keys, key)
	}
	return keys
}

func (s *ServiceKeyService) GetKeyByID(id string) *ServiceKey {
	for _, key := range s.keys {
		if key.ID == id {
			return key
		}
	}
	return nil
}

func generateServiceKey(prefix string) string {
	randomPart := make([]byte, 24)
	rand.Read(randomPart)
	return fmt.Sprintf("sk_%s_%s", strings.ToLower(sanitizePrefix(prefix)), hex.EncodeToString(randomPart))
}

func hashServiceKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func sanitizePrefix(prefix string) string {
	replacer := strings.NewReplacer(" ", "_", "/", "_", "\\", "_")
	return replacer.Replace(strings.ToLower(prefix))
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	rand.Read(result)
	for i, b := range result {
		result[i] = charset[int(b)%len(charset)]
	}
	return string(result)
}
