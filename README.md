# rtools
RocketMQ 命令行测试工具

* 消费主题数据
```sh
rtools  consumer -n 127.0.0.1:9876  --topic test  --group taaaa
```
* 向topic发送消息
```sh
rtools  producer  -n 172.20.132.1:9876  --topic test  --group taaaa
```