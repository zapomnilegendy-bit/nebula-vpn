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

type AppleIAPClient struct {
	sharedSecret string
	httpClient   *http.Client
	logger       *zap.Logger
	environment  string // "sandbox" or "production"
}

type AppleVerifyReceiptRequest struct {
	ReceiptData string `json:"receipt-data"`
	Password    string `json:"password,omitempty"`
	ExcludeOldTransactions bool `json:"exclude-old-transactions,omitempty"`
}

type AppleVerifyReceiptResponse struct {
	Status    int    `json:"status"`
	Receipt   AppleReceipt `json:"receipt"`
	LatestReceiptInfo []AppleInAppPurchase `json:"latest_receipt_info"`
	LatestReceipt     string `json:"latest_receipt"`
}

type AppleReceipt struct {
	BundleID              string `json:"bundle_id"`
	ApplicationVersion    string `json:"application_version"`
	OriginalApplicationVersion string `json:"original_application_version"`
	InApp []AppleInAppPurchase `json:"in_app"`
}

type AppleInAppPurchase struct {
	Quantity              string `json:"quantity"`
	ProductID             string `json:"product_id"`
	TransactionID         string `json:"transaction_id"`
	OriginalTransactionID string `json:"original_transaction_id"`
	PurchaseDate          string `json:"purchase_date"`
	ExpiresDate           string `json:"expires_date"`
	IsTrialPeriod         string `json:"is_trial_period"`
	IsInIntroOfferPeriod  string `json:"is_in_intro_offer_period"`
	CancellationDate      string `json:"cancellation_date"`
	WebOrderLineItemID    string `json:"web_order_line_item_id"`
}

const (
	appleProductionURL = "https://buy.itunes.apple.com/verifyReceipt"
	appleSandboxURL    = "https://sandbox.itunes.apple.com/verifyReceipt"
)

func NewAppleIAPClient(sharedSecret, environment string, logger *zap.Logger) *AppleIAPClient {
	return &AppleIAPClient{
		sharedSecret: sharedSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		logger:       logger,
		environment:  environment,
	}
}

func (a *AppleIAPClient) VerifyReceipt(ctx context.Context, receiptData string) (*AppleVerifyReceiptResponse, error) {
	reqBody := AppleVerifyReceiptRequest{
		ReceiptData: receiptData,
		Password:    a.sharedSecret,
		ExcludeOldTransactions: true,
	}

	resp, err := a.verifyReceiptWithURL(ctx, appleProductionURL, reqBody)
	if err != nil {
		return nil, err
	}

	// If status is 21007, receipt is from sandbox — retry
	if resp.Status == 21007 && a.environment == "sandbox" {
		a.logger.Info("production verification failed (21007), retrying with sandbox")
		return a.verifyReceiptWithURL(ctx, appleSandboxURL, reqBody)
	}

	return resp, nil
}

func (a *AppleIAPClient) verifyReceiptWithURL(ctx context.Context, url string, reqBody AppleVerifyReceiptRequest) (*AppleVerifyReceiptResponse, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	var result AppleVerifyReceiptResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	if result.Status != 0 {
		return &result, fmt.Errorf("apple receipt verification failed with status %d", result.Status)
	}

	return &result, nil
}

func (a *AppleIAPClient) ValidatePurchase(subscription *AppleInAppPurchase) bool {
	// Check if the purchase has expired
	if subscription.ExpiresDate != "" {
		expiresTime, err := time.Parse(time.RFC3339, subscription.ExpiresDate)
		if err != nil {
			a.logger.Warn("failed to parse apple expires_date", zap.Error(err))
			return false
		}
		if time.Now().After(expiresTime) {
			return false
		}
	}

	// Check for cancellation
	if subscription.CancellationDate != "" {
		return false
	}

	return true
}
