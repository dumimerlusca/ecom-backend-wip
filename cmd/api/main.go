package main

import (
	"database/sql"
	"ecom-backend/internal/handlers"
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"ecom-backend/pkg/sqldb"
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	cfg        config
	logger     *jsonlog.Logger
	db         *sql.DB
	services   *service.Services
	middleware *handlers.Middleware
}

func main() {
	fmt.Println("Server runninng")

	var cfg config
	logger := jsonlog.NewLoger(os.Stdout, jsonlog.LevelInfo)

	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25,
		"PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25,
		"PostgreSQL max open idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m",
		"PostgreSQL max connection idle time")

	flag.Parse()

	db, err := sqldb.OpenDB(sqldb.DbConfig{Dsn: cfg.db.dsn, MaxOpenConns: cfg.db.maxOpenConns, MaxIdleConns: cfg.db.maxIdleConns, MaxIdleTime: cfg.db.maxIdleTime})

	if err != nil {
		logger.PrintFatal(err, nil)
	}

	app := application{cfg: cfg, logger: logger, db: db}
	app.initServices()
	app.middleware = handlers.NewMiddleware(logger, app.services)

	err = app.serve()

	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func (app *application) initServices() {
	db := app.db
	models := model.NewModels(db)

	app.services = service.NewServices(db, models)
}
