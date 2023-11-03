package commands

import (
	"fmt"

	"github.com/opensecrecy/cryptctl/commands/utils"
	"github.com/spf13/cobra"
)

func init() {

	initCmd.Flags().StringVarP(&Namespace, "namespace", "n", "", "namespace to use (required)")
	initCmd.Flags().StringVarP(&Provider, "provider", "p", "", "provider to use (required)")
	_ = initCmd.MarkFlagRequired("provider")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [flags]",
	Short: "init",
	Long:  "Init initializes the encrypted-secrets CLI",
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if Provider == "k8s" && Namespace == "" {
			return fmt.Errorf("namespace is required for k8s provider")
		}
		return nil

	},

	RunE: func(_ *cobra.Command, args []string) error {

		switch Provider {
		case "k8s":
			return utils.InitK8s(Namespace)

		case "aws-kms":
			return utils.InitAwsKms(Namespace)
		}

		return nil

	},
}
