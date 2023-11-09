package pkg

import (
	"fmt"
	"io"

	"github.com/diogenxs/dxs/utils"
)

func MyPublicIP(verbose bool) string {
	urlList := []string{
		"http://icanhazip.com",
		"https://ifconfig.me",
		"https://api.ipify.org",
		"http://checkip.amazonaws.com",
	}
	result := make(chan string, 1)
	for _, url := range urlList {
		go makeHttpRequest(url, result, verbose)
	}

	return <-result
}

func makeHttpRequest(url string, c chan string, verbose bool) {
	// resp, err := http.Get(url)
	resp, err := utils.MyHTTP(url, "GET", verbose)
	if err != nil {
		fmt.Errorf("Error getting public IP: %s", err)
	}
	defer resp.Body.Close()

	result, _ := io.ReadAll(resp.Body)

	c <- string(result)
}
