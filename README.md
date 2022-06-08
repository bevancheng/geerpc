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

## Day2 

一个函数需要能够被远程调用，需要满足条件：

Method's type is exported
Method is exported
Method has two args, both exported types
the second arg is a pointer
method has return type error

func (t *T) MethodName(argType T1, replyType *T2) error
## Day3

![geerpcday3](https://raw.githubusercontent.com/bevancheng/imgrepo/main/202206081259912.jpg)

## Day4
增加连接超时的处理机制
增加服务端处理超时的处理机制

超时处理是rpc框架一个比较基本的能力，如果缺少超时处理，服务端和客户端都容易一因为网络或其他其他错误导致挂死，降低服务的可用性。


需要客户端处理超时的地方：
- 与服务端**建立连接**，导致的超时
- 发送请求到服务端，**写报文**导致的超时
- 等待服务端处理，等待**处理报文**导致的超时（服务端不响应）
- 从服务端接收响应，**读报文**导致的超时

需要服务端处理的：
- 读取客户端请求报文时，**读报文**导致的超时
- 发送响应报文时，**写报文**导致的超时
- 调用映射服务的方法时，**处理报文**导致的超时

geerpc实现：
- 客户端创建连接时
- 客户端Client.Call()整个过程的超时（发送报文，等待处理，接收报文所有阶段）
- 服务端处理报文，Server.handRequest超时

实现的代码中client.go的dialTimeout和server.go的handleRequest有内存泄露问题

## Day6
随机选择和Round Robin轮询调度算法实现负载均衡

## Day7
注册中心，支持服务注册、接收心跳
客户端基于注册中心的服务发现