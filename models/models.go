package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
)

type Fields []struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserSchema struct {
	User struct {
		Fields Fields `json:"fields"`
	} `json:"User"`
}

func GenerateSchema(configPath string, schema interface{}) error {
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

func GeneratePostgreSQLTable(tableName string, structInterface interface{}) string {
	var columns []string
	switch tableName {
	case "users", "Users":
		columns = append(columns, "id INTEGER PRIMARY KEY AUTOINCREMENT")
	}
	typeOf := reflect.TypeOf(structInterface)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		column := field.Name
		fieldType := ""
		switch column {
		case "Password", "password":
			fieldType = "VARCHAR(255) NOT NULL"
		case "Email", "email":
			fieldType = "VARCHAR(255) NOT NULL UNIQUE"
		case "Username", "username":
			fieldType = "VARCHAR(255) NOT NULL UNIQUE"
		default:
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
		}
		columns = append(columns, fmt.Sprintf("%s %s", column, fieldType))
	}
	fmt.Printf("query: CREATE TABLE IF NOT EXISTS \"%s\"(%s);\n", tableName, strings.Join(columns, ","))
	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS \"%s\"(%s);", tableName, strings.Join(columns, ","))
}

func InsertDataFromStruct(tableName string, structInterface interface{}) string {
	var columns []string
	var values []string
	typeOf := reflect.TypeOf(structInterface)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		column := field.Name
		columns = append(columns, column)
		value := reflect.ValueOf(structInterface).FieldByName(column)
		values = append(values, fmt.Sprintf("'%v'", value))
	}
	fmt.Printf("query: INSERT INTO \"%s\"(%s) VALUES( %s );\n", tableName, strings.Join(columns, ","), strings.Join(values, ","))
	return fmt.Sprintf("INSERT INTO \"%s\"(%s) VALUES( %s );", tableName, strings.Join(columns, ","), strings.Join(values, ","))
}

func QueryAdminUserDB(ut string, userStruct interface{}) string {
	var columns []string
	var values []string
	var combinedStuff string
	typeOf := reflect.TypeOf(userStruct)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		column := field.Name
		switch column {
		case regexp.MustCompile(`(?i)password`).FindString(column):
			continue
		}
		columns = append(columns, column)
		value := reflect.ValueOf(userStruct).FieldByName(column)
		values = append(values, fmt.Sprintf("'%v'", value))
	}
	for i := 0; i < len(columns); i++ {
		combinedStuff += fmt.Sprintf("%s = %s", columns[i], values[i])
		if i != len(columns)-1 {
			combinedStuff += " AND "
		}
	}
	fmt.Printf("query: SELECT * FROM %s WHERE %s;\n ", ut, combinedStuff)

	return fmt.Sprintf("SELECT * FROM %s WHERE %s;\n ", ut, combinedStuff)
}

func IsUserInitialized(db *sql.DB) bool {
	result, err := db.Exec("SELECT * FROM users;")
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
