package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"

	"github.com/hanzoai/o11y-mcp/pkg/types"
)

const (
	EnvAPIKey = "O11Y-API-KEY"
	ContentType  = "Content-Type"
)

type Client struct {
	baseURL string
	apiKey  string
	logger  *zap.Logger
}

func NewClient(log *zap.Logger, url, apiKey string) *Client {
	return &Client{logger: log, baseURL: url, apiKey: apiKey}
}

func (s *Client) ListMetricKeys(ctx context.Context) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/metrics", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Making request to Client API", zap.String("method", "GET"), zap.String("endpoint", "/v1/o11y/metrics"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	s.logger.Debug("Successfully retrieved metric keys", zap.String("url", url), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) SearchMetricByText(ctx context.Context, searchText string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/autocomplete/aggregate_attributes?dataSource=metrics&searchText=%s", s.baseURL, searchText)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Searching metric names (aggregate_attributes)", zap.String("searchText", searchText))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.String("searchText", searchText), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	s.logger.Debug("Successfully searched metric names", zap.String("searchText", searchText), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) ListAlerts(ctx context.Context) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/alerts", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching alerts from SigNoz")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	s.logger.Debug("Successfully retrieved alerts", zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetAlertByRuleID(ctx context.Context, ruleID string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/rules/%s", s.baseURL, ruleID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching alert rule details", zap.String("ruleID", ruleID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.String("ruleID", ruleID), zap.Error(err))
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("read error: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	s.logger.Debug("Successfully retrieved alert rule", zap.String("ruleID", ruleID), zap.Int("status", resp.StatusCode))
	return body, nil
}

// ListDashboards filters data as it returns too much of data even the ui tags
// so we filter and only return required information which might help to get
// detailed info of a dashboard.
func (s *Client) ListDashboards(ctx context.Context) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/dashboards", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching dashboards from SigNoz")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var rawResponse map[string]interface{}
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		s.logger.Error("Failed to parse response", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if data, ok := rawResponse["data"].([]interface{}); ok {
		simplifiedDashboards := make([]map[string]interface{}, 0, len(data))

		for _, dashboard := range data {
			if dash, ok := dashboard.(map[string]interface{}); ok {
				var (
					dashData map[string]interface{}
					name     any
					desc     any
					tags     any
				)
				if v, ok := dash["data"].(map[string]interface{}); ok {
					dashData = v
					name = dashData["title"]
					desc = dashData["description"]
					tags = dashData["tags"]
				}

				simplified := map[string]interface{}{
					"uuid":        dash["id"],
					"name":        name,
					"description": desc,
					"tags":        tags,
					"createdAt":   dash["createdAt"],
					"updatedAt":   dash["updatedAt"],
					"createdBy":   dash["createdBy"],
					"updatedBy":   dash["updatedBy"],
				}
				simplifiedDashboards = append(simplifiedDashboards, simplified)
			}
		}

		simplifiedResponse := map[string]interface{}{
			"status": rawResponse["status"],
			"data":   simplifiedDashboards,
		}

		simplifiedJSON, err := json.Marshal(simplifiedResponse)
		if err != nil {
			s.logger.Error("Failed to marshal simplified response", zap.Error(err))
			return nil, fmt.Errorf("failed to marshal simplified response: %w", err)
		}

		s.logger.Debug("Successfully retrieved and simplified dashboards", zap.Int("count", len(simplifiedDashboards)), zap.Int("status", resp.StatusCode))
		return simplifiedJSON, nil
	}

	s.logger.Debug("Successfully retrieved dashboards", zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetDashboard(ctx context.Context, uuid string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/dashboards/%s", s.baseURL, uuid)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching dashboard details", zap.String("uuid", uuid))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.String("uuid", uuid), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved dashboard", zap.String("uuid", uuid), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) ListServices(ctx context.Context, start, end string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/services", s.baseURL)

	payload := map[string]string{
		"start": start,
		"end":   end,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	s.logger.Debug("Fetching services from SigNoz", zap.String("start", start), zap.String("end", end))

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.String("start", start), zap.String("end", end), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved services", zap.String("start", start), zap.String("end", end), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetServiceTopOperations(ctx context.Context, start, end, service string, tags json.RawMessage) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/service/top_operations", s.baseURL)

	payload := map[string]any{
		"start":   start,
		"end":     end,
		"service": service,
		"tags":    tags,
	}
	bodyBytes, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	s.logger.Debug("Fetching service top operations", zap.String("start", start), zap.String("end", end), zap.String("service", service))

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.String("start", start), zap.String("end", end), zap.String("service", service), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved service top operations", zap.String("start", start), zap.String("end", end), zap.String("service", service), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) QueryBuilderV5(ctx context.Context, body []byte) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/query_range", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("sending request", zap.String("url", url), zap.ByteString("body", body))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return b, nil
}

func (s *Client) GetAlertHistory(ctx context.Context, ruleID string, req types.AlertHistoryRequest) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/rules/%s/history/timeline", s.baseURL, ruleID)
	// includes ruleid to get history
	// eg: /v1/o11y/rules/<ruleID>/history/timeline

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set(ContentType, "application/json")
	httpReq.Header.Set(EnvAPIKey, s.apiKey)

	s.logger.Debug("sending request", zap.String("url", url), zap.ByteString("body", reqBody))

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	httpReq = httpReq.WithContext(ctx)

	s.logger.Debug("Making request to Client API",
		zap.String("method", "POST"),
		zap.String("endpoint", fmt.Sprintf("/v1/o11y/rules/%s/history/timeline", ruleID)))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.Error(err))
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		s.logger.Error("API request failed",
			zap.String("url", url),
			zap.Int("status", resp.StatusCode),
			zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (s *Client) ListLogViews(ctx context.Context) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/explorer/views?sourcePage=logs", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching log views from SigNoz")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved log views", zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetLogView(ctx context.Context, viewID string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/explorer/views/%s", s.baseURL, viewID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching log view details", zap.String("viewID", viewID))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.String("viewID", viewID), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved log view", zap.String("viewID", viewID), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetTraceFieldValues(ctx context.Context, fieldName string, searchText string) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/fields/values?signal=traces&name=%s&searchText=%s&metricName=&source=meter", s.baseURL, fieldName, searchText)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching trace field values", zap.String("fieldName", fieldName), zap.String("searchText", searchText))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", url), zap.String("fieldName", fieldName), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", url), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", url), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved trace field values", zap.String("fieldName", fieldName), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetTraceAvailableFields(ctx context.Context, searchText string) (json.RawMessage, error) {
	encodedSearchText := url.QueryEscape(searchText)
	urlStr := fmt.Sprintf("%s/v1/o11y/autocomplete/attribute_keys?aggregateOperator=noop&searchText=%s&dataSource=traces&aggregateAttribute=&tagType=", s.baseURL, encodedSearchText)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching trace available fields", zap.String("searchText", searchText))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", urlStr), zap.String("searchText", searchText), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", urlStr), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", urlStr), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved trace available fields", zap.String("searchText", searchText), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetLogsAvailableFields(ctx context.Context, searchText string) (json.RawMessage, error) {
	encodedSearchText := url.QueryEscape(searchText)
	urlStr := fmt.Sprintf("%s/v1/o11y/filter_suggestions?searchText=%s&dataSource=logs&existingFilter=e30", s.baseURL, encodedSearchText)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching logs available fields", zap.String("searchText", searchText))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", urlStr), zap.String("searchText", searchText), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", urlStr), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", urlStr), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved logs available fields", zap.String("searchText", searchText), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetMetricsAvailableFields(ctx context.Context, searchText string) (json.RawMessage, error) {
	encodedSearchText := url.QueryEscape(searchText)
	urlStr := fmt.Sprintf("%s/v1/o11y/autocomplete/aggregate_attributes?aggregateOperator=avg&searchText=%s&dataSource=metrics", s.baseURL, encodedSearchText)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching metrics available fields", zap.String("searchText", searchText))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", urlStr), zap.String("searchText", searchText), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", urlStr), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", urlStr), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved metrics available fields", zap.String("searchText", searchText), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetLogsFieldValues(ctx context.Context, fieldName string, searchText string) (json.RawMessage, error) {
	encodedFieldName := url.QueryEscape(fieldName)
	encodedSearchText := url.QueryEscape(searchText)
	urlStr := fmt.Sprintf("%s/v1/o11y/autocomplete/attribute_values?aggregateOperator=noop&dataSource=logs&aggregateAttribute=&attributeKey=%s&searchText=%s&filterAttributeKeyDataType=string&tagType=resource", s.baseURL, encodedFieldName, encodedSearchText)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching logs field values", zap.String("fieldName", fieldName), zap.String("searchText", searchText))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", urlStr), zap.String("fieldName", fieldName), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", urlStr), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", urlStr), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved logs field values", zap.String("fieldName", fieldName), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetMetricsFieldValues(ctx context.Context, fieldName string, searchText string) (json.RawMessage, error) {
	encodedFieldName := url.QueryEscape(fieldName)
	encodedSearchText := url.QueryEscape(searchText)
	urlStr := fmt.Sprintf("%s/v1/o11y/fields/keys?signal=metrics&metricName=%s&searchText=%s&fieldContext=&fieldDataType=&source=", s.baseURL, encodedFieldName, encodedSearchText)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(ContentType, "application/json")
	req.Header.Set(EnvAPIKey, s.apiKey)

	ctx, cancel := context.WithTimeout(ctx, 600*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	s.logger.Debug("Fetching metrics field values", zap.String("fieldName", fieldName), zap.String("searchText", searchText))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.Error("HTTP request failed", zap.String("url", urlStr), zap.String("fieldName", fieldName), zap.Error(err))
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn("Failed to close response body", zap.Error(err))
		}
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("API request failed", zap.String("url", urlStr), zap.Int("status", resp.StatusCode), zap.String("response", string(body)))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed to read response body", zap.String("url", urlStr), zap.Error(err))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Debug("Successfully retrieved metrics field values", zap.String("fieldName", fieldName), zap.Int("status", resp.StatusCode))
	return body, nil
}

