package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes/scheme"

	secretsv1alpha1 "github.com/opensecrecy/encrypted-secrets/api/v1alpha1"
)

var (
	Namespace   string
	Provider    string
	Filename    string
	Version     string
	versionFlag bool
)

const (
	SecretApiVersion    = "secrets.opensecrecy.org/v1alpha1"
	DecryptedSecretKind = "DecryptedSecret"
)

var rootCmd = &cobra.Command{
	Use:           "cryptctl",
	Short:         "cryptctl is a command line tool for managing EncryptedSecrets",
	SilenceErrors: true,
	RunE: func(_ *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Printf("cryptctl version: %s\n", Version)
			return nil
		}
		if len(args) == 0 {
			return fmt.Errorf("no command specified")
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVar(&versionFlag, "version", false, "--version")
	_ = secretsv1alpha1.AddToScheme(scheme.Scheme)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
}
