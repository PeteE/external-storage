package cmd

import (
    "fmt"
    "os"
    "github.com/Sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "strings"
)

var log = logrus.New()

var RootCmd = &cobra.Command{
    Use:   "freenas-controller",
    Short: "an iscsi freenas dynamic provisioner for kubernetes",
    Long:  "an iscsi freenas dynamic provisioner for kubernetes",
}

func Execute() {
    if err := RootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(-1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)
    RootCmd.PersistentFlags().String("log-level", "debug", "log level")
    viper.BindPFlag("log-level", RootCmd.PersistentFlags().Lookup("log-level"))
}

func initConfig() {
    viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
    viper.AutomaticEnv()
}

func initLog() {
    var err error
    log.Level, err = logrus.ParseLevel(viper.GetString("log-level"))
    if err != nil {
        log.Fatalln(err)
    }
}
