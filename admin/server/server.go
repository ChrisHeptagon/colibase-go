package server

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"encoding/base64"

	"github.com/ChrisHeptagon/colibase/admin/models"
	"github.com/ChrisHeptagon/colibase/admin/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
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
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
		URL:  "/favicon.ico",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4321",
	}))

	publicAdminUIPaths := []string{
		"login",
		"init",
	}
	privateAdminUIPaths := []string{
		"dashboard",
		"settings",
		"users",
		"tables",
	}
	app.Get("/admin-entry/logout", func(c *fiber.Ctx) error {
		return handleUserLogout(c, store)
	})
	app.Post("/api/login", func(c *fiber.Ctx) error {
		return handleUserLogin(c, db, store)
	})
	app.Get("/api/login-schema", func(c *fiber.Ctx) error {
		return loginSchema(c)
	})
	app.Get("/api/user-initialization-status", func(c *fiber.Ctx) error {
		return handleUserInitializatonStatus(c, db, "users")
	})
	app.Post("/api/init-login", func(c *fiber.Ctx) error {
		return handleUserInitializaton(c, db)
	})

	publicGroup := app.Group("/admin-entry")
	privateGroup := app.Group("/admin-ui", func(c *fiber.Ctx) error {
		return AuthMiddleware(c, store)
	})

	for _, path := range publicAdminUIPaths {
		publicGroup.Get(path, func(c *fiber.Ctx) error {
			c.Render("./admin-ui/dist/index.html", fiber.Map{})
			return c.Render("index", fiber.Map{})
		})
	}

	for _, path := range privateAdminUIPaths {
		privateGroup.Get(path, func(c *fiber.Ctx) error {
			return c.Render("index", fiber.Map{})
		})
	}
	app.Static("/assets", "./admin-ui/dist/assets")

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
		return c.Redirect("/admin-entry/login")
	}

	return c.Next()
}

func handleUserLogout(c *fiber.Ctx, st *session.Store) error {
	sess, err := st.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	sess.Destroy()
	c.ClearCookie("colibase")
	return c.Redirect("/admin-entry/login")
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
	c.Redirect("/admin-ui/dashboard")
	return c.JSON(
		fiber.Map{
			"message": "Login Successful",
			"status":  fiber.StatusOK,
		},
	)
}

func handleUserInitializaton(c *fiber.Ctx, db *sql.DB) error {
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

	for key, value := range formData {
		switch key {
		case regexp.MustCompile(`(?i)email`).FindString(key):
			if !regexp.MustCompile(`(?i)^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`).MatchString(value.(string)) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "invalid email",
				})
			}
		case regexp.MustCompile(`(?i)password`).FindString(key):
			HashedPassword, err := utils.HashPassword(value.(string))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			formData[key] = HashedPassword

			if len(value.(string)) < 1 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "password too short",
				})
			}
		default:
			if !regexp.MustCompile(`(?i)^[\w]+$`).MatchString(value.(string)) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "invalid characters",
				})
			}
			continue
		}
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
		if regexp.MustCompile(`(?i)has no column named`).MatchString(err.Error()) {
			db.Exec("DROP TABLE users;")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.SendString("User Initialization Successful")
}

func handleUserInitializatonStatus(c *fiber.Ctx, db *sql.DB, tn string) error {
	if models.IsUserInitialized(db) {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
}

func loginSchema(c *fiber.Ctx) error {
	var userSchema models.DefaultUserSchema
	var userSchemaInterface []interface{}
	emptyStruct := reflect.ValueOf(&userSchema).Elem()
	for i := 0; i < emptyStruct.NumField(); i++ {
		field := emptyStruct.Type().Field(i)
		userSchemaInterface = append(userSchemaInterface, map[string]interface{}{
			"name":      field.Name,
			"form_type": field.Tag.Get("form_type"),
			"required":  field.Tag.Get("required"),
		})
	}
	return c.JSON(userSchemaInterface)
}
