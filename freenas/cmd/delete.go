package cmd

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/kubernetes-incubator/external-storage/freenas/pkg/provisioner"
)

var deletecontrollerCmd = &cobra.Command{
    Use:   "delete",
    Short: "delete an iscsi volume",
    Long:  `delete an iscsi volume`,
    Run: func(cmd *cobra.Command, args []string) {
        initLog()
        log.Debugln("delete called")
		url := fmt.Sprintf("%s://%s/", viper.GetString("freenas-scheme"), viper.GetString("freenas-address"))
        config := &provisioner.FreeNasConfig{
            Url: url,
            Username: viper.GetString("freenas-username"),
            Password: viper.GetString("freenas-password"),
            Pool: viper.GetString("pool-name"),
        }
		log.Debugf("Connecting to %s", config.Url)
        log.Debugf("Username: %s", config.Username)
        log.Debugf("Password: %s", config.Password)
        log.Debugf("Using pool: %s", config.Pool)

        name := viper.GetString("name")
        log.Debugf("Volume Name: %s", name)
        err := provisioner.DeleteVolume(config, name)
        if err != nil {
            log.Fatal("Error deleing volume: Name = %s", name)
            os.Exit(1)
        }
    },
}
func init() {
    deletecontrollerCmd.Flags().String("name", "", "Volume Name")
    viper.BindPFlag("name", deletecontrollerCmd.Flags().Lookup("name"))
	RootCmd.AddCommand(deletecontrollerCmd)
}


