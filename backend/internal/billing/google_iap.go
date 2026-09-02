package billing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type GoogleIAPClient struct {
	serviceAccount string
	packageName    string
	httpClient     *http.Client
	logger         *zap.Logger
}

type GoogleProductPurchase struct {
	Kind                  string `json:"kind"`
	PurchaseTimeMillis    string `json:"purchaseTimeMillis"`
	ConsumptionState      string `json:"consumptionState"`
	DeveloperPayload      string `json:"developerPayload"`
	OrderID               string `json:"orderId"`
	PurchaseState         int    `json:"purchaseState"` // 0=purchased, 1=canceled, 2=pending
	PurchaseType          int    `json:"purchaseType"`
	AcknowledgementState  int    `json:"acknowledgementState"`
}

type GoogleSubscriptionPurchase struct {
	Kind                        string `json:"kind"`
	StartTimeMillis             string `json:"startTimeMillis"`
	ExpiryTimeMillis            string `json:"expiryTimeMillis"`
	AutoResumeTimeMillis        string `json:"autoResumeTimeMillis"`
	CancelReason                int    `json:"cancelReason"`
	UserCancellationTimeMillis  string `json:"userCancellationTimeMillis"`
	OrderID                     string `json:"orderId"`
	PurchaseState               int    `json:"purchaseState"` // 0=purchased, 1=canceled, 2=pending
	AutoRenewing                bool   `json:"autoRenewing"`
	PriceCurrencyCode           string `json:"priceCurrencyCode"`
	PriceAmountMicros           string `json:"priceAmountMicros"`
	CountryCode                 string `json:"countryCode"`
	AcknowledgementState        int    `json:"acknowledgementState"`
}

const googlePlayAPIBase = "https://androidpublisher.googleapis.com/androidpublisher/v3"

func NewGoogleIAPClient(serviceAccount, packageName string, logger *zap.Logger) *GoogleIAPClient {
	return &GoogleIAPClient{
		serviceAccount: serviceAccount,
		packageName:    packageName,
		httpClient:     &http.Client{Timeout: 30 * time.Second},
		logger:         logger,
	}
}

func (g *GoogleIAPClient) VerifyProductPurchase(ctx context.Context, productID, purchaseToken string) (*GoogleProductPurchase, error) {
	url := fmt.Sprintf("%s/applications/%s/purchases/products/%s/tokens/%s",
		googlePlayAPIBase, g.packageName, productID, purchaseToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// In production, use Google OAuth2 JWT for authentication
	req.Header.Set("Authorization", "Bearer "+g.getAccessToken())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api error: %s", string(body))
	}

	var purchase GoogleProductPurchase
	if err := json.Unmarshal(body, &purchase); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &purchase, nil
}

func (g *GoogleIAPClient) VerifySubscriptionPurchase(ctx context.Context, subscriptionID, purchaseToken string) (*GoogleSubscriptionPurchase, error) {
	url := fmt.Sprintf("%s/applications/%s/purchases/subscriptions/%s/tokens/%s",
		googlePlayAPIBase, g.packageName, subscriptionID, purchaseToken)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.getAccessToken())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api error: %s", string(body))
	}

	var purchase GoogleSubscriptionPurchase
	if err := json.Unmarshal(body, &purchase); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &purchase, nil
}

func (g *GoogleIAPClient) AcknowledgePurchase(ctx context.Context, productID, purchaseToken string) error {
	url := fmt.Sprintf("%s/applications/%s/purchases/products/%s/tokens/%s:acknowledge",
		googlePlayAPIBase, g.packageName, productID, purchaseToken)

	body := `{"developerPayload": ""}`
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+g.getAccessToken())
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("acknowledge failed: %s", string(respBody))
	}

	return nil
}

func (g *GoogleIAPClient) getAccessToken() string {
	// In production, generate a JWT signed with the service account key
	// and exchange it for a Google OAuth2 access token
	return "google-oauth-token"
}

func (g *GoogleIAPClient) ValidateSubscription(sub *GoogleSubscriptionPurchase) bool {
	// Check purchase state
	if sub.PurchaseState != 0 { // Not purchased
		return false
	}

	// Check expiry
	if sub.ExpiryTimeMillis != "" {
		expiryMs, err := parseInt64(sub.ExpiryTimeMillis)
		if err != nil {
			return false
		}
		if time.Now().UnixMilli() > expiryMs {
			return false
		}
	}

	return true
}

func parseInt64(s string) (int64, error) {
	var n int64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("not a number: %s", s)
		}
		n = n*10 + int64(c-'0')
	}
	return n, nil
}
