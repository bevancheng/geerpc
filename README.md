#
学习 https://geektutu.com/post/geerpc.html

## RPC框架
RPC(Remote Procedure Call)是一种计算机通信协议，允许调用不同进程空间的程序。RPC的客户端和服务器可以在一台机器上也可以不在。程序员使用时，可以无需关注内部实现细节。
浏览器和服务器间广泛使用基于HTTP协议的Restful API。Restful API有相对统一的标准，通用性更好，兼容性好。
缺点也很明显：
- Restful API需要额外定义，无论客户端还是服务端，都需要额外代码，RPC调用接近直接调用。
- 基于HTTP的Restful报文冗杂，RPC通常用自定义协议格式。
- RPC用更高效的序列化协议
- RPC更容易扩展和集成注册中心、负载均衡等功能。

需要解决的问题：传输协议（TCP、HTTP、Unix Socket），报文编码（JSON、XML、Protobuf），可用性问题（连接超时、异步请求、并发）

服务端实例很多的话，客户端不关心这些实例的地址和部署位置，只关心自己能否得到期待的结果，这时就要用到注册中心（registry）和负载均衡（load balance）。服务端启动将自己注册到注册中心，客户端调用时，从注册中心获取可用实例，选择调用。注册中心需要实现服务动态添加、删除，使用心跳保证服务处于可用状态。

成熟的RPC框架：grpc、rpcx、go-micro。
rpc是微服务框架的一个子集，微服务框架可以自己实现rpc部分，也可以选择不同的rpc框架作为通信基座。
第三方库：protobuf、etcd、zookeeper。
geerpc：protocol exchange、registry、service discovery、load balance、timeout processing

## Day1 服务端与消息编码


通信过程
报文最开始规定固定的字节，来协商相关的信息。比如第1个字节表示序列化方式，第二个字节表示压缩方式，3-6表示header长度，7-10表示body长度

为了实现简单，协商阶段还是用JSON编码OPTION，根据OPTION字段的CodecType获取接下来的编码方式
| Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
| <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
一次连接可能是：
| Option | Header1 | Body1 | Header2 | Body2 | ... 