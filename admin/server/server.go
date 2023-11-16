package server

import (
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"encoding/base64"

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

func AuthMiddleware(c *fiber.Ctx, st *session.Store) error {
	sess, err := st.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if sess.Get("email") == nil {
		return c.Redirect("/admin/entry/login")
	}

	return c.Next()
}

func handleUserLogin(c *fiber.Ctx, db *sql.DB, store *session.Store) error {
	var formData map[string]interface{}
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
	var invalidKeys []string
	var fieldErrors []string
	for key, value := range formData {
		switch value.(string) {
		case "":
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "empty field(s)",
			})
		}

		if regexp.MustCompile(`(?i)email`).MatchString(key) {
			if !regexp.MustCompile(`(?i)^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`).MatchString(value.(string)) {
				invalidKeys = append(invalidKeys, key)
				fieldErrors = append(fieldErrors, fmt.Sprintf("Invalid %s", key))
				continue
			}
			continue
		} else if regexp.MustCompile(`(?i)password`).MatchString(key) {
			if len(value.(string)) < 1 {
				invalidKeys = append(invalidKeys, key)
				fieldErrors = append(fieldErrors, fmt.Sprintf("%s too short", key))
				continue
			}
			continue
		} else {
			if !regexp.MustCompile(`(?i)^[\w]+$`).MatchString(value.(string)) {
				invalidKeys = append(invalidKeys, key)
				fieldErrors = append(fieldErrors, fmt.Sprintf("Invalid characters in %s", key))
				continue
			}
			continue
		}
	}

	if len(invalidKeys) < 1 || invalidKeys == nil {
		fmt.Println("no errors")
	} else if len(invalidKeys) == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fieldErrors,
		})
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": fieldErrors,
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
	for key, value := range userDeets {
		switch key {
		case regexp.MustCompile(`(?i)password`).FindString(key):
			fmt.Println(value.(string))
			fmt.Println(formData[key].(string))
			if utils.CheckPassword(formData[key].(string), value.(string)) != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid password",
				})
			} else if utils.CheckPassword(formData[key].(string), value.(string)) == nil {
				continue
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid credentials",
				})
			}
		case regexp.MustCompile(`(?i)id`).FindString(key):
			continue
		default:
			if value != formData[key] {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "invalid credentials",
				})
			}
		}
	}
	for key, value := range userDeets {
		switch key {
		case regexp.MustCompile(`(?i)email`).FindString(key):
			valueThing := base64.StdEncoding.EncodeToString([]byte(value.(string)))
			c.Cookie(&fiber.Cookie{
				Name:        "colibase",
				Value:       valueThing,
				Expires:     time.Now().Add(24 * 7 * time.Hour),
				SameSite:    fiber.CookieSameSiteStrictMode,
				MaxAge:      24 * 7 * 60 * 60,
				HTTPOnly:    true,
				SessionOnly: true,
			})
			sess, err := store.Get(c)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			fmt.Println(sess.Get("email"))
			keys := sess.Keys()
			fmt.Println(keys)
			sess.Set("email", valueThing)
			sess.SetExpiry(24 * 7 * time.Hour)
			err = sess.Save()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

		default:
			continue
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
			fmt.Println(schema[key])
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

		}
	}
	fmt.Println(formData)
	if len(formErrors) > 0 || formData["failure"] == "true" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": formErrors,
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
	var jsonSchema []map[string]interface{}
	for key, value := range schema {
		jsonSchema = append(jsonSchema, map[string]interface{}{
			"name":   key,
			"values": value,
		})
	}
	fmt.Println(jsonSchema)
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
		return c.JSON(jsonSchema)
	}
}
