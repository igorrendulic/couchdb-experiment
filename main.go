package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/igorrendulic/couchdb-experiment/api"
	"github.com/igorrendulic/couchdb-experiment/docs"
	"github.com/igorrendulic/couchdb-experiment/global"
	services "github.com/igorrendulic/couchdb-experiment/services/couchdb"

	cfg "github.com/mailio/go-web3-kit/config"
	w3srv "github.com/mailio/go-web3-kit/gingonic"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Web3 Go Kit basic server
// @version 1.0
// @description This is a basic server example using go-web3-kit

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	var (
		configFile string
	)

	// configuration file optional path. Default:  current dir with  filename conf.yaml
	flag.StringVar(&configFile, "c", "conf.yaml", "Configuration file path.")
	flag.StringVar(&configFile, "config", "conf.yaml", "Configuration file path.")
	flag.Usage = usage
	flag.Parse()

	// var conf g.Config
	err := cfg.NewYamlConfig(configFile, &global.Conf)
	if err != nil {
		global.Logger.Log(err, "conf.yaml failed to load")
		panic("Failed to load conf.yaml")
	}

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Go Web3 Kit API"
	docs.SwaggerInfo.Description = "This is a basic server example using go-web3-kit"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", global.Conf.Host, global.Conf.Port)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{global.Conf.Scheme}
	ginSwagger.DefaultModelsExpandDepth(1)

	// server wait to shutdown monitoring channels
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	// init routing (for endpoints)
	router := w3srv.NewAPIRouter(&global.Conf.YamlConfig)

	couchDBService := services.NewCouchDB()

	root := router.Group("/api")
	{

		root.GET("/pong", api.PingPongAPI)
		root.POST("/v1/register", api.NewUserAPI(couchDBService).RegisterUser)
		root.POST("/v1/email", api.NewEmailAPI(couchDBService).AddEmail)
		// root.GET("/v1/email", api.NewEmailAPI(couchDBService).ListEmails)
	}

	// start server
	srv := w3srv.Start(&global.Conf.YamlConfig, router)
	// wait for server shutdown
	go w3srv.Shutdown(srv, quit, done)

	global.Logger.Log("Server is ready to handle requests at", global.Conf.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		global.Logger.Log("Could not listen on %s: %v\n", global.Conf.Port, err)
	}

	<-done

}

// usage will print out the flag options for the server.
func usage() {
	usageStr := `Usage: operator [options]
	Server Options:
	-c, --config <file>              Configuration file path
`
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}
