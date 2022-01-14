## description
The plug-in acts on the active request device side of the service center. It must be included in your code. When a new request is routed to the plug-in. The plug-in will send the request to the request queue of the service center. You must implement the following methods


//NewActiveCallPlugin represent you will receive a callback function for the push request,You only need to implement the following method, which will be automatically found and called by the service center  
//websocket.ActiveCallHandler - The active request callback function provided by the service center. When the plug-in receives a new request, you only need to execute the callback function  
**func NewActiveCallPlugin(handler websocket.ActiveCallHandler)**  


## support
- [x] current implementation
  - [x] LocalService  
  - [x] RPCX

