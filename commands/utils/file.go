// this file contains the file related functions

package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"

	secretsv1alpha1 "github.com/opensecrecy/encrypted-secrets/api/v1alpha1"
)

const (
	SecretApiVersion = "secrets.shubhindia.xyz/v1alpha1"
)

type Objects struct {
	EncryptedSecret secretsv1alpha1.EncryptedSecret
	Objects         []string
}

func ParseYaml(filename string) (*Objects, error) {
	encryptedFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s", err.Error())
	}

	yamlDocuments := strings.Split(string(encryptedFile), "---\n")

	codecs := serializer.NewCodecFactory(scheme.Scheme, serializer.EnableStrict)
	obj, _, err := codecs.UniversalDeserializer().Decode([]byte(yamlDocuments[0]), &schema.GroupVersionKind{
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

	objs := Objects{
		EncryptedSecret: *encryptedSecret,
		Objects:         yamlDocuments[1:],
	}
	return &objs, nil
}
