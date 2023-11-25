package jsonHelper

import(
	"net/http"
	"encoding/json"
	"time"
)

var client *http.Client

func GetJson(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}