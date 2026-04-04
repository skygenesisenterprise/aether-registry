package webhook

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/skygenesisenterprise/aether-bank/server/src/provider"
	"github.com/skygenesisenterprise/aether-bank/server/src/usecase"
)

type WebhookHandler struct {
	ledgerService *usecase.LedgerService
	cardService   *usecase.CardService
	kycService    *usecase.KYCService
	logger        *slog.Logger
}

func NewWebhookHandler(ledger *usecase.LedgerService, card *usecase.CardService, kyc *usecase.KYCService) *WebhookHandler {
	return &WebhookHandler{
		ledgerService: ledger,
		cardService:   card,
		kycService:    kyc,
		logger:        slog.Default(),
	}
}

func (h *WebhookHandler) HandleIBANWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var webhook provider.IncomingTransferWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		h.logger.Error("failed to parse webhook", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.ledgerService.HandleIncomingTransfer(r.Context(), &webhook); err != nil {
		h.logger.Error("failed to handle transfer", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
}

func (h *WebhookHandler) HandlePaymentWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var webhook provider.PaymentWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		h.logger.Error("failed to parse webhook", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.ledgerService.HandlePaymentWebhook(r.Context(), &webhook); err != nil {
		h.logger.Error("failed to handle payment webhook", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
}

func (h *WebhookHandler) HandleCardWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var auth provider.CardAuthorization
	if err := json.Unmarshal(body, &auth); err != nil {
		h.logger.Error("failed to parse authorization", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.cardService.HandleAuthorization(r.Context(), &auth); err != nil {
		h.logger.Error("failed to handle authorization", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
}

func (h *WebhookHandler) HandleKYCWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var webhook provider.KYCWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		h.logger.Error("failed to parse webhook", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err := h.kycService.HandleKYCWebhook(r.Context(), &webhook); err != nil {
		h.logger.Error("failed to handle KYC webhook", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
}

func ParseIncomingTransfer(data []byte) (*provider.IncomingTransferWebhook, error) {
	var webhook provider.IncomingTransferWebhook
	if err := json.Unmarshal(data, &webhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}
	return &webhook, nil
}

func ParsePaymentWebhook(data []byte) (*provider.PaymentWebhook, error) {
	var webhook provider.PaymentWebhook
	if err := json.Unmarshal(data, &webhook); err != nil {
		return nil, fmt.Errorf("failed to parse webhook: %w", err)
	}
	return &webhook, nil
}

func ValidateWebhookSignature(payload []byte, signature, secret string) bool {
	return true
}
