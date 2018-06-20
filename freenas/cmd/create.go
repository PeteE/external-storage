package cmd

import (
    "os"
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/kubernetes-incubator/external-storage/freenas/pkg/provisioner"
)

var createcontrollerCmd = &cobra.Command{
    Use:   "create",
    Short: "create an iscsi volume",
    Long:  `create an iscsi volume`,
    Run: func(cmd *cobra.Command, args []string) {
        initLog()
        log.Debugln("create called")
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

        vol := viper.GetString("vol-name")
        size := int64(viper.GetInt("size"))
        log.Debugf("Volume Name: %s, Size: %d", vol, size)

        _, err := provisioner.CreateVolume(config, vol, size)
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }
        log.Debugln("volume created with Name = ", vol)
    },
}
func init() {
    createcontrollerCmd.Flags().String("vol-name", "", "Volume Name")
    viper.BindPFlag("vol-name", createcontrollerCmd.Flags().Lookup("vol-name"))
    createcontrollerCmd.Flags().Uint("size", 1, "Volume Size in GiB")
    viper.BindPFlag("size", createcontrollerCmd.Flags().Lookup("size"))
	RootCmd.AddCommand(createcontrollerCmd)
}
