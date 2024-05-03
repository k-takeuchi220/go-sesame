package main

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aead/cmac"
)

const (
	ApiBaseURL = "https://api.candyhouse.co/v3/sesame/"
)

type Command int

const (
	// CommandLock   Command = 82
	// CommandUnlock Command = 83
	CommandToggle Command = 88
)

func sendCommand(apiKey, deviceUUID, secretKey string, command Command) error {
	headers := map[string]string{
		"x-api-key": apiKey,
	}

	history := fmt.Sprintf("call api. command:%d", command)
	encodedHistory := base64.StdEncoding.EncodeToString([]byte(history))

	sign, err := generateCMAC(secretKey)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://app.candyhouse.co/api/sesame2/%s/cmd", deviceUUID)
	body := map[string]interface{}{
		"cmd":     command,
		"history": encodedHistory,
		"sign":    hex.EncodeToString(sign),
	}
	jsonBody, _ := json.Marshal(body)

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	return nil
}

func generateCMAC(secretKey string) ([]byte, error) {
	key, err := hex.DecodeString(secretKey)
	if err != nil {
		panic(err)
	}

	date := uint32(time.Now().Unix())
	dateBytes := make([]byte, 4)
	dateBytes[0] = byte(date & 0xff)
	dateBytes[1] = byte((date >> 8) & 0xff)
	dateBytes[2] = byte((date >> 16) & 0xff)
	dateBytes[3] = byte((date >> 24) & 0xff)

	message := dateBytes[1:4]

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cmac, err := cmac.Sum(message, block, block.BlockSize())
	if err != nil {
		panic(err)
	}

	return cmac, nil
}
