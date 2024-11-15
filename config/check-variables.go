package config

import (
	"log"
	"os"
	logger "ruantiengo/log"

	"github.com/joho/godotenv"
)

func CheckVariables() {
	// Primeiro carrega o arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
	logger.Message(logger.Info, "✔️ Arquivo .env carregado com sucesso.")

	// Depois verifica as variáveis
	requiredEnvVars := []string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",

		"DEBUG_MODE",

		"RABBITMQ_URI",
	}

	missingVars := []string{}
	for _, envVar := range requiredEnvVars {
		if value := os.Getenv(envVar); value == "" {
			missingVars = append(missingVars, envVar)
		}
	}

	if len(missingVars) > 0 {
		logger.Message(logger.Error, "------------------------------------------------------")
		logger.Message(logger.Error, "Variáveis de ambiente faltando:")
		for _, envVar := range missingVars {
			logger.Message(logger.Error, "  ➡️ %s", envVar)
		}
		logger.Message(logger.Error, "------------------------------------------------------")
		logger.Message(logger.Error, "❌ A aplicação não pode ser executada até que todas as variáveis de ambiente sejam configuradas.")
		logger.Message(logger.Error, "Consulte o arquivo .env para definir as variáveis faltantes.")
		os.Exit(1)
	}

	logger.Message(logger.Info, "✔️ Todas as variáveis de ambiente estão configuradas corretamente.")
}
