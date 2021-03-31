package configurations

import (
	"log"
	"os"
)

func requireEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	} else {
		log.Fatalf("Environment variable %s isn't set", key)
		return ""
	}
}
