// Filename: test2/cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"time"

	_ "github.com/lib/pq"
	"michaelgomez.net/internal/data"
	"michaelgomez.net/internal/jsonlog"
)

// Version number
const version = "1.0.0"

// configuration settings
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleCoons int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64 //request per second
		burst   int     //how many requests at the intial momment
		enabled bool    //rate limiting toggle
	}
}

// dependency injection
type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

// main
func main() {
	var cfg config

	//server flags
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production)")

	//database flags
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MREF_DB_DSN"), "PostgreSQL DSN") //please remember to create the dsn michael!
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleCoons, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idel time")

	//flags for the rate limiter
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum request per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Rate limiter enabled")

	flag.Parse()

	//creating the logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	//create the connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	//logging the successful connection pool
	logger.PrintInfo("database connection pool established", nil)

	//instance of app struct
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}
	err = app.serve() //starting the server
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

// Database function, will be refactored for test 2
func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleCoons)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	//context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
