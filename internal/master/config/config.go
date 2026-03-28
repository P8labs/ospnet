package config

import "strings"

type Config struct {
	Port           string
	DatabaseURL    string
	BaseURL        string
	TsAuthKey      string
	TailnetID      string
	FrontendURL    string
	AllowedOrigins []string
}

func Load() *Config {

	cors := GetEnv("CORS_ORIGINS", "*")

	raw := strings.Split(cors, ",")
	allowedCors := make([]string, 0, len(raw))

	for _, o := range raw {
		o = strings.TrimSpace(o)
		if o != "" {
			allowedCors = append(allowedCors, o)
		}
	}

	cfg := &Config{
		Port:           GetEnv("PORT", "8000"),
		DatabaseURL:    MustEnv("DATABASE_URL"),
		BaseURL:        GetEnv("BASE_URL", "http://localhost:8000"),
		FrontendURL:    GetEnv("FRONTEND_URL", "http://localhost:5173"),
		TsAuthKey:      MustEnv("TS_AUTH_KEY"),
		TailnetID:      MustEnv("TS_ID"),
		AllowedOrigins: allowedCors,
	}

	return cfg
}
