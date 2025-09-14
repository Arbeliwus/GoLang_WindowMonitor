package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr string
	DSN      string
}

func Load() Config {
	_ = godotenv.Load(".env")

	dsn := os.Getenv("DB_URL")

	if dsn == "" {
		user := getenv("DB_USER", "postgres")
		pass := os.Getenv("DB_PASS") // 密碼不給預設，避免誤連
		host := getenv("DB_HOST", "127.0.0.1")
		port := getenv("DB_PORT", "5432")
		name := getenv("DB_NAME", "iot_evergrain")
		ssl := getenv("DB_SSLMODE", "disable")

		if pass == "" {
			log.Fatal("missing DB_PASS (database password)")
		}

		// key=value 形式；密碼含特殊字元用單引號即可
		dsn = fmt.Sprintf(
			"user=%s password='%s' host=%s port=%s dbname=%s sslmode=%s",
			user, pass, host, port, name, ssl,
		)
	}

	return Config{
		HTTPAddr: getenv("HTTP_ADDR", ":8100"),
		DSN:      dsn,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
