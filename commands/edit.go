package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/opensecrecy/encrypted-secrets/pkg/providers"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"

	secretsv1alpha1 "github.com/opensecrecy/encrypted-secrets/api/v1alpha1"
)

var (
	whitespaceRegexp *regexp.Regexp
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit [flags]",
	Short: "edit encryptedSecrets manifest",
	Long:  "Edit an EncryptedSecret manifest file that contains encrypted secret values",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("filename is required")
		}
		return nil
	},
	RunE: func(_ *cobra.Command, args []string) error {
		fileName := args[0]
		encryptedFile, err := os.ReadFile(fileName)
		if err != nil {
			return fmt.Errorf("error reading file %s", err.Error())
		}

		codecs := serializer.NewCodecFactory(scheme.Scheme, serializer.EnableStrict)
		obj, _, err := codecs.UniversalDeserializer().Decode(encryptedFile, &schema.GroupVersionKind{
			Group:   secretsv1alpha1.GroupVersion.Group,
			Version: SecretApiVersion,
			Kind:    "EncryptedSecret",
		}, nil)
		if err != nil {
			if ok, _ := regexp.MatchString("no kind(.*)is registered for version", err.Error()); ok {
				panic("no kind(.*)is registered for version")
			}
			panic(err)
		}

		// convert the runtimeObj to encryptedSecret object
		encryptedSecret, ok := obj.(*secretsv1alpha1.EncryptedSecret)
		if !ok {
			// should never happen
			panic("failed to convert runtimeObject to encryptedSecret")
		}

		// get the decryptedSecret
		decryptedObj, err := providers.DecodeAndDecrypt(encryptedSecret)
		if err != nil {
			return fmt.Errorf("failed to decrypt value for %s", err.Error())
		}

		// marshal into yaml
		decrypted, err := yaml.Marshal(&decryptedObj)
		if err != nil {
			return fmt.Errorf("error marshaling decryptedSecret %s", err.Error())
		}

		// open editor to edit the decryptedSecret
		editedManitest, err := editObjects(decrypted)
		if err != nil {
			return fmt.Errorf("error editing objects %s", err.Error())
		}

		// unmarshal the edited yaml into decryptedSecrets again to encrypt the new secrets
		var newDecryptedSecret secretsv1alpha1.DecryptedSecret

		err = yaml.Unmarshal(editedManitest, &newDecryptedSecret)
		if err != nil {
			return fmt.Errorf("error unmarshaling file %s", err.Error())
		}

		// encrypt the modified data again
		encryptedObj, err := providers.EncryptAndEncode(newDecryptedSecret)
		if err != nil {
			return fmt.Errorf("error encrypting new secrets %s", err)
		}

		// write the contents to yaml
		newEncrypted, err := yaml.Marshal(&encryptedObj)
		if err != nil {
			return fmt.Errorf("error marshaling encryptedSecret %s", err.Error())
		}

		// just a simple hack to remove status field from the yaml
		re := regexp.MustCompile(`status:[\s\S]*?(?:---|$)`)
		result := re.ReplaceAll(newEncrypted, nil)
		// finally, write the encryptedSecret yaml
		err = os.WriteFile(fileName, result, 0600)
		if err != nil {
			return fmt.Errorf("error writing EncryptedSecret %s", err)
		}

		return nil

	},
}

func editObjects(data []byte) ([]byte, error) {
	manifestBuf := bytes.Buffer{}
	_, _ = manifestBuf.Write(data)

	for {

		// make the YAML to show in the editor.
		editorBuf := bytes.Buffer{}
		_, _ = manifestBuf.WriteTo(&editorBuf)
		editorReader := bytes.NewReader(editorBuf.Bytes())

		// open a temporary file.
		tmpfile, err := os.CreateTemp("", ".*.yaml")
		if err != nil {
			return nil, fmt.Errorf("error making tempfile %s", err.Error())
		}
		defer tmpfile.Close()
		defer os.Remove(tmpfile.Name())
		_, _ = editorReader.WriteTo(tmpfile)
		_ = tmpfile.Sync()

		// show the editor.
		err = runEditor(tmpfile.Name())
		if err != nil {
			return nil, fmt.Errorf("error running editor %s", err.Error())
		}

		// re-read the edited file.
		afterTmpfile, err := os.Open(tmpfile.Name())
		if err != nil {
			return nil, fmt.Errorf("error re-opening tempfile %s %s", tmpfile.Name(), err)
		}
		defer afterTmpfile.Close()
		afterBuf := bytes.Buffer{}
		_, err = afterBuf.ReadFrom(afterTmpfile)
		if err != nil {
			return nil, fmt.Errorf("error reading tempfile %s %s", tmpfile.Name(), err)
		}

		// check if the file was edited at all.
		if bytes.Equal(editorBuf.Bytes(), afterBuf.Bytes()) {
			fmt.Println("Edit cancelled. No changes made")
			os.Exit(0)
		} else {
			// return the edited bytes
			return afterBuf.Bytes(), nil
		}

	}

}

func runEditor(filename string) error {
	whitespaceRegexp = regexp.MustCompile(`\s+`)
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("nso $EDITOR set")
	}

	// deal with an editor that has options.
	editorParts := whitespaceRegexp.Split(editor, -1)
	executable := editorParts[0]
	executable, _ = exec.LookPath(executable)

	editorParts = append(editorParts, filename)
	cmd := exec.Command(executable, editorParts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running editor %s", err)
	}
	return nil
}
