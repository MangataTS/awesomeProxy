# AwesomeProxy


## 预计功能

### 正向代理
- [x] HTTP、HTTPS、WS、WSS、SOCKET5多种协议数据包抓包代理

- [x] 支持数据拦截和自定义修改

- [x] 自动识别入站协议，按照消息头识别不同的协议进行处理

- [ ] 正则匹配黑名单HOST过滤

- [ ] 用AC自动机去做一个模式串匹配，处理数据包过滤+网站分析

- [ ] 自动安装证书，代理加密数据包

### 反向代理

- [x] 多种算法实现负载均衡

- [ ] 过滤爬虫、流量攻击

- [ ] IP封禁

### 其他功能

- [ ] 一个独立的Golang日志库，支持彩色打印，分为不同日志级别，线程安全，可以参考：[https://github.com/donnie4w/go-logger](https://github.com/donnie4w/go-logger) 和 [https://github.com/davyxu/golog](https://github.com/davyxu/golog)

- [ ] 生成流量分析报告
  - [ ] 反向代理分析恶意攻击和爬虫
  - [ ] 正向代理分析对用户上网行为分析，违禁网站访问等
  
## config.json

先给一个简单的配置例子：

```json
{
  "ProxyMethod":true,
  "ReProxy":{
    "port":"9090",
    "BalanceMethod":1,
    "backend":[{"host":"127.0.0.1:10001","Weghit":1},{"host":"127.0.0.1:10002","Weghit":1},{"host":"127.0.0.1:10003","Weghit":1}]
  },
  "CoProxy":{
    "port":"9090",
    "MultiListenNum":5,
    "nagle":true,
    "filt":["http://*.csdn.*"]
  }
}
```

- `ProxyMethod` 表示的代理方式，`true` 为反向代理， `false` 为正向代理
- `ReProxy` 表示反向代理的配置
  - `prot` 表示的是反向代理的代理端口
  - `BalanceMethod` 表示的是负载均衡的方法，本项目支持六种负载均衡算法，分别是：
    1. 基于随机算法的负载均衡
    2. 基于RoundRobin算法的负载均衡
    3. 基于带权重RoundRobin算法的负载均衡
    4. 基于一致性hash算法的负载均衡
    5. 基于洗牌算法的负载均衡
    6. 基于优化洗牌算法的负载均衡
  - 项目默认基于RoundRobin算法的负载均衡，如果选用 `基于带权重RoundRobin算法的负载均衡`，请配置 `config.json` 中的 `ReProxy.backend.Weghit` 部分，另外如果只有一台服务器，默认 `RoundRobin`算法，可以提升效率
  - `backend` 表示的是服务器集群，在这里可以配置多个服务器进行负载均衡
    -  `host` 表示的是被代理服务器的ip+端口，或者是url
    -  `weight` 表示的是被代理服务器的权重优先级，范围是[1,10]，值越小权重越大
- CoProxy
  - `port` 表示的是正向代理的代理端口
  - `MultiListenNum` 表示的多线程数，可以设置为CPU的核心数
  - `nagle` 表示是否开启nagle算法优化网络传输，在一些对时延要求较高的交互式操作环境中可以设置 `false` ，默认开启
  - `filt` 表示的是正则过滤的网站，可以将一些违禁网站过滤