# nms-proxy
   
第一次创建
生成proto文件
在pb目录下
切换到root
配置好gopath和path中的gopath/bin
protoc -I/usr/local/include -I.   -I$GOPATH/src   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis   --go_out=plugins=grpc:. netServer.proto
protoc -I/usr/local/include -I.   -I$GOPATH/src   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. netServer.proto
在server目录下
go install




在client目录下
go build


server &

gateway &

./client -R -appname="app1" -portname="portsend1" -rioid=11 -slotsize=1024 -remoteappname="app2" -remoteportname="portrecv1" -porttype=0 &


./client -R -appname="app2" -portname="portrecv1" -rioid=11 -slotsize=1024 -remoteappname="app1" -remoteportname="portsend1" -porttype=1 &


./client -R -appname="app1" -portname="portsend2" -rioid=11 -slotsize=1024 -remoteappname="app2" -remoteportname="portrecv2" -porttype=0 &


./client -R -appname="app2" -portname="portrecv2" -rioid=11 -slotsize=1024 -remoteappname="app1" -remoteportname="portsend2" -porttype=1 &

 
curl -X POST -k http://localhost:8080/links/echo 