func (s *Client) GetTraceDetails(ctx context.Context, traceID string, includeSpans bool, startTime, endTime int64) (json.RawMessage, error) {
	if startTime == 0 || endTime == 0 {
		return nil, fmt.Errorf("start and end time parameters are required")
	}

	filterExpression := fmt.Sprintf("traceID = '%s'", traceID)
	limit := 1000

	queryPayload := types.BuildTracesQueryPayload(startTime, endTime, filterExpression, limit)
	queryJSON, err := json.Marshal(queryPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query payload: %w", err)
	}

	return s.QueryBuilderV5(ctx, queryJSON)
}

func (s *Client) GetTraceErrorAnalysis(ctx context.Context, startTime, endTime int64, serviceName string) (json.RawMessage, error) {
	filterExpression := "hasError = true"
	if serviceName != "" {
		filterExpression += fmt.Sprintf(" AND service.name in ['%s']", serviceName)
	}

	limit := 1000
	queryPayload := types.BuildTracesQueryPayload(startTime, endTime, filterExpression, limit)
	queryJSON, err := json.Marshal(queryPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query payload: %w", err)
	}

	return s.QueryBuilderV5(ctx, queryJSON)
}

func (s *Client) GetTraceSpanHierarchy(ctx context.Context, traceID string, startTime, endTime int64) (json.RawMessage, error) {
	if startTime == 0 || endTime == 0 {
		return nil, fmt.Errorf("start and end time parameters are required")
	}

	filterExpression := fmt.Sprintf("traceID = '%s'", traceID)
	limit := 1000
	queryPayload := types.BuildTracesQueryPayload(startTime, endTime, filterExpression, limit)
	queryJSON, err := json.Marshal(queryPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query payload: %w", err)
	}

	return s.QueryBuilderV5(ctx, queryJSON)
}

func (s *Client) CreateDashboard(ctx context.Context, dashboard types.Dashboard) (json.RawMessage, error) {
	url := fmt.Sprintf("%s/v1/o11y/dashboards", s.baseURL)

	dashboardJSON, err := json.Marshal(dashboard)
	if err != nil {
		return nil, fmt.Errorf("marshal dashboard: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(dashboardJSON))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set(EnvAPIKey, s.apiKey)
	req.Header.Set(ContentType, "application/json")

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req.WithContext(timeoutCtx))
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return body, nil
}

func (s *Client) UpdateDashboard(ctx context.Context, id string, dashboard types.Dashboard) error {
	url := fmt.Sprintf("%s/v1/o11y/dashboards/%s", s.baseURL, id)

	dashboardJSON, err := json.Marshal(dashboard)
	if err != nil {
		return fmt.Errorf("marshal dashboard: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(dashboardJSON))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set(EnvAPIKey, s.apiKey)
	req.Header.Set(ContentType, "application/json")

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Do(req.WithContext(timeoutCtx))
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
