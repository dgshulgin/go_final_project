package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config int

func (c Config) getEnv(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("переменная %s не задана", name)
	}
	return value, nil
}

func (c Config) GetEnvAsString(name string, defaultValue string) string {
	value, err := c.getEnv(name)
	if err != nil {
		return defaultValue
	}
	return value
}
func (c Config) GetEnvAsInt(name string, defaultValue int) int {
	value, err := c.getEnv(name)
	if err != nil {
		return defaultValue
	}
	value0, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return value0
}
