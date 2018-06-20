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

	RootCmd.PersistentFlags().String("freenas-scheme", "https", "scheme of the freenas connection, can be http or https")
	viper.BindPFlag("freenas-scheme", RootCmd.PersistentFlags().Lookup("freenas-scheme"))

    RootCmd.PersistentFlags().String("freenas-username", "", "Freenas username")
	viper.BindPFlag("freenas-username", RootCmd.PersistentFlags().Lookup("freenas-username"))

	RootCmd.PersistentFlags().String("freenas-password", "", "Freenas password")
	viper.BindPFlag("freenas-password", RootCmd.PersistentFlags().Lookup("freenas-password"))

	RootCmd.PersistentFlags().String("freenas-address", "", "Freenas address")
	viper.BindPFlag("freenas-address", RootCmd.PersistentFlags().Lookup("freenas-address"))

	RootCmd.PersistentFlags().String("freenas-iqn", "iqn.2005-10.org.freenas.ctl", "Freenas IQN")
	viper.BindPFlag("freenas-iqn", RootCmd.PersistentFlags().Lookup("freenas-iqn"))

	RootCmd.PersistentFlags().String("pool-name", "wheatpool", "name of the freenas zpool")
	viper.BindPFlag("pool-name", RootCmd.PersistentFlags().Lookup("pool-name"))

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
