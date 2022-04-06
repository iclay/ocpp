# brief introduction

Open charge point protocol (ocpp) is a global open communication standard. It is a network protocol used for communication between electric vehicle charger and central background system. It is mainly used to solve various difficulties caused by communication between private charging networks. Ocpp has been applied to multiple charging facilities in 49 countries, so it has essentially become the industry standard for charging facilities and network communication, Ocpp supports seamless communication management between the charging station and the central management system of each supplier. The communication protocol supports WS and WSS, and supports the business end to integrate into the whole communication system service in the form of user-defined plug-ins. The protocol realizes two optional versions based on go native goroutine and reactor + epoll, At present, it only supports all functions of ocpp1.6   

Please note that this is a library rather than an application. You can refer to the center_ system. Go to implement your own main method

Current support:

- [x] OCPP 1.6
- [ ] OCPP 2.0 

### User defined plug-in instructions
If you want to integrate the custom function plug-in into the communication service, you must implement the callback function defined by the interface in the plugin directory, which contains two subdirectories active and passive

- **active**: Under this directory, the user-defined plug-in that the charging system actively sends data to the charging pile needs to be realized. At present, local and rpcx plug-ins are supported. The usage of the user-defined plug-in is described below. See the usage for details of active/README.md  
  ```go
    //This parameter is a closure callback function provided by the charging system. The plug-in only needs to call this callback function to issue a command to the charging pile. Please refer to the implementation method of (active/local/plugin.go)
    func NewActiveCallPlugin(handler ocpp16server.ActiveCallHandler)
  ```

- **passive**: Under this directory, the charging system needs to realize the user-defined function plug-in after receiving the active request of the charging point. At present, local and rpcx plug-ins are supported. See passive/readme for the use of user-defined plug-ins md
    ```go
    //The plug-in must implement the actionplugin interface. Please refer to the implementation method of (passive/local/plugin.go)
    type ActionPlugin interface {
        //Pass in parameter action: the same as ocpp1 6 agreement action
        //The first return value returns a user-defined callback function about action (this function is the callback function of the charging point requesting the corresponding action of the charging system)
        //The second return value is whether the plug-in should the action. If it is not supported, the charging system will return an error message agreed in ocpp1.6
	    RequestHandler(action string) (protocol.RequestHandler, bool) 
        //Pass in parameter action: the same as ocpp1.6. Agreement action,
        //The first return value returns a user-defined callback function about action (this function is the user-defined callback function of the charging point responding to the charging system's active command)
        //The second return value is whether the plug-in supports the action. If not, the charging system will send an error message agreed in ocpp1.6
	    ResponseHandler(action string) (protocol.ResponseHandler, bool) 
        }  
    ```   

### Example of plug-in integrated into charging system
```go
import (
	"fmt"
	"ocpp16/config" 
	"ocpp16/logwriter"
        //An active plug-in implemented locally, which is used in stand-alone services
	// active "ocpp16/plugin/active/local" 
        //The passive plug-in implemented locally is used in stand-alone services
	 // passive "ocpp16/plugin/passive/local"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
        //The active plug-in implemented by rpcx is used in distributed services
	active "ocpp16/plugin/active/rpcx"
        //The passive plug-in implemented by rpcx is used in distributed services
	passive "ocpp16/plugin/passive/rpcx" 
        //Bottom communication library of server
	ocpp16server "ocpp16/server" 
	"os"
	"time"
)
func main() {
        //Configuration file
	config.ParseFile(c.String("config")) 
	config.Print()
	conf := config.GCONF
	lg := initLogger()
	ocpp16server.SetLogger(lg)
        //Start a default charging system service
	server := ocpp16server.NewDefaultServer() 
	defer server.Stop() 
        //Customize the passive plug-in. The rpcx plug-in is currently used
	actionPlugin := passive.NewActionPlugin() 
        //Integrate the passive plug-in into the charging system, and the charging system will proxy the plug-in to perform the custom functions in the plug-in
	server.RegisterActionPlugin(actionPlugin)
        //Custom callback function of charging point connected to charging system
	server.SetConnectHandlers(func(ws *ocpp16server.Wsconn) error { 
		lg.Debugf("id(%s) connect,time(%s)", ws.ID(), time.Now().Format(time.RFC3339))
		return nil
	})
        //Custom callback function for charging point disconnection
	server.SetDisconnetHandlers(func(ws *ocpp16server.Wsconn) error { 
		lg.Debugf("id(%s) disconnect,time(%s)", ws.ID(), time.Now().Format(time.RFC3339))
		return nil
	}, func(ws *ocpp16server.Wsconn) error {
		return actionPlugin.ChargingPointOffline(ws.ID())
	})
        //The user-defined active plug-in is integrated into the charging system. Currently, the rpcx plug-in is used. The charging  system sends commands to the charging pile on behalf of the plug-in
	server.RegisterActiveCallHandler(server.HandleActiveCall, active.NewActiveCallPlugin) 
	ServiceAddr, ServiceURI := conf.ServiceAddr, conf.ServiceURI
	if conf.WsEnable {
		wsAddr := fmt.Sprintf("%s:%d", ServiceAddr, conf.WsPort)
                //Server starts ws service
		server.Serve(wsAddr, ServiceURI) 
	}
	if conf.WssEnable && conf.TLSCertificate != "" && conf.TLSCertificateKey != "" {
		wssAddr := fmt.Sprintf("%s:%d", ServiceAddr, conf.WssPort)
                //Server starts the wss service
		server.ServeTLS(wssAddr, ServiceURI, conf.TLSCertificate, conf.TLSCertificateKey)
	}
	return nil
}
```
