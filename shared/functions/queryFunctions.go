package functions

import (
	"fmt"
	"strconv"
	"strings"
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
