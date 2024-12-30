package main

// TODO: include connection with redis
// TODO: include connection with mongodb
// TODO: include connection with Oracle
// TODO: include connection with Mysql
// TODO: include authentication using JWT
// TODO: implement cli arguments
// TODO: improve the readme with examples
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	cf "github.com/marinellirubens/dbwrapper/config"
	pg "github.com/marinellirubens/dbwrapper/database"
	logs "github.com/marinellirubens/dbwrapper/logger"
	cli "github.com/urfave/cli/v2"
)

func exportNotes(cCtx *cli.Context) error {
	fileExport := cCtx.String("output")
	fmt.Println("Notes exported to file "+cCtx.String("output"), fileExport)
	return nil
}

// TODO: need to create some treatment on the path variable to understand how to do that without any framework
func ServeApiNative(address string, port int, app *pg.App) {
	server_path := fmt.Sprintf("%v:%v", address, port)
	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		fmt.Println("Root handler:", r.URL.Path)
	}))

	// process requests with path arguments
	mux.Handle(
		"/pg/",
		http.StripPrefix("/pg/", http.HandlerFunc(app.ProcessPostgresRequestHandlePath)),
	)

	mux.HandleFunc("/oracle", app.ProcessOracleRequest)
	mux.HandleFunc("/mongodb", app.ProcessMongoRequest)

	app.Log.Info(fmt.Sprintf("Starting server on %v", server_path))
	http.ListenAndServe(server_path, mux)
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
