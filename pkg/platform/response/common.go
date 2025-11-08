package response

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// CommonResponse platform server common response format
type CommonResponse struct {
	Code   int    `json:"code"`
	Data   any    `json:"data"`
	Status string `json:"status"`
}

// ParseResponse parse platform response and return data part
func ParseResponse(respBody []byte, target any) error {
	var commonResp CommonResponse

	if err := json.Unmarshal(respBody, &commonResp); err != nil {
		log.Error().Err(err).Msg("Failed to parse platform response")
		return fmt.Errorf("Failed to parse platform response: %w", err)
	}

	// Check response status code
	if commonResp.Code != 0 {
		log.Error().
			Int("code", commonResp.Code).
			Str("status", commonResp.Status).
			Msg("Platform returned error")
		return fmt.Errorf("Platform returned error: code=%d, status=%s",
			commonResp.Code, commonResp.Status)
	}

	// Check status
	if commonResp.Status != "OK" {
		log.Error().
			Int("code", commonResp.Code).
			Str("status", commonResp.Status).
			Msg("Platform status abnormal")
		return fmt.Errorf("Platform status abnormal: code=%d, status=%s",
			commonResp.Code, commonResp.Status)
	}

	// Parse data part
	if commonResp.Data != nil {
		dataBytes, err := json.Marshal(commonResp.Data)
		if err != nil {
			log.Error().Err(err).Msg("Failed to serialize response data")
			return fmt.Errorf("Failed to serialize response data: %w", err)
		}

		if err := json.Unmarshal(dataBytes, target); err != nil {
			log.Error().Err(err).Msg("Failed to parse response data")
			return fmt.Errorf("Failed to parse response data: %w", err)
		}
	}

	return nil
}

// ParseRawResponse parse platform response and return raw CommonResponse
func ParseRawResponse(respBody []byte) (*CommonResponse, error) {
	var commonResp CommonResponse

	if err := json.Unmarshal(respBody, &commonResp); err != nil {
		log.Error().Err(err).Msg("Failed to parse platform response")
		return nil, fmt.Errorf("Failed to parse platform response: %w", err)
	}

	return &commonResp, nil
}
