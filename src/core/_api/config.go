package api

import (
	"os"
)

// GetEnv : python like get env :)
func GetEnv(key ...string) string {
	switch len(key) {
	case 0:
		panic("getEnv: No Key provided!")
	case 1:
		val, ok := os.LookupEnv(key[0])
		if ok {
			return val
		} else {
			panic("getEnv: No Values for provided Key !")
		}
	case 2:
		val, ok := os.LookupEnv(key[0])
		if ok {
			return val
		} else {
			return key[1]
		}
	default:
		panic("getEnv: Only Two parameters is needed!")
	}
}

var (
	MONGO_URL = GetEnv("MONGO_URL", "mongodb://127.0.0.1:27017")
	PORT      = GetEnv("PORT", "8080")
)
