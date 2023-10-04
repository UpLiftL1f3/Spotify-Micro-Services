package functions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func BuildUpdateQuery(tableName string, updateFields map[string]interface{}, conditionField string, conditionValue interface{}) (string, []interface{}) {
	// Construct the placeholders for the SET clause
	var setPlaceholders []string
	var setValues []interface{}
	index := 1
	for key, value := range updateFields {
		setPlaceholders = append(setPlaceholders, key+" = $"+strconv.Itoa(index))
		setValues = append(setValues, value)
		index++
	}

	// Construct the SET clause
	setClause := strings.Join(setPlaceholders, ", ")

	// Construct the query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", tableName, setClause, conditionField, index)
	setValues = append(setValues, conditionValue)

	return query, setValues
}

func BuildFindOneQuery(tableName string, conditionField string) string {
	// Construct the query
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, conditionField)

	return query
}

// BuildFindOneQueryWithToken constructs a query for finding one record based on multiple condition fields
func BuildFindOneQueryDynamic(tableName string, conditionFields map[string]interface{}) (string, []interface{}) {
	// Check if at least one condition field is provided
	if len(conditionFields) < 1 {
		panic("BuildFindOneQueryWithToken: At least one condition field is required")
	}

	// Construct the placeholders for the WHERE clause
	var wherePlaceholders []string
	var setValues []interface{}
	index := 1
	for key, value := range conditionFields {
		fmt.Println("KEY:", key)
		fmt.Println("value:", value)
		// Use different placeholders based on the type of the value
		switch v := value.(type) {
		case []string:
			// If it's an array, check if the array contains the specific token
			placeholder := fmt.Sprintf("%s = ANY($%d)", key, index)
			wherePlaceholders = append(wherePlaceholders, placeholder)
			setValues = append(setValues, pq.Array(v))
		default:
			placeholder := fmt.Sprintf("%s = $%d", key, index)
			wherePlaceholders = append(wherePlaceholders, placeholder)
			setValues = append(setValues, v)
		}

		index++
	}

	// Construct the WHERE clause
	whereClause := strings.Join(wherePlaceholders, " AND ")

	// Construct the query
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", tableName, whereClause)

	return query, setValues
}

func BuildDeleteQuery(tableName string, conditionField string) string {
	// Construct the query
	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, conditionField)

	return query
}

func BuildInsertQuery(tableName string, fields map[string]interface{}) (string, []interface{}) {
	var fieldNames []string
	var placeholders []string
	var values []interface{}

	index := 1
	for key, value := range fields {
		fieldNames = append(fieldNames, key)
		placeholders = append(placeholders, fmt.Sprintf("$%d", index))
		values = append(values, value)
		index++
	}

	fieldClause := strings.Join(fieldNames, ", ")
	placeholderClause := strings.Join(placeholders, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", tableName, fieldClause, placeholderClause)

	return query, values
}
