package server

import (
	"database/sql"
	"fmt"
	"regexp"
	"sort"
	"time"

	"encoding/json"

	"github.com/ChrisHeptagon/colibase/admin/models"
	"github.com/ChrisHeptagon/colibase/admin/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3/v2"
)

func MainServer(db *sql.DB) {
	app := fiber.New()
	storageDB := sqlite3.New(sqlite3.Config{
		Database: "./db/sessions.db",
		Table:    "sessions",
	})
	store := session.New(
		session.Config{
			Expiration: 24 * 7 * time.Hour,
			Storage:    storageDB,
			KeyLookup:  "cookie:colibase",
		})
	app.Use(compress.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Post("/api/login", func(c *fiber.Ctx) error {
		return handleUserLogin(c, db, store)
	})
	app.Get("/api/login-schema", func(c *fiber.Ctx) error {
		return loginSchema(db, c)
	})
	app.Get("/api/user-initialization-status", func(c *fiber.Ctx) error {
		return handleUserInitializatonStatus(c, db, "users")
	})
	app.Post("/api/init-login", func(c *fiber.Ctx) error {
		return handleUserInitializaton(c, db)
	})

	app.Listen(":6700")
}

func handleUserLogin(c *fiber.Ctx, db *sql.DB, store *session.Store) error {
	var formData map[string]string
	if err := c.BodyParser(&formData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(formData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty field(s)",
		})
	}
	var formErrors []string
	schema, err := models.GenAdminSchema(db, "admin_schema")
	fmt.Println("Login Schema: ", schema["Password"]["pattern"])
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	for key, value := range formData {
		switch value {
		case "":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "empty field(s)",
			})
		}
		if schema[key] == nil {
			defer delete(formData, key)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid Field: %s", key),
			})
		}
		if schema[key] != nil {
			if schema[key]["required"] == "true" {
				if value == "" {
					formErrors = append(formErrors, fmt.Sprintf("Empty Field: %s", key))
				}
			}
			if schema[key]["pattern"] != "" {
				if !regexp.MustCompile(schema[key]["pattern"]).MatchString(value) {
					formErrors = append(formErrors, fmt.Sprintf("Invalid %s", key))
				}
			}

		}

	}

	if len(formErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": formErrors,
		})
	}

	structFormData, err := models.MapToStruct(formData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	rows, err := models.QueryAdminUserDB(db, "users", structFormData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userDeets := make(map[string]interface{})

	defer rows.Close()
	column, err := rows.Columns()
	columnsInterface := make([]interface{}, len(column))
	for rows.Next() {
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		for i := range column {
			columnsInterface[i] = &columnsInterface[i]
		}
		err = rows.Scan(columnsInterface...)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		for i := range column {
			userDeets[column[i]] = columnsInterface[i]
		}
	}
	if len(userDeets) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}
	fmt.Println(userDeets)
	fmt.Println(formData)
	for key, value := range userDeets {
		if regexp.MustCompile(`(?i)password`).FindString(key) != "" {
			if utils.CheckPassword(formData[key], value.(string)) == nil {
				continue
			} else if utils.CheckPassword(formData[key], value.(string)) != nil {
				fmt.Println("Password:", formData[key])
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid credentials",
				})
			}
		} else if regexp.MustCompile(`(?i)id`).FindString(key) != "" {
			continue
		} else {
			if value != formData[key] {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid credentials",
				})
			}
		}
	}
	return c.JSON(
		fiber.Map{
			"message": "Login Successful",
			"status":  fiber.StatusOK,
		},
	)
}

func handleUserInitializaton(c *fiber.Ctx, db *sql.DB) error {
	var formData map[string]string
	if err := c.BodyParser(&formData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if len(formData) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "empty field(s)",
		})
	}
	var formErrors []string
	schema, err := models.GenAdminSchema(db, "admin_schema")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	for key, value := range formData {
		if schema[key] == nil {
			defer delete(formData, key)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid Field: %s", key),
			})
		}
		if schema[key] != nil {
			if schema[key]["required"] == "true" {
				if value == "" {
					formErrors = append(formErrors, fmt.Sprintf("Empty Field: %s", key))
				}
			}
			if schema[key]["pattern"] != "" {
				if !regexp.MustCompile(schema[key]["pattern"]).MatchString(value) {
					formData = map[string]string{
						"failure": "true",
					}
					formErrors = append(formErrors, fmt.Sprintf("Invalid %s\n", key))
				}
			}
			if regexp.MustCompile(`(?i)password`).FindString(key) != "" {
				hashedPass, err := utils.HashPassword(value)
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"error": err.Error(),
					})
				}
				formData[key] = hashedPass
			}
		}
	}
	if len(formErrors) > 0 || formData["failure"] == "true" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": formErrors,
		})
	}
	structFormData, err := models.MapToStruct(formData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = models.GenerateAdminTable(db, "users", structFormData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = models.InsertDataFromStruct(db, "users", structFormData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(
		fiber.Map{
			"message": "User Initialized",
			"status":  fiber.StatusOK,
		},
	)
}

func handleUserInitializatonStatus(c *fiber.Ctx, db *sql.DB, tn string) error {
	if models.IsUserInitialized(db) {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func loginSchema(db *sql.DB, c *fiber.Ctx) error {
	schema, err := models.GenAdminSchema(db, "admin_schema")
	var jsonSchema []map[string]string
	for key, value := range schema {
		modVal := make(map[string]string)
		modVal["name"] = key
		for k, v := range value {
			modVal[k] = v
		}
		jsonSchema = append(jsonSchema, modVal)
	}
	sort.Slice(jsonSchema, func(i, j int) bool {
		return jsonSchema[i]["name"] < jsonSchema[j]["name"]
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if schema == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "schema is nil",
		})
	} else if len(schema) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "schema is empty",
		})
	} else {
		jsonTest, err := json.Marshal(jsonSchema)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		fmt.Println(string(jsonTest))
		return c.JSON(jsonSchema)
	}
}
