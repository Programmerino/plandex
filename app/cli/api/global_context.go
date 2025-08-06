package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	shared "plandex-shared"
)

type GlobalContextRequest struct {
	Content string `json:"content"`
}

func (a *Api) GetGlobalContext() (string, *shared.ApiError) {
	serverUrl := fmt.Sprintf("%s/global_context", GetApiHost())
	
	resp, err := authenticatedFastClient.Get(serverUrl)
	if err != nil {
		return "", &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error sending request: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil // No global context set
	}

	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		apiErr := HandleApiError(resp, errorBody)
		return "", apiErr
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error reading response: %v", err)}
	}

	var respData map[string]string
	err = json.Unmarshal(body, &respData)
	if err != nil {
		return "", &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error unmarshalling response: %v", err)}
	}

	return respData["content"], nil
}

func (a *Api) UpdateGlobalContext(content string) *shared.ApiError {
	serverUrl := fmt.Sprintf("%s/global_context", GetApiHost())
	
	req := GlobalContextRequest{Content: content}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error marshalling request: %v", err)}
	}

	resp, err := authenticatedFastClient.Post(serverUrl, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error sending request: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		apiErr := HandleApiError(resp, errorBody)
		return apiErr
	}

	return nil
}

func (a *Api) DeleteGlobalContext() *shared.ApiError {
	serverUrl := fmt.Sprintf("%s/global_context", GetApiHost())
	
	req, err := http.NewRequest("DELETE", serverUrl, nil)
	if err != nil {
		return &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error creating request: %v", err)}
	}

	resp, err := authenticatedFastClient.Do(req)
	if err != nil {
		return &shared.ApiError{Type: shared.ApiErrorTypeOther, Msg: fmt.Sprintf("error sending request: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		errorBody, _ := io.ReadAll(resp.Body)
		apiErr := HandleApiError(resp, errorBody)
		return apiErr
	}

	return nil
}