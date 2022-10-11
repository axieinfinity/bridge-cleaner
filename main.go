package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/axieinfinity/bridge-cleaner/cleaners"
	"github.com/axieinfinity/bridge-cleaner/configs"
	"github.com/axieinfinity/bridge-core/stores"
	storesBr "github.com/axieinfinity/bridge-v2/stores"
	"github.com/ethereum/go-ethereum/log"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"
)

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(configs.AppConfig.LogLevel), log.StreamHandler(os.Stderr, log.TerminalFormat(true))))

	db, err := stores.MustConnectDatabase(&stores.Database{
		Host:     configs.AppConfig.DB.Host,
		User:     configs.AppConfig.DB.User,
		Password: configs.AppConfig.DB.Password,
		DBName:   configs.AppConfig.DB.DBName,
		Port:     configs.AppConfig.DB.Port,
	}, false)
	if err != nil {
		panic(err)
	}
	listenHandlerStore := storesBr.NewBridgeStore(db)
	bridgeStore := stores.NewMainStore(db)

	cleaner := cleaners.NewCleaner(listenHandlerStore, bridgeStore)
	cleaner.Start()
	log.Info("Start running")

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
}

func init() {
	if _, err := configs.New(); err != nil {
		panic(err)
	}
}
