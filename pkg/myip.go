package pkg

import (
	"fmt"
	"io/ioutil"

	"github.com/diogenxs/dxs/utils"
)

func MyPublicIP(verbose bool) string {
	result := make(chan string, 1)
	go makeHttpRequest("http://icanhazip.com", result, verbose)
	go makeHttpRequest("https://ifconfig.me", result, verbose)
	go makeHttpRequest("https://api.ipify.org", result, verbose)
	go makeHttpRequest("http://checkip.amazonaws.com", result, verbose)

	return <-result
}

func makeHttpRequest(url string, c chan string, verbose bool) {
	// resp, err := http.Get(url)
	resp, err := utils.MyHTTP(url, "GET", verbose)
	if err != nil {
		fmt.Errorf("Error getting public IP: %s", err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	c <- string(result)
}
