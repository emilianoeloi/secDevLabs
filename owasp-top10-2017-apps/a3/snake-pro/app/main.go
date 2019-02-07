package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/globocom/secDevLabs/owasp-top10-2017-apps/a3/snake-pro/app/api"
	"github.com/globocom/secDevLabs/owasp-top10-2017-apps/a3/snake-pro/app/config"
	db "github.com/globocom/secDevLabs/owasp-top10-2017-apps/a3/snake-pro/app/db/mongo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

func main() {

	fmt.Println("[*] Starting Snake Pro...")

	// loading viper
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		errorAPI(err)
	}
	if err := viper.Unmarshal(&config.APIconfiguration); err != nil {
		errorAPI(err)
	}

	// check if MongoDB is acessible and credentials received are working.
	if _, err := checkMongoDB(); err != nil {
		fmt.Println("[X] ERROR MONGODB: ", err)
		os.Exit(1)
	}

	fmt.Println("[*] MongoDB: OK!")
	fmt.Println("[*] Viper loaded: OK!")

	echoInstance := echo.New()
	echoInstance.HideBanner = true

	echoInstance.Use(middleware.Logger())
	echoInstance.Use(middleware.Recover())
	echoInstance.Use(middleware.RequestID())

	echoInstance.GET("/healthcheck", api.HealthCheck)
	APIport := fmt.Sprintf(":%d", getAPIPort())
	echoInstance.Logger.Fatal(echoInstance.Start(APIport))
}

func errorAPI(err error) {
	fmt.Println("[x] Error starting Snake Pro:")
	fmt.Println("[x]", err)
	os.Exit(1)
}

func getAPIPort() int {
	apiPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		apiPort = 10033
	}
	return apiPort
}

func checkMongoDB() (*db.DB, error) {
	return db.Connect()
}
