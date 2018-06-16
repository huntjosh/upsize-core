package models

import (
	"strings"
	"strconv"
	"database/sql"
)

type QueryBuilder struct {
	query     string
	hasWhereClause bool
	paramCount  int
	params []interface{}
}

func (q QueryBuilder) AddQueryString(query string, hasWhere bool) QueryBuilder {
	q.query = query
	q.paramCount = 0
	q.hasWhereClause = hasWhere

	return q
}

func (q QueryBuilder) AddParams(params []interface{}) QueryBuilder {
	for _, val := range params {
		q.params = append(q.params, val)
		q.paramCount++
	}
	return q
}

func (q QueryBuilder) AddWhereClause(whereClause string, params []interface{}) QueryBuilder {
	if q.hasWhereClause {
		q.query += " AND " + whereClause
	} else {
		q.query += " WHERE " + whereClause
	}

	for _, val := range params {
		q.params = append(q.params, val)
		q.paramCount++
	}

	return q
}

func (q QueryBuilder) Get(db *sql.DB) (*sql.Rows, error) {
	queryParts := strings.Split(q.query, "$1")
	queryString := ""
	for i := 1; i <= q.paramCount; i++ {
		queryString += queryParts[i - 1] + "$" + strconv.Itoa(i)
	}

	return db.Query(queryString, q.params...)
}