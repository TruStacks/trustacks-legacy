package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	toolchain "github.com/trustacks/trustacks/pkg/provision"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// toolchain cli command flags.
var (
	toolchainName       string
	toolchainConfig     string
	toolchainKubeconfig string
	toolchainForce      bool
)

// toolchainCmd contains subcommands for managing factories.
var toolchainCmd = &cobra.Command{
	Use:   "toolchain",
	Short: "manage toolchains",
}

// toolchainInstallCmd install the toolchain.
var toolchainInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "install a toolchain",
	Run: func(cmd *cobra.Command, args []string) {
		if err := toolchain.Install(toolchainConfig, toolchainForce, git.PlainClone); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var toolchainDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroy a toolchain",
	Run: func(cmd *cobra.Command, args []string) {
		stdin := bufio.NewReader(os.Stdin)
		fmt.Println("\033[0;93mWARNING: \033[3mthis action is destructive\033[0m")
		fmt.Printf("please type the name of the toolchain to proceed [\033[1;95m%s\033[0m]:\n> ", toolchainName)
		line, _, err := stdin.ReadLine()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if string(line) != toolchainName {
			fmt.Println("the toolchain name did not match. aborting")
			os.Exit(1)
		}
		config, err := clientcmd.BuildConfigFromFlags("", toolchainKubeconfig)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err := toolchain.Destory(toolchainName, clientset); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("the toolchain has been deleted. the kubernetes resources will be cleaned up in the background")
	},
}

func init() {
	toolchainCmd.AddCommand(toolchainInstallCmd)
	toolchainInstallCmd.Flags().StringVar(&toolchainConfig, "config", "", "configuration file")
	if err := toolchainInstallCmd.MarkFlagRequired("config"); err != nil {
		log.Fatal(err)
	}
	toolchainInstallCmd.Flags().BoolVar(&toolchainForce, "force", false, "force update (experimental: use at your own risk)")
	rootCmd.AddCommand(toolchainCmd)

	toolchainCmd.AddCommand(toolchainDestroyCmd)
	toolchainDestroyCmd.Flags().StringVar(&toolchainName, "name", "", "name of the toolchain")
	if err := toolchainDestroyCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}
	// add the kubeconfig
	if home := homedir.HomeDir(); home != "" {
		toolchainDestroyCmd.Flags().StringVar(&toolchainKubeconfig, "kubeconfig", filepath.Join(home, ".kube", "config"), "kubeconfig path (absolute path)")
	} else {
		toolchainDestroyCmd.Flags().StringVar(&toolchainKubeconfig, "kubeconfig", "", "kubeconfig path (absolute path)")
		if err := toolchainDestroyCmd.MarkFlagRequired("kubeconfig"); err != nil {
			log.Fatal(err)
		}
	}
}
