package configs

import (
	"errors"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func FromEnv() (DatabaseConfig, error) {
	dbConfig := DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
	var err error

	dbConfig.Port, err = strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return DatabaseConfig{}, errors.New("Could not get DB_PORT:")
	}

	if dbConfig.Host == "" || dbConfig.User == "" || dbConfig.Password == "" || dbConfig.DBName == "" {
		return DatabaseConfig{}, errors.New("One or more database configuration parameters are empty")
	}

	return dbConfig, nil
}

func FromEnvRedis() (RedisConfig, error) {
	redisConfig := RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	if redisConfig.Host == "" || redisConfig.Port == "" || redisConfig.Password == "" {
		return RedisConfig{}, errors.New("One or more redis configuration parameters are empty")
	}

	return redisConfig, nil
}
