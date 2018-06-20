package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/kubernetes-incubator/external-storage/lib/controller"
    "github.com/kubernetes-incubator/external-storage/freenas/pkg/provisioner"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "github.com/Sirupsen/logrus"
)


// start-controllerCmd represents the start-controller command
var startcontrollerCmd = &cobra.Command{
    Use:   "start",
    Short: "Start a freenas iscsi dynamic provisioner",
    Long:  `Start a freenas iscsi dynamic provisioner`,
    Run: func(cmd *cobra.Command, args []string) {
        initLog()
        log.Debugln("start called")
        var config *rest.Config
        var err error
        master := viper.GetString("master")
        log.Debugln(master)
        kubeconfig := viper.GetString("kubeconfig")
        log.Debugln(kubeconfig)
        // creates the in-cluster config
        log.Debugln("creating in cluster default kube client config")
        if master != "" || kubeconfig != "" {
            config, err = clientcmd.BuildConfigFromFlags(master, kubeconfig)
        } else {
            config, err = rest.InClusterConfig()
        }
        if err != nil {
            log.Fatalln(err)
        }
        log.WithFields(logrus.Fields{
            "config-host": config.Host,
        }).Debugln("kube client config created")

        // creates the clientset
        log.Debugln("creating kube client set")
        kubernetesClientSet, err := kubernetes.NewForConfig(config)
        if err != nil {
            log.Fatalln(err)
        }
        log.Debugln("kube client set created")

        // The controller needs to know what the server version is because out-of-tree
        // provisioners aren't officially supported until 1.5
        serverVersion, err := kubernetesClientSet.Discovery().ServerVersion()
        if err != nil {
            log.Fatalf("Error getting server version: %v", err)
        }
        url := fmt.Sprintf("%s://%s/", viper.GetString("scheme"), viper.GetString("address"))
        iqn := viper.GetString("iqn")
        log.Debugf("IQN: %s", iqn)
        freeNasConfig := &provisioner.FreeNasConfig{
            Uri: url,
            Pool: viper.GetString("pool"),
            Username: viper.GetString("username"),
            Password: viper.GetString("password"),
            Portal: fmt.Sprintf("%s:%d", viper.GetString("address"), 3260),
            IQN: iqn,
        }
        freenasProvisioner := provisioner.NewFreenasProvisioner(freeNasConfig)

        pc := controller.NewProvisionController(kubernetesClientSet, viper.GetString("provisioner-name"), freenasProvisioner, serverVersion.GitVersion)
        log.Debugln("freenas iscsi controller created, running forever...")
        pc.Run(wait.NeverStop)

    },
}
func init() {
    RootCmd.AddCommand(startcontrollerCmd)

    startcontrollerCmd.Flags().String("scheme", "https", "scheme of the freenas connection, can be http or https")
    viper.BindPFlag("scheme", startcontrollerCmd.Flags().Lookup("scheme"))
    startcontrollerCmd.Flags().String("username", "", "Freenas username")
    viper.BindPFlag("username", startcontrollerCmd.Flags().Lookup("username"))
    startcontrollerCmd.Flags().String("password", "", "Freenas password")
    viper.BindPFlag("password", startcontrollerCmd.Flags().Lookup("password"))
    startcontrollerCmd.Flags().String("address", "", "Freenas address")
    viper.BindPFlag("address", startcontrollerCmd.Flags().Lookup("address"))
    startcontrollerCmd.Flags().String("iqn", "iqn.2005-10.org.freenas.ctl", "Freenas IQN")
    viper.BindPFlag("iqn", startcontrollerCmd.Flags().Lookup("iqn"))
    startcontrollerCmd.Flags().String("pool", "tank", "name of the freenas zpool")
    viper.BindPFlag("pool", startcontrollerCmd.Flags().Lookup("pool"))


    startcontrollerCmd.Flags().String("master", "", "Master URL")
    viper.BindPFlag("master", startcontrollerCmd.Flags().Lookup("master"))
    startcontrollerCmd.Flags().String("kubeconfig", "", "Absolute path to the kubeconfig")
    viper.BindPFlag("kubeconfig", startcontrollerCmd.Flags().Lookup("kubeconfig"))

    startcontrollerCmd.Flags().Duration("resync-period", controller.DefaultResyncPeriod, "how often to poll the master API for updates")
    viper.BindPFlag("resync-period", startcontrollerCmd.Flags().Lookup("resync-period"))
    startcontrollerCmd.Flags().Bool("exponential-backoff-on-error", controller.DefaultExponentialBackOffOnError, "exponential-backoff-on-error doubles the retry-period everytime there is an error")
    viper.BindPFlag("exponential-backoff-on-error", startcontrollerCmd.Flags().Lookup("exponential-backoff-on-error"))
    startcontrollerCmd.Flags().Int("fail-retry-threshold", controller.DefaultFailedProvisionThreshold, "Threshold for max number of retries on failure of provisioner")
    viper.BindPFlag("fail-retry-threshold", startcontrollerCmd.Flags().Lookup("fail-retry-threshold"))
    startcontrollerCmd.Flags().Duration("lease-period", controller.DefaultLeaseDuration, "LeaseDuration is the duration that non-leader candidates will wait to force acquire leadership. This is measured against time of last observed ack")
    viper.BindPFlag("lease-period", startcontrollerCmd.Flags().Lookup("lease-period"))
    startcontrollerCmd.Flags().Duration("renew-deadline", controller.DefaultRenewDeadline, "RenewDeadline is the duration that the acting master will retry refreshing leadership before giving up")
    viper.BindPFlag("renew-deadline", startcontrollerCmd.Flags().Lookup("renew-deadline"))
    startcontrollerCmd.Flags().Duration("retry-period", controller.DefaultRetryPeriod, "RetryPeriod is the duration the LeaderElector clients should wait between tries of actions")
    viper.BindPFlag("retry-period", startcontrollerCmd.Flags().Lookup("retry-period"))
    startcontrollerCmd.Flags().Duration("term-limit", controller.DefaultTermLimit, "TermLimit is the maximum duration that a leader may remain the leader to complete the task before it must give up its leadership. 0 for forever or indefinite.")
    viper.BindPFlag("term-limit", startcontrollerCmd.Flags().Lookup("term-limit"))

    startcontrollerCmd.Flags().String("provisioner-name", "freenas-provisioner", "name of this provisioner, must match what is passed int the storage class annotation")
    viper.BindPFlag("provisioner-name", startcontrollerCmd.Flags().Lookup("provisioner-name"))

    startcontrollerCmd.Flags().String("default-fs", "ext4", "filesystem to use when not specified")
    viper.BindPFlag("default-fs", startcontrollerCmd.Flags().Lookup("default-fs"))

}
