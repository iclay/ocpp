## description
It should be noted that the plug-in is required and must be included in your code. Similarly, you can customize the plug-in on the premise that you must implement the following methods,If you implement the action method, the center service will automatically discover and execute

//RequestHandler represent device active request Center  
//action - ocpp protocol action  
//proto.RequestHandler - callback method about action  
**RequestHandler(action string) (proto.RequestHandler, bool)**  

//ResponseHandler represent The device reply to the center request  
//action - ocpp protocol action  
//proto.ResponseHandler - callback method about action  
**ResponseHandler(action string) (proto.ResponseHandler, bool)**  
## support
- [x] current implementation
  - [x] LocalService
  - [x] RPCX