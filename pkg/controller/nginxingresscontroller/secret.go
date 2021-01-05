package nginxingresscontroller

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	k8sv1alpha1 "github.com/nginxinc/nginx-ingress-operator/pkg/apis/k8s/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const sslHost = "example.com"

func defaultSecretForNginxIngressController(instance *k8sv1alpha1.NginxIngressController) (*corev1.Secret, error) {
	crt, key, err := generateFakeCertAndKey()

	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Data: map[string][]byte{
			corev1.TLSCertKey:       crt,
			corev1.TLSPrivateKeyKey: key,
		},
		Type: corev1.SecretTypeTLS,
	}

	return secret, nil
}

func generateFakeCertAndKey() ([]byte, []byte, error) {
	var priv interface{}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	// This certificate is valid for 365 days
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)

	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"NGINX Inc"},
			CommonName:   "example.com",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{sslHost},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.(*rsa.PrivateKey).PublicKey, priv)
	if err != nil {
		return nil, nil, err
	}

	cert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv.(*rsa.PrivateKey))})

	if cert == nil || key == nil {
		return nil, nil, fmt.Errorf("error encoding ket/crt to PEM")
	}

	return cert, key, nil
}
