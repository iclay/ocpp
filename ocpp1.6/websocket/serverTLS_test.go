package websocket

func createTLSCertificate(certificateFilename string, keyFilename string, cn string, ca *x509.Certificate, caKey *ecdsa.PrivateKey) error {
	// Generate ed25519 key-pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
	// Create self-signed certificate
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Hour * 24)
	keyUsage := x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"ocpp-go"},
			CommonName:   cn,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           extKeyUsage,
		BasicConstraintsValid: true,
		DNSNames:              []string{cn},
	}
	var derBytes []byte
	if ca != nil && caKey != nil {
		// Certificate signed by CA
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, ca, &privateKey.PublicKey, caKey)
	} else {
		// Self-signed certificate
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	}
	if err != nil {
		return err
	}
	// Save certificate to disk
	certOut, err := os.Create(certificateFilename)
	if err != nil {
		return err
	}
	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}
	// Save key to disk
	keyOut, err := os.Create(keyFilename)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	privateBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateBytes})
	if err != nil {
		return err
	}
	return nil
}
