package cmd

import (
    "os"
    "os/signal"
    "syscall"

    "github.com/victoryang/api-gateway/api"
    "github.com/victoryang/api-gateway/config"
    //"github.com/victoryang/api-gateway/db"
    Log "github.com/victoryang/api-gateway/log"

    "github.com/spf13/cobra"
)

func handleSignals(apiserver *api.Server) error {
    signal.Ignore()
    signalQueue := make(chan os.Signal)
    signal.Notify(signalQueue, syscall.SIGHUP, os.Interrupt)

    for {
        sig := <-signalQueue
        switch sig {
        //TODO:
        //case syscall.SIGHUP:
            //reload config file
        default:
            apiserver.Shutdown()

            stopAdminServer()

            //db.CloseDB()

            Log.CloseFile()

            return nil
        }
    }
}

func runDaemon(cmd *cobra.Command, args []string) error {
    cfg := LoadConfig()

    if err := ConfigServerLog(cfg); err!=nil {
        return returnError(ERR_OPEN_LOG_FILE_FAIL)
    }

    startAdminServer(cfg.Admin)

    //SetUpDatabase(cfg.Databases)

    apiserver := api.NewApiServer(cfg)
    if apiserver == nil {
        Log.Error("Error in starting apiserver")
        return returnError(ERR_START_APISERVER)
    }
    apiserver.Run()

    return handleSignals(apiserver)
}