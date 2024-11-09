package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sped-efinanceira/database"
	"sped-efinanceira/database/seeders"
	"sped-efinanceira/routes"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Carrega as variáveis de ambiente do arquivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}

	// Configurações do MongoDB
	dbURL := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")

	if dbURL == "" || dbName == "" {
		log.Fatal("As variáveis de ambiente DB_URL e DB_NAME devem ser definidas.")
	}

	// Testa a conexão com o MongoDB
	client, _, err := database.Connect(dbURL, dbName)

	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer client.Disconnect(nil)

	if database.CheckConnection(client) {
		log.Println("✅ Conexão com o MongoDB confirmada e funcionando!")
	} else {
		log.Println("⚠️ Falha ao verificar a conexão com o MongoDB")
	}

	// Verifica se o diretório "logs" existe, caso contrário, cria-o
	err = os.Mkdir("logs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Cria o arquivo de log
	logFile, err := os.Create("logs/server.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	// Configura o logger para usar o arquivo de log
	log.SetOutput(logFile)
	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)

	// Sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	// Cria um roteador principal com Mux
	router := routes.ConfiguraRotas(client)

	// Configuração personalizada do CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Configuração dos repositórios
	perfilRepo := seeders.ConfiguraRepositorios()

	// Criação das Seeders
	seeders.SeedPerfis(&perfilRepo)

	// Obtém o IP local
	ip, err := getLocalIP()
	if err != nil {
		log.Fatalf("Erro ao obter o IP local: %v", err)
	}

	// Define a porta do servidor
	port := ":8080"
	address := fmt.Sprintf("%s%s", ip, port)

	// Inicia o servidor com CORS
	log.Printf("🟢 Servidor rodando em http://%s\n", address)
	if err := http.ListenAndServe(address, cors(router)); err != nil {
		log.Printf("Erro ao iniciar servidor: %v\n", err)
	}
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("não foi possível encontrar o IP local")
}
