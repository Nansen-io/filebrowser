package chainfs

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetLoginUrl fetches the Azure AD B2C login URL from ChainFS API
func GetLoginUrl(baseUrl string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/NansenFile/LoginURL", baseUrl)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to fetch login URL from ChainFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ChainFS API returned status %d when fetching login URL", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login URL response: %w", err)
	}

	return string(body), nil
}

// GetLogoutUrl fetches the logout URL from ChainFS API
func GetLogoutUrl(baseUrl string) (string, error) {
	endpoint := fmt.Sprintf("%s/api/NansenFile/LogoutURL", baseUrl)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to fetch logout URL from ChainFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ChainFS API returned status %d when fetching logout URL", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read logout URL response: %w", err)
	}

	return string(body), nil
}
