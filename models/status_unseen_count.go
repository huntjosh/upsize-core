package models

import "database/sql"

type StatusUnseenCount struct {
	Status string
	Count  int
}

func GetContractorUnseenCounts(db *sql.DB, contractorId string) ([]StatusUnseenCount, error) {
	rows, err := db.Query(
		"SELECT status, count(*) AS jobs FROM contractor_jobs WHERE contractor_id = $1 GROUP BY status", contractorId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	counts := make([]StatusUnseenCount, 0)
	for rows.Next() {
		var count StatusUnseenCount
		if err := rows.Scan(&count.Status, &count.Count); err != nil {
			return nil, err
		}
		counts = append(counts, count)
	}

	return counts, nil
}
