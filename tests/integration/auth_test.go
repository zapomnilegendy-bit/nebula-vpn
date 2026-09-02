package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nebula-vpn/backend/internal/auth"
)

func TestJWTTokenLifecycle(t *testing.T) {
	manager := auth.NewJWTManager(
		"test-access-secret-32chars!",
		"test-refresh-secret-32chars!",
		15*time.Minute,
		7*24*time.Hour,
		"nebula-vpn-test",
	)

	userID := uuid.New()
	email := "test@nebula-vpn.io"

	// Generate tokens
	pair, err := manager.GenerateTokenPair(userID, email)
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	if pair.AccessToken == "" {
		t.Fatal("access token should not be empty")
	}
	if pair.RefreshToken == "" {
		t.Fatal("refresh token should not be empty")
	}
	if pair.ExpiresIn <= 0 {
		t.Fatal("expires_in should be positive")
	}

	// Validate access token
	claims, err := manager.ValidateAccessToken(pair.AccessToken)
	if err != nil {
		t.Fatalf("failed to validate access token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected user_id %s, got %s", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}

	// Validate refresh token
	refreshClaims, err := manager.ValidateRefreshToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("failed to validate refresh token: %v", err)
	}
	if refreshClaims.UserID != userID {
		t.Errorf("expected user_id %s, got %s", userID, refreshClaims.UserID)
	}
}

func TestTokenRefresh(t *testing.T) {
	manager := auth.NewJWTManager(
		"test-access-secret-32chars!",
		"test-refresh-secret-32chars!",
		15*time.Minute,
		7*24*time.Hour,
		"nebula-vpn-test",
	)

	userID := uuid.New()
	email := "test@nebula-vpn.io"

	pair, err := manager.GenerateTokenPair(userID, email)
	if err != nil {
		t.Fatalf("failed to generate token pair: %v", err)
	}

	// Refresh using the refresh token
	newPair, err := manager.RefreshAccessToken(pair.RefreshToken)
	if err != nil {
		t.Fatalf("failed to refresh access token: %v", err)
	}

	if newPair.AccessToken == pair.AccessToken {
		t.Error("refreshed access token should be different from the original")
	}

	// Validate new access token
	_, err = manager.ValidateAccessToken(newPair.AccessToken)
	if err != nil {
		t.Fatalf("failed to validate new access token: %v", err)
	}
}

func TestInvalidTokenRejection(t *testing.T) {
	manager := auth.NewJWTManager(
		"test-access-secret-32chars!",
		"test-refresh-secret-32chars!",
		15*time.Minute,
		7*24*time.Hour,
		"nebula-vpn-test",
	)

	// Test expired token (manually construct)
	_, err := manager.ValidateAccessToken("invalid-token")
	if err == nil {
		t.Error("expected error for invalid token")
	}

	// Test access token used as refresh token
	userID := uuid.New()
	pair, err := manager.GenerateTokenPair(userID, "test@test.com")
	if err != nil {
		t.Fatalf("failed to generate tokens: %v", err)
	}

	_, err = manager.ValidateRefreshToken(pair.AccessToken)
	if err == nil {
		t.Error("expected error when using access token as refresh token")
	}
}
