package server

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"math/big"
	randn "math/rand"
	"net/http"
	"net/url"
	local "ocpp16/plugin/passive/local"
	"ocpp16/protocol"
	"os"
	"sync"
	"testing"
	"time"
)

var wss_addr = flag.String("wss_addr", "localhost:8091", "websocket tls service address")

func createTLSCertificate(certificateFilename string, keyFilename string, cn string, ca *x509.Certificate, caKey *ecdsa.PrivateKey) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return err
	}
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
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, ca, &privateKey.PublicKey, caKey)
	} else {
		derBytes, err = x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	}
	if err != nil {
		return err
	}
	certOut, err := os.Create(certificateFilename)
	if err != nil {
		return err
	}
	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}
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

func TLSClientHandler(ctx context.Context, t *testing.T, d *dispatcher, serverCertName string, clientCertName string, clientKeyName string) {
	flag.Parse()
	name, id := RandString(5), RandString(5)
	// name, id := "qinglianyun", "sujunkang"
	path := fmt.Sprintf("/ocpp/%s/%s", name, id)
	u := url.URL{Scheme: "wss", Host: "localhost:8091", Path: path}
	certPool := x509.NewCertPool()
	data, err := ioutil.ReadFile(serverCertName)
	require.Nil(t, err)
	ok := certPool.AppendCertsFromPEM(data)
	require.True(t, ok)
	loadedCert, err := tls.LoadX509KeyPair(clientCertName, clientKeyName)
	require.Nil(t, err)
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		Subprotocols:     []string{"ocpp1.5", "ocpp1.6"},
		TLSClientConfig: &tls.Config{
			RootCAs:      certPool,
			Certificates: []tls.Certificate{loadedCert},
		},
	}
	c, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		t.Fatal("dial:", err)
	}
	defer c.Close()
	ch := make(chan string, 10)
	var closed bool
	defer func() {
		closed = true
		// close(ch)
	}()
	queue := NewRequestQueue()
	var waitgroup sync.WaitGroup
	var mtx sync.Mutex
	waitgroup.Add(1)
	go func() {
		for range time.Tick(time.Second * 10) {
			mtx.Lock()
			if err = c.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
				t.Error(err)
			}
			mtx.Unlock()
		}
	}()
	go func() { //test for center request
		defer waitgroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				call := &protocol.Call{
					MessageTypeID: protocol.CALL,
					UniqueID:      RandString(7),
					Action:        "BootNotification",
					Request:       fnBootNotificationRequest(),
				}
				queue.Push(call.UniqueID)
				if err := d.appendRequest(context.Background(), fmt.Sprintf("%s-%s", name, id), call); err != nil {
					return
				}
				// time.Sleep(time.Second*1)
				time.Sleep(time.Second * time.Duration(randn.Intn(5)) / 5)
			}
		}
	}()
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case res_uniqueid := <-ch:
				rep_uniqueid, _ := queue.Pop()
				next_uniqueid, _ := queue.Peek()
				t.Logf("ws_id(%s), res_uniqueid(%s),rep_uniqueid(%s),queue remain(%d), next_uniqueid(%v)", fmt.Sprintf("%s-%s", name, id), res_uniqueid, rep_uniqueid, queue.Len(), next_uniqueid)
				if res_uniqueid != rep_uniqueid {
					t.Errorf("ws_id(%s), res_uniqueid(%s) != rep_uniqueid(%s)", fmt.Sprintf("%s-%s", name, id), res_uniqueid, rep_uniqueid)
				}
			}
		}
	}()
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					return
				}
				fields, err := parseMessage(message)
				if err != nil {
					return
				}
				switch fields[0].(float64) {
				case float64(protocol.CALL):
					go func() {
						Interval := 10
						uniqueid := fields[1].(string)
						callResult := &protocol.CallResult{
							MessageTypeID: protocol.CALL_RESULT,
							UniqueID:      uniqueid,
							Response: &protocol.BootNotificationResponse{
								CurrentTime: time.Now().Format(time.RFC3339),
								Interval:    &Interval,
								Status:      "Accepted",
							},
						}
						callResultMsg, err := json.Marshal(callResult)
						if err != nil {
							return
						}
						// time.Sleep(time.Second * time.Duration(randn.Intn(3)) / 10)
						t.Logf("test for center call: recv msg(%+v), resp_msg(%+v)", string(message), string(callResultMsg))
						mtx.Lock()
						err = c.WriteMessage(websocket.TextMessage, callResultMsg)
						mtx.Unlock()
						if err != nil {
							return
						}
						if !closed {
							ch <- callResult.UniqueID
						}
					}()
				case float64(protocol.CALL_RESULT), float64(protocol.CALL_ERROR):
					t.Logf("test for client call: recv msg(%s), ", string(message))
				default:
					t.Log(string(message))
				}

			}
		}
	}()
	//test for client call
	waitgroup.Add(1)
	go func() {
		defer waitgroup.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				var action = "BootNotification"
				call := &protocol.Call{
					MessageTypeID: protocol.CALL,
					UniqueID:      RandString(7),
					Action:        action,
				}
				switch action {
				case "StatusNotification":
					call.Request = fnStatusNotificationRequest()
				case "Authorize":
					call.Request = fnAuthorizeRequest()
				case "BootNotification":
					call.Request = fnBootNotificationRequest()
				case "MeterValues":
					call.Request = fnMeterValueRequest()
					t.Logf("%+v", call.Request)
				case "StartTransaction":
					call.Request = fnStartTransactionRequest()
				case "StopTransaction":
					call.Request = fnStopTransactionRequest()
				default:
				}
				callMsg, err := json.Marshal(call)
				if err != nil {
					t.Error(err)
					return
				}
				mtx.Lock()
				err = c.WriteMessage(websocket.TextMessage, callMsg)
				mtx.Unlock()
				if err != nil {
					t.Error(err)
					return
				}
				time.Sleep(time.Second * 100)
			}
		}
	}()
	waitgroup.Wait()
	t.Logf("(%s) grace exit gorutine", path)
}

func WssHandler(t *testing.T, waitGroup *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*40)
	clientCertName := "../cert/client/client.pem"
	clientKeyName := "../cert/client/client_key.pem"
	err := createTLSCertificate(clientCertName, clientKeyName, "localhost", nil, nil)
	require.Nil(t, err)
	defer os.Remove(clientCertName)
	defer os.Remove(clientKeyName)
	serverCertName := "../cert/server/cert.pem"
	serverKeyName := "../cert/server/key.pem"
	err = createTLSCertificate(serverCertName, serverKeyName, "localhost", nil, nil)
	require.Nil(t, err)
	defer os.Remove(serverCertName)
	defer os.Remove(serverKeyName)
	lg := initLogger()
	SetLogger(lg)
	server := NewDefaultServer()
	plugin := local.NewActionPlugin()
	server.RegisterActionPlugin(plugin)
	go func() {
		server.ServeTLS(*wss_addr, "/ocpp/:name/:id", serverCertName, serverKeyName)
	}()
	for i := 0; i < 2; i++ { //numbers of client
		time.Sleep(time.Second / 10)
		go func() {
			TLSClientHandler(ctx, t, server.dispatcher, serverCertName, clientCertName, clientKeyName)
		}()
	}
	select {
	case <-ctx.Done():
		time.Sleep(time.Second * 10)
		cancel()
	}
	waitGroup.Done()
}
