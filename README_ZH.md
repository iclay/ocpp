# 简介

实现的开放充电协议基于websocket底层信框架,通信协议支持ws和wss,支持业务端以自定义插件的形式集成到整个通信系统服务中，该协议实现了基于go原生goroutine和reactor+epoll的两个可选版本,目前仅支持ocpp1.6的所有功能   

当前支持:

- [x] OCPP 1.6
- [ ] OCPP 2.0 

### 自定义插件使用说明
如果要将自定义功能插件集成到通信服务中,必须要在plugin目录下实现接口定义的回调函数，plugin目录下包含两个子目录active以及passive   

- **active**: 该目录下需要实现充电系统主动下发数据到充电桩桩的自定义插件，当前已经支持本地和rpcx插件,自定义插件使用介绍如下,使用可以详见 active/README.md  
  ```go
    //该参数是充电系统提供的一个闭包回调函数,插件只需调用这个回调函数便可下发命令到充电桩，可以参考(active/local/plugin.go)的实现方式
    func NewActiveCallPlugin(handler websocket.ActiveCallHandler)
  ```

- **passive**: 该目录下需要实现充电系统接收充电桩主动请求后自定义功能插件，当前已经支持本地和rpcx插件,自定义插件使用详见 passive/README.md      
    ```go
    //插件必须实现ActionPlugin接口，可以参考(passive/local/plugin.go)的实现方式
    type ActionPlugin interface {
        //传入参数action:同ocpp1.6协议约定action
        //第一个返回值返回的是一个关于action自定义回调函数(该函数是充电桩请求充电系统对应action的回调函数)
        //第二个返回值是插件是否该action，如果不支持，则充电系统会返回桩一个ocpp1.6中约定的错误信息
	    RequestHandler(action string) (protocol.RequestHandler, bool) 
        //传入参数action:同ocpp1.6协议约定action,
        //第一个返回值返回的是一个关于action自定义回调函数(该函数是充电桩应答充电系统主动下发命令的自定义回调函数)
        //第二个返回值是插件是否支持支持该action，如果不支持，则充电系统一个ocpp1.6中约定的错误信息
	    ResponseHandler(action string) (protocol.ResponseHandler, bool) 
        }  
    ```   

### 插件集成到充电系统使用示例
```go
import (
	"fmt"
	"ocpp16/config" //配置文件
	"ocpp16/logwriter"
	// active "ocpp16/plugin/active/local" //本地实现的active插件，用于在单机服务中
	// passive "ocpp16/plugin/passive/local" //本地实现的passive插件，用于在单机服务中
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	active "ocpp16/plugin/active/rpcx" //rpcx实现的active插件，用于在分布式服务中
	passive "ocpp16/plugin/passive/rpcx" //rpcx实现的passive插件，用于在分布式服务中
	"ocpp16/websocket" //websocket底层通信库
	"os"
	"time"
)
func main() {
        //解析自定义配置文件
	config.ParseFile(c.String("config")) 
	config.Print()
	conf := config.GCONF
	lg := initLogger()
	websocket.SetLogger(lg)
      //启动一个默认充电系统服务
	server := websocket.NewDefaultServer() 
	defer server.Stop() 
        //自定义passive插件，当前使用的是rpcx插件
	actionPlugin := passive.NewActionPlugin() 
        //将该passive插件集成到充电系统中，充电系统来代理插件来执行插件内的自定义功能
	server.RegisterActionPlugin(actionPlugin)
        //充电桩连接到充电系统的自定义回调函数
	server.SetConnectHandlers(func(ws *websocket.Wsconn) error { 
		lg.Debugf("id(%s) connect,time(%s)", ws.ID(), time.Now().Format(time.RFC3339))
		return nil
	})
        //充电桩断开连接的自定义回调函数
	server.SetDisconnetHandlers(func(ws *websocket.Wsconn) error { 
		lg.Debugf("id(%s) disconnect,time(%s)", ws.ID(), time.Now().Format(time.RFC3339))
		return nil
	}, func(ws *websocket.Wsconn) error {
		return actionPlugin.ChargingPointOffline(ws.ID())
	})
        //将自定义的active插件集成到充电系统中，当前使用的是rpcx插件,充电系统来代理插件下发命令到充电桩
	server.RegisterActiveCallHandler(server.HandleActiveCall, active.NewActiveCallPlugin) 
	ServiceAddr, ServiceURI := conf.ServiceAddr, conf.ServiceURI
	if conf.WsEnable {
		wsAddr := fmt.Sprintf("%s:%d", ServiceAddr, conf.WsPort)
                //server启动ws服务
		server.Serve(wsAddr, ServiceURI) 
	}
	if conf.WssEnable && conf.TLSCertificate != "" && conf.TLSCertificateKey != "" {
		wssAddr := fmt.Sprintf("%s:%d", ServiceAddr, conf.WssPort)
              //server启动wss服务
		server.ServeTLS(wssAddr, ServiceURI, conf.TLSCertificate, conf.TLSCertificateKey)
	}
	return nil
}
```
