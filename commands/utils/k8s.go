package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/opensecrecy/encrypted-secrets/pkg/providers/utils"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func InitK8s(namespace string) error {
	// create a certificate which is valid for 10 years
	cert, err := GeneratePrivateKeyAndCert(2048, 10*365*24*time.Hour, "cryptctl-key")
	if err != nil {
		panic(err)

	}

	// create a secret with the certificate
	client, err := utils.GetKubeClient()
	if err != nil {
		return err
	}

	secret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cryptctl-key",
			Namespace: namespace,
		},
		Type: "Opaque",
		Data: map[string][]byte{
			"tls.crt": cert,
		},
	}

	// create the secret only if it doesn't exist, else return error that it already exists
	_, err = client.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func GeneratePrivateKeyAndCert(keySize int, validFor time.Duration, cn string) ([]byte, error) {
	r := rand.Reader
	privKey, err := rsa.GenerateKey(r, keySize)
	if err != nil {
		return nil, err
	}
	serialNo, err := rand.Int(r, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}

	cert := x509.Certificate{
		SerialNumber: serialNo,
		KeyUsage:     x509.KeyUsageEncipherOnly,
		NotBefore:    time.Now().UTC(),
		NotAfter:     time.Now().Add(validFor).UTC(),
		Issuer: pkix.Name{
			CommonName: cn,
		},
		Subject: pkix.Name{
			CommonName: cn,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	data, err := x509.CreateCertificate(r, &cert, &cert, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	_ = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: data,
	})

	return certPEM.Bytes(), nil
}
