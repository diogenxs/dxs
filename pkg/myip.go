package pkg

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func MyPublicIP() string {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		fmt.Errorf("Error getting public IP: %s", err)
		return ""
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
