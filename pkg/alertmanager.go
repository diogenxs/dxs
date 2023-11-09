package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/diogenxs/dxs/models"
	"github.com/spf13/viper"
)

// Alert is the structure that represents an alert from Alertmanager

func SyncAlerts() error {
	db, err := models.GetDB()
	if err != nil {
		return err
	}

	alerts, err := GetFilteredAlerts()
	if err != nil {
		return err
	}
	for _, alert := range alerts {
		fmt.Println(alert.Fingerprint + ": " + alert.Labels["alertname"])
		err = models.UpsertAlert(db, alert)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetFilteredAlerts() ([]models.Alert, error) {
	m := viper.GetStringMap("alertmanager")
	if len(m) == 0 {
		return nil, fmt.Errorf("No keys added to config file: %s", viper.ConfigFileUsed())
	}
	var alerts []models.Alert
	for _, v := range m {

		config := v.(map[string]interface{})

		url, urlOk := config["url"].(string)
		if !urlOk {
			return nil, fmt.Errorf("error getting url from config file: %s", viper.ConfigFileUsed())
		}
		a, err := getAlerts(url)
		if err != nil {
			return nil, err
		}

		filters := config["filters"].(map[string]interface{})

		filteredAlerts := filterAlerts(a, filters)
		alerts = append(alerts, filteredAlerts...)
	}
	return alerts, nil

}

// getAlerts queries the Alertmanager API and returns a slice of Alert structs
func getAlerts(apiURL string) ([]models.Alert, error) {
	alertManagerURI := "/api/v2/alerts"

	resp, err := http.Get(apiURL + alertManagerURI)
	if err != nil {
		return nil, fmt.Errorf("error making the GET request to Alertmanager API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var alerts []models.Alert
	err = json.Unmarshal(body, &alerts)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return alerts, nil
}

func filterAlerts(alerts []models.Alert, filter map[string]interface{}) []models.Alert {
	var filtered []models.Alert
	for _, alert := range alerts {
		for k, v := range filter {
			value := v.(string)
			if alert.Labels[k] == value {
				filtered = append(filtered, alert)
			}
		}
	}
	return filtered
}
