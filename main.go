package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	bolt "github.com/boltdb/bolt"
	logger "github.com/0187773933/Logger/v1/logger"
	utils "github.com/0187773933/BLANK_SERVER/v1/utils"
	server "github.com/0187773933/BLANK_SERVER/v1/server"
)

var s server.Server
var DB *bolt.DB

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		logger.Log.Println( "\r- Ctrl+C pressed in Terminal" )
		DB.Close()
		logger.Log.Printf( "Shutting Down %s Server" , s.Config.Name )
		s.FiberApp.Shutdown()
		logger.CloseDB()
		os.Exit( 0 )
	}()
}

func main() {
	config := utils.GetConfig()
	// utils.GenerateNewKeysWrite( &config )
	// utils.GenerateNewKeys()
	defer utils.SetupStackTraceReport()
	logger.New( &config.Log )
	DB , _ = bolt.Open( config.Bolt.Path , 0600 , &bolt.Options{ Timeout: ( 3 * time.Second ) } )
	s = server.New( &config , logger.Log , DB )
	SetupCloseHandler()
	s.Start()
}