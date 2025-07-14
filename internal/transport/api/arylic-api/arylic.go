package arylic_api

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// https://developer.arylic.com/httpapi/#http-api

type ArylicAPI struct {
	IP         string
	httpClient *http.Client
	//Network     *NetworkApi
	//DeviceInfo  *DeviceInfoApi
	PlayBack *PlayBackApi
	//USBPlayback *USBPlaybackApi
	//MultiRoom   *MultiRoomApi
}

func NewAPI(IP string) *ArylicAPI {
	client := newHTTPClient(5 * time.Second) // общий таймаут ответа
	api := &ArylicAPI{
		IP:         IP,
		httpClient: client,
	}
	api.PlayBack = NewPlayBackApi(api)
	return api
}

func newHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   3 * time.Second, // таймаут установки соединения
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 3 * time.Second,
		},
	}
}

// DoAPIRequest выполняет HTTP-запрос к API устройства и заполняет переданную структуру v результатами запроса.
// Если v == nil, то просто проверяет, что ответ "OK".
func (a *ArylicAPI) DoAPIRequest(method, command string, v interface{}) error {
	url := fmt.Sprintf("http://%s/httpapi.asp?command=%s", a.IP, command)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	resp, err := a.doRequestWithRetry(a.httpClient, req, 3, 500*time.Millisecond)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	s := string(body)
	if s == "OK" {
		return nil
	}
	if v != nil {
		if err := json.Unmarshal(body, v); err == nil {
			return nil
		}
	}
	return fmt.Errorf("unexpected response: %s", s)
}

// doRequestWithRetry выполняет HTTP-запрос с повторными попытками в случае неудачи.
func (a *ArylicAPI) doRequestWithRetry(client *http.Client, req *http.Request, retries int, delay time.Duration) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i < retries; i++ {
		resp, err = client.Do(req)
		if err == nil {
			return resp, nil
		}
		time.Sleep(delay)
	}
	return nil, fmt.Errorf("request failed after %d retries: %w", retries, err)
}
