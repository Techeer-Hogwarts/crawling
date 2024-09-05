package config

import (
	"log"
	"os"
)

func GetEnv(key string, defaultVal string) string {
	value, found := os.LookupEnv(key)
	if !found {
		return defaultVal
	}
	log.Printf("env %s: found \n", key)
	return value
}
