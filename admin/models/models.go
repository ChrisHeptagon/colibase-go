package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DefaultUserSchema struct {
	Email    string `form_type:"email" required:"true" pattern:"^[^\\s]+@[^\\s]+\\.\\w+$" order:"1"`
	Username string `form_type:"text" required:"true" pattern:"^[a-zA-Z0-9]+$" order:"2"`
	Password string `form_type:"password" required:"true" pattern:"^[^\\s]+$" order:"3"`
}

var SchemaFields = []string{
	"form_type",
	"required",
	"pattern",
	"order",
}

func DeleteCookie(db *sql.DB, tableName string, cookieValue string) error {
	result, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE cookie = ?;", tableName))
	if err != nil {
		return err
	}
	defer result.Close()
	_, err = result.Exec(cookieValue)
	if err != nil {
		return err
	}
	return nil
}

func CheckCookie(db *sql.DB, tableName string, cookieValue string) error {
	makeResult, err := db.Prepare(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, cookie VARCHAR(255) UNIQUE NOT NULL);", tableName))
	if err != nil {
		return err
	}
	defer makeResult.Close()
	_, err = makeResult.Exec()
	if err != nil {
		return err
	}
	result, err := db.Prepare(fmt.Sprintf("SELECT cookie FROM %s WHERE cookie = ? LIMIT 1;", tableName))
	if err != nil {
		return err
	}
	defer result.Close()
	rows, err := result.Query(cookieValue)
	if err != nil {
		return err
	}
	if rows.Err() == sql.ErrNoRows {
		return fmt.Errorf("invalid cookie")
	}
	defer rows.Close()
	var cookie string
	for rows.Next() {
		err = rows.Scan(&cookie)
		if err != nil {
			return err
		}
	}
	fmt.Println("Cookie: ", cookie)
	fmt.Println("Cookie Value: ", cookieValue)
	if cookie == "" {
		result, err := db.Prepare(fmt.Sprintf("INSERT INTO %s(cookie) VALUES(?);", tableName))
		if err != nil {
			return err
		}
		defer result.Close()
		_, err = result.Exec(cookieValue)
		if err != nil {
			return err
		}
	}
	return nil
}

func GenAdminSchema(db *sql.DB, tableName string) (map[string]map[string]string, error) {
	var schema = make(map[string]map[string]string)
	makeQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, schema BLOB UNIQUE);", tableName)
	db.Exec(makeQuery)
	firstQuery := fmt.Sprintf("SELECT schema FROM %s;", tableName)
	fmt.Println(firstQuery)
	rows, err := db.Query(firstQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemaByte []byte
	for rows.Next() {
		for i := 0; i < len(SchemaFields); i++ {
			err = rows.Scan(&schemaByte)
			if err != nil {
				return nil, err
			}
		}

	}

	if len(schemaByte) > 1 || schemaByte != nil {
		err = json.Unmarshal(schemaByte, &schema)
		if err != nil {
			return nil, err
		}
	} else {
		defaultSchema := reflect.TypeOf(&DefaultUserSchema{}).Elem()
		for i := 0; i < defaultSchema.NumField(); i++ {
			field := defaultSchema.Field(i)
			schema[field.Name] = make(map[string]string)
			for _, schemaField := range SchemaFields {
				schema[field.Name][schemaField] = field.Tag.Get(schemaField)
			}
		}
	}
	return schema, nil
}

func InitDB() (*sql.DB, error) {
	// db, err := sql.Open("sqlite3", fmt.Sprintf("file:./db/%s.sqlite?cache=shared&mode=rwc&_journal_mode=WAL&_synchronous=NORMAL&_foreign_keys=ON", os.Getenv("DB_NAME")))
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:./db/%s.sqlite?cache=shared&mode=rwc&_synchronous=NORMAL&_foreign_keys=ON", os.Getenv("DB_NAME")))
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}
	db.Exec(fmt.Sprintf("CREATE DATABASE %s;", os.Getenv("DB_NAME")))
	return db, nil
}

func MapToStruct[T any](mapping map[string]T) (interface{}, error) {
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

func GenerateAdminTable(db *sql.DB, tableName string, structInterface interface{}) error {
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
			fieldType = fmt.Sprintf("VARCHAR(255) NOT NULL UNIQUE CHECK(%s LIKE '%%@%%.%%')", column)
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
	formattedColumns := strings.Join(columns, ", ")
	query := fmt.Sprintf(" CREATE TABLE IF NOT EXISTS \"%s\" (%s);", tableName, formattedColumns)
	fmt.Println(query)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func InsertDataFromStruct(db *sql.DB, tableName string, structInterface interface{}) error {
	var columns []string
	var values []string
	var valuesInterface []interface{}
	typeOf := reflect.TypeOf(structInterface)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		column := field.Name
		columns = append(columns, column)
		value := reflect.ValueOf(structInterface).FieldByName(column).Interface()
		values = append(values, value.(string))
	}
	for _, value := range values {
		valuesInterface = append(valuesInterface, value)
	}
	blankedValues := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		blankedValues[i] = "?"
	}
	fmt.Printf("INSERT INTO \"%s\"(%s) VALUES( %s );", tableName, strings.Join(columns, ", "), strings.Join(blankedValues, ","))
	query := fmt.Sprintf("INSERT INTO \"%s\"(%s) VALUES( %s );", tableName, strings.Join(columns, ", "), strings.Join(blankedValues, ","))
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	fmt.Println("VALUES: ", valuesInterface)
	_, err = stmt.Exec(valuesInterface...)
	if err != nil {
		return err
	}
	return nil
}

func QueryAdminUserDB(db *sql.DB, ut string, userStruct interface{}) (*sql.Rows, error) {
	var columns []string
	var combinedStuff string
	var values []interface{}
	typeOf := reflect.TypeOf(userStruct)
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		column := field.Name
		switch column {
		case regexp.MustCompile(`(?i)password`).FindString(column):
			continue
		}
		columns = append(columns, column)
		values = append(values, reflect.ValueOf(userStruct).FieldByName(column).Interface())
	}
	for i := 0; i < len(columns); i++ {
		combinedStuff += fmt.Sprintf("%s = ?", columns[i])
		if i != len(columns)-1 {
			combinedStuff += " AND "
		}
	}
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s;\n ", ut, combinedStuff)
	fmt.Print(query)
	stmt, err := db.Prepare(query)
	if err != nil {
		return &sql.Rows{}, err
	}
	defer stmt.Close()
	result, err := stmt.Query(values...)
	if err != nil {
		return &sql.Rows{}, err
	}
	return result, nil
}

func IsUserInitialized(db *sql.DB) bool {
	rows, err := db.Query("SELECT id FROM users LIMIT 1;")
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		switch err := rows.Err(); err {
		case sql.ErrNoRows:
			return false
		case nil:
			return true
		default:
			return true
		}
	}
	return false
}
