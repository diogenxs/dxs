package models

import (
	"database/sql"
	"encoding/json"
)

type Alert struct {
	Fingerprint string            `json:"fingerprint"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Status      struct {
		State string `json:"state"`
	} `json:"status"`
	GeneratorURL string `json:"generatorURL"`
	Acknowledged bool   `json:"acknowledged"`
}

// InsertAlert inserts a new Alert into the SQLite database.
func InsertAlert(db *sql.DB, alert Alert) error {
	// Convert the labels and annotations maps to JSON for storage
	labelsJSON, err := json.Marshal(alert.Labels)
	if err != nil {
		return err
	}
	annotationsJSON, err := json.Marshal(alert.Annotations)
	if err != nil {
		return err
	}

	// Prepare SQL statement for inserting data
	stmt, err := db.Prepare("INSERT INTO alerts(fingerprint, labels, annotations, state, generatorURL) values(?,?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(alert.Fingerprint, labelsJSON, annotationsJSON, alert.Status.State, alert.GeneratorURL)
	return err
}

// UpsertAlert inserts a new Alert into the SQLite database or updates it if it already exists.
func UpsertAlert(db *sql.DB, alert Alert) error {
	// Convert the labels and annotations maps to JSON for storage
	labelsJSON, err := json.Marshal(alert.Labels)
	if err != nil {
		return err
	}
	annotationsJSON, err := json.Marshal(alert.Annotations)
	if err != nil {
		return err
	}

	// Prepare SQL statement for inserting or replacing data
	stmt, err := db.Prepare("INSERT INTO alerts(fingerprint, labels, annotations, state, generatorURL) values(?,?,?,?,?) ON CONFLICT (fingerprint) DO UPDATE SET labels=excluded.labels , annotations=excluded.annotations, state=excluded.state, generatorURL=excluded.generatorURL")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute SQL statement
	_, err = stmt.Exec(alert.Fingerprint, labelsJSON, annotationsJSON, alert.Status.State, alert.GeneratorURL)
	return err
}

func ListAlerts(db *sql.DB) ([]Alert, error) {
	const query = "SELECT fingerprint, labels, annotations, state, generatorURL FROM alerts WHERE acknowledged = FALSE"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []Alert
	for rows.Next() {
		var (
			fingerprint     string
			labelsJSON      string
			annotationsJSON string
			state           string
			generatorURL    string
		)

		err := rows.Scan(&fingerprint, &labelsJSON, &annotationsJSON, &state, &generatorURL)
		if err != nil {
			return nil, err
		}

		var labels map[string]string
		var annotations map[string]string
		err = json.Unmarshal([]byte(labelsJSON), &labels)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(annotationsJSON), &annotations)
		if err != nil {
			return nil, err
		}

		alerts = append(alerts, Alert{
			Fingerprint: fingerprint,
			Labels:      labels,
			Annotations: annotations,
			Status: struct {
				State string `json:"state"`
			}{State: state},
			GeneratorURL: generatorURL,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alerts, nil
}

func AckAlert(db *sql.DB, fingerprint string) error {
	const query = "UPDATE alerts SET acknowledged = TRUE WHERE fingerprint = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(fingerprint)
	return err
}
