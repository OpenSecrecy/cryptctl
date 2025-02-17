package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	utils "github.com/opensecrecy/cryptctl/commands/utils"
	secretsv1alpha1 "github.com/opensecrecy/encrypted-secrets/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var creationTimeStampRegexp *regexp.Regexp

func init() {

	creationTimeStampRegexp = regexp.MustCompile(`.*creationTimestamp.*\n`)
	createCmd.Flags().StringVarP(&Provider, "provider", "p", "", "provider to use (required)")
	_ = createCmd.MarkFlagRequired("provider")
	createCmd.Flags().StringVarP(&Filename, "filename", "f", "", "filename to use (required)")
	_ = createCmd.MarkFlagRequired("filename")
	rootCmd.AddCommand(createCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [flags] <cluster_name>",
	Short: "create encryptedSecrets manifest",
	Long:  "Create an EncryptedSecretinstance manifest file that will contain encrypted secret values",
	RunE: func(_ *cobra.Command, args []string) error {

		sampleEncryptedSecret := &secretsv1alpha1.EncryptedSecret{
			ObjectMeta: v1.ObjectMeta{
				Name:      "encryptedsecret-sample",
				Namespace: "default",
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: SecretApiVersion,
				Kind:       "EncryptedSecret",
			},
		}

		// set some default labels
		sampleEncryptedSecret.SetLabels(map[string]string{
			"app.kubernetes.io/name":       "encryptedsecret",
			"app.kubernetes.io/part-of":    "encryted-secrets",
			"app.kubernetes.io/created-by": "encryted-secrets",
		})

		// set provider annotation
		sampleEncryptedSecret.SetAnnotations(map[string]string{
			"secrets.opensecrecy.org/provider": Provider,
		})

		// write the contents to yaml
		newEncrypted, err := yaml.Marshal(&sampleEncryptedSecret)
		if err != nil {
			return fmt.Errorf("error marshaling encryptedSecret %s", err.Error())
		}

		yaml := utils.RemoveStatus(newEncrypted)
		yaml = creationTimeStampRegexp.ReplaceAll(yaml, nil)

		err = os.WriteFile(Filename, yaml, 0600)
		if err != nil {
			return fmt.Errorf("error writing EncryptedSecret %s", err)
		}

		return nil

	},
}
