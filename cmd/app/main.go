package main

// TODO: include connection with redis
// TODO: include connection with mongodb
// TODO: include connection with Oracle
// TODO: include connection with Mysql
// TODO: include authentication using JWT
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	pg "github.com/marinellirubens/dbwrapper/database"
	cf "github.com/marinellirubens/dbwrapper/internal/config"
	logs "github.com/marinellirubens/dbwrapper/internal/logger"
	cli "github.com/urfave/cli/v2"
)

const VERSION = "1.0.0"

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func SetupRoutes(mux *http.ServeMux, app *pg.App) (http.Handler, error) {
	// TODO: need to create some treatment on the path variable to understand how to do path handling without any framework
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.Log.Debug(fmt.Sprintf("Requested ping by %v", ReadUserIP(r)))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("I'm alive.\n"))

		if err != nil {
			app.Log.Error(fmt.Sprintf("Error trying to get server. %v", err))
			panic(err)
		}
	}))

	// process base request for postgresl
	mux.HandleFunc("/pg", app.ProcessPostgresRequest)

	// process requests with path arguments
	mux.Handle(
		"/pg/",
		http.StripPrefix("/pg/", http.HandlerFunc(app.ProcessPostgresRequestHandlePath)),
	)

	mux.HandleFunc("/oracle", app.ProcessOracleRequest)
	mux.HandleFunc("/mongodb", app.ProcessMongoRequest)
	return corsMiddleware(mux), nil
}

func ServeApiNative(address string, port int, app *pg.App) {
	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()

	handler, err := SetupRoutes(mux, app)
	if err != nil {
		panic("Error setting up the routes")
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	app.Log.Info(fmt.Sprintf("Starting server on %v", server_path))

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

func run_server(cfgPath string) {
	cfg, err := cf.GetInfoFile(cfgPath)
	if err != nil {
		log.Fatal("error processing configuration file")
		os.Exit(1)
	}
	listDelimiter := cfg.Section("DEFAULT").Key("LIST_DELIMITER").String()
	//fmt.Println("List delimiter", listDelimiter)
	//fmt.Println(cfg.SectionStrings())
	//fmt.Println(cfg.Section("POSTGRES").Key("barlist").Strings(listDelimiter))
	for i, v := range cfg.Section("POSTGRES").Key("barlist").Strings(listDelimiter) {
		fmt.Printf("Barlist item %d: %v\n", i, v)
	}
	logger, err := logs.CreateLogger(
		cfg.Section("SERVER").Key("LOGGER_FILE").String(),
		logs.DEBUG,
		[]uint16{logs.STREAM_WRITER, logs.FILE_WRITER},
	)
	if err != nil {
		log.Fatal(err)
	}

	host := cfg.Section("SERVER").Key("SERVER_ADDRESS").String()
	port, _ := cfg.Section("SERVER").Key("SERVER_PORT").Int()

	psqlInfom := pg.GetConnectionInfo(cfg)
	db, err := pg.ConnectToPsql(psqlInfom)
	if err != nil {
		panic(err)
	}
	defer pg.CloseConn(db)

	if err := db.Ping(); err != nil {
		panic(err)
	}
	application := &pg.App{Log: logger}
	application.IncludeDbConnection(db, reflect.TypeOf(pg.PostgresHandler{}), psqlInfom)
	ServeApiNative(host, port, application)
}

func main() {
	var cfgPath string

	app := &cli.App{
		UseShortOptionHandling: true,
		EnableBashCompletion:   true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config_file",
				Usage:       "Path for the configuration file",
				Aliases:     []string{"f"},
				Value:       cf.DefaultCfgFilePath,
				Destination: &cfgPath,
			},
			&cli.BoolFlag{
				Name:    "version",
				Usage:   "Path for the configuration file",
				Aliases: []string{"v"},
				Value:   false,
				Action: func(ctx *cli.Context, b bool) error {
					fmt.Printf("Version: %v\n", VERSION)
					os.Exit(0)
					return nil
				},
			},
		},
		Action: func(cCtx *cli.Context) error {
			fmt.Println("path: ", cfgPath)
			run_server(cfgPath)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
