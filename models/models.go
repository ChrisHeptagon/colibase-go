package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type UserSchema struct {
	User struct {
		Fields []struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"fields"`
	} `json:"User"`
}

type SchemaInterface interface{}

func GenerateSchema(configPath string, schema SchemaInterface) error {
	configJSON, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Unmarshal JSON into the provided schema
	if err := json.Unmarshal(configJSON, schema); err != nil {
		return err
	}

	return nil
}

func MapToStruct(mapping map[string]interface{}) (interface{}, error) {
	var fields []reflect.StructField
	for key, value := range mapping {
		fieldType := reflect.TypeOf(value)
		fields = append(fields, reflect.StructField{
			Name: key,
			Type: fieldType,
		})
	}
	structType := reflect.StructOf(fields)
	structValue := reflect.New(structType).Elem()
	for key, value := range mapping {
		field := structValue.FieldByName(key)
		if field.IsValid() {
			field.Set(reflect.ValueOf(value))
		}
	}
	return structValue.Interface(), nil
}

func GenerateSQLTable(tableName string, structInterface interface{}) string {
	var columns []string
	typeOf := reflect.TypeOf(structInterface)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		column := field.Tag.Get("db")

		if column == "" {
			column = strings.ToLower(field.Name)
		}
		fieldType := ""
		switch field.Type.Kind() {
		case reflect.String:
			fieldType = "TEXT"
		case reflect.Int:
			fieldType = "INT"
		case reflect.Bool:
			fieldType = "BOOLEAN"
		default:
			fieldType = "TEXT"
		}
		columns = append(columns, fmt.Sprintf("%s %s", column, fieldType))
	}
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(columns, ", "))
}

func IsUserInitialized(db *sql.DB, ut string) bool {
	result, err := db.Exec("SELECT * FROM %s", ut)
	if err != nil {
		return false
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false
	}
	if rowsAffected > 0 {
		return true
	}
	if rowsAffected == 0 {
		return false
	}
	return false
}
