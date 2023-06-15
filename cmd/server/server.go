package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	application "github.com/MyLi2tlePony/messenger/internal/app"
	loggerConfig "github.com/MyLi2tlePony/messenger/internal/config/logger"
	migrationConfig "github.com/MyLi2tlePony/messenger/internal/config/migration"
	serverConfig "github.com/MyLi2tlePony/messenger/internal/config/server"
	"github.com/MyLi2tlePony/messenger/internal/migration"
	httpsrv "github.com/MyLi2tlePony/messenger/internal/server/http"
	"github.com/MyLi2tlePony/messenger/internal/storage/postgres"

	databaseConfig "github.com/MyLi2tlePony/messenger/internal/config/database"
)

var configPath string

func init() {
	defaultConfigPath := path.Join("configs", "server", "config.toml")
	flag.StringVar(&configPath, "config", defaultConfigPath, "Path to configuration file")
}

func main() {
	// Получение конфигурации базы данных
	dbConfig, err := databaseConfig.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Получение конфигурации миграции
	migrationConf, err := migrationConfig.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Получение конфигурации базы данных
	dbConnectionString := dbConfig.GetConnectionString()
	postgresDsn := os.Getenv("POSTGRES_DSN")
	if postgresDsn != "" {
		dbConnectionString = postgresDsn
	}
	ctx := context.Background()

	// Инициализация структуры, отвечающей за миграции
	migrator, err := migration.Connect(ctx, dbConnectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Выполнение миграции
	err = migrator.Migrate(ctx, migrationConf.UpPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Подключение к базе данных
	db, err := postgres.Connect(ctx, dbConnectionString)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Инициализация уровня приложения
	apps := application.New(db)

	// Получение конфигурации логера
	logConfig, err := loggerConfig.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Инициализация сервера
	srv := httpsrv.NewServer(logConfig.GetLevel(), apps)
	srvConfig, err := serverConfig.New(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		// Остановка сервера
		if err := srv.Stop(); err != nil {
			fmt.Println("failed to stop serv: " + err.Error())
		}
	}()

	fmt.Println("app is running...")

	// Запуск сервера
	err = srv.Start(srvConfig.GetHostPort())
	if err != nil {
		fmt.Println(err)
	}
}
