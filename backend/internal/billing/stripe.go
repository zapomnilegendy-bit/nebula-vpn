package billing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nebula-vpn/backend/internal/config"
	"go.uber.org/zap"
)

type StripeClient struct {
	httpClient    *http.Client
	secretKey     string
	webhookSecret string
	logger        *zap.Logger
	baseURL       string
}

type StripePrice struct {
	ID       string `json:"id"`
	Currency string `json:"currency"`
	UnitAmount int64 `json:"unit_amount"`
	Product   string `json:"product"`
}

type StripeSubscription struct {
	ID        string `json:"id"`
	Customer  string `json:"customer"`
	Status    string `json:"status"`
	CurrentPeriodStart int64 `json:"current_period_start"`
	CurrentPeriodEnd   int64 `json:"current_period_end"`
	CancelAtPeriodEnd  bool  `json:"cancel_at_period_end"`
	Items    struct {
		Data []struct {
			Price StripePrice `json:"price"`
		} `json:"data"`
	} `json:"items"`
}

type StripeCheckoutSession struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Status string `json:"status"`
}

func NewStripeClient(cfg *config.StripeConfig, logger *zap.Logger) *StripeClient {
	return &StripeClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		secretKey:     cfg.SecretKey,
		webhookSecret: cfg.WebhookSecret,
		logger:        logger,
		baseURL:       "https://api.stripe.com/v1",
	}
}

func (s *StripeClient) CreateCheckoutSession(ctx context.Context, priceID string, userID string, successURL, cancelURL string) (*StripeCheckoutSession, error) {
	body := fmt.Sprintf(
		"mode=subscription&line_items[0][price]=%s&line_items[0][quantity]=1&client_reference_id=%s&success_url=%s&cancel_url=%s",
		priceID, userID, successURL, cancelURL,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", s.baseURL+"/checkout/sessions", bytes.NewBufferString(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+s.secretKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("stripe request failed: %w", err)
	}
	defer resp.Body.Close()

	var session StripeCheckoutSession
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to decode stripe response: %w", err)
	}

	s.logger.Info("stripe checkout session created",
		zap.String("session_id", session.ID),
		zap.String("user_id", userID),
	)

	return &session, nil
}

func (s *StripeClient) GetSubscription(ctx context.Context, subscriptionID string) (*StripeSubscription, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.baseURL+"/subscriptions/"+subscriptionID, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.secretKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sub StripeSubscription
	if err := json.NewDecoder(resp.Body).Decode(&sub); err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *StripeClient) CancelSubscription(ctx context.Context, subscriptionID string) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", s.baseURL+"/subscriptions/"+subscriptionID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+s.secretKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stripe cancel failed: %s", string(body))
	}

	return nil
}

type StripeWebhookEvent struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Data    struct {
		Object json.RawMessage `json:"object"`
	} `json:"data"`
}

func (s *StripeClient) VerifyWebhook(payload []byte, signature string) (*StripeWebhookEvent, error) {
	// In production, use stripe-go library with proper signature verification
	var event StripeWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}
	return &event, nil
}
