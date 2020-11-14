package tool

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

func ParseCertsPEM(pemCerts []byte) ([]*x509.Certificate, error) {
	ok := false
	var certs []*x509.Certificate
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return certs, err
		}

		certs = append(certs, cert)
		ok = true
	}

	if !ok {
		return certs, errors.New("invalid cert pem")
	}
	return certs, nil
}

func ParsePrivKey(keyData []byte) (*rsa.PrivateKey, error) {
	var privateKeyPemBlock *pem.Block
	for {
		privateKeyPemBlock, keyData = pem.Decode(keyData)
		if privateKeyPemBlock == nil {
			break
		}

		switch privateKeyPemBlock.Type {
		case "RSA PRIVATE KEY":
			if key, err := x509.ParsePKCS1PrivateKey(privateKeyPemBlock.Bytes); err == nil {
				return key, nil
			}
		}
	}

	return nil, errors.New("invalid key data")
}

func GenCertAndPrivkey(caPrivKey *rsa.PrivateKey, caCert *x509.Certificate, domains []string) (privKey []byte, cert []byte, err error) {
	// 签证书
	r, err := rand.Int(rand.Reader, big.NewInt(2020))
	certTpl := &x509.Certificate{
		SerialNumber: big.NewInt(r.Int64()),
		Subject: pkix.Name{
			Country:    []string{"CN"},
			Locality:   []string{"BJ"},
			Province:   []string{"BJ"},
			CommonName: domains[0],
		},
		DNSNames:     domains,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return []byte{}, []byte{}, errors.New("create private key failed")
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, certTpl, caCert, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return []byte{}, []byte{}, errors.New("create cert failed")
	}

	keyBuffer := bytes.Buffer{}
	err = pem.Encode(&keyBuffer, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})
	if err != nil {
		return []byte{}, []byte{}, errors.New("encode private key failed")
	}
	certBuffer := bytes.Buffer{}
	err = pem.Encode(&certBuffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return []byte{}, []byte{}, errors.New("encode cert failed")
	}

	return keyBuffer.Bytes(), certBuffer.Bytes(), nil
}
