# 基于kafka封装的消息队列, changed from go-common(b~i~l~i~b~i~l~i)
## 说明~仅供学习和借鉴
* 服务启动后，会监听配置的端口，等待生产或消费client连接;
* client连接上后用redis协议认证auth
* 认证成功后，通过认证串里包含的信息来连接kafka集群（使用库github.com/Shopify/sarama），使用对应的topic和groupname
* 根据role是pub还是sub来执行生产或消费操作
* 生产或消费 使用的redis协议中的set和mget命令

## 时序图
```sequence
title:databus数据时序图

participant producerBiz
participant consumerBiz
participant dataBus
participant kafka

producerBiz->dataBus:Auth
dataBus-->>producerBiz:Auth Success
Note over producerBiz:send Msg
dataBus->kafka:produce Msg

consumerBiz->dataBus:Auth
dataBus-->>consumerBiz:Auth Success
Note over consumerBiz:pull Msg
dataBus->kafka:consume Msg
```

## 环境
### 认证 可选redis或者mysql
#### redis(databus连接认证所用) 使用redis_orm
* redis 建表所用模型（请见 https://github.com/weikaishio/redis_orm_workbench 建表和手动录入都可以）
```go
type AuthTb struct {
	Id         int64  `redis_orm:"pk autoincr comment 'ID'"`
	AppId      int32  `redis_orm:"dft '' comment 'AppId'"`
	Group      string `redis_orm:"index dft '' comment '组名'"`
	Operation  int8   `redis_orm:"dft '0' comment '操作类型'"`
	Topic      string `redis_orm:"dft '' comment '主题名'"`
	CreatedAt  int64  `redis_orm:"created_at comment '创建时间'"`
	UpdatedAt  int64  `redis_orm:"updated_at comment '更新时间'"`
}
type AppTb struct {
	Id        int64  `redis_orm:"pk autoincr comment 'ID'"`
	AppKey    string `redis_orm:"dft '' comment 'key'"`
	AppSecret string `redis_orm:"dft '' comment 'secret'"`
	Cluster   string `redis_orm:"dft '' comment '集群名'"`
	CreatedAt int64  `redis_orm:"created_at comment '创建时间'"`
	UpdatedAt int64  `redis_orm:"updated_at comment '更新时间'"`
}
```
#### mysql(databus连接认证所用)，认证有进程内缓存，每分钟拉数据库数据更新一次，最多一分钟脏数据
* mysql 建表
```sql
create database databus_db;
use databus_db;
CREATE TABLE `auth` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_id` INT(11) NOT NULL,
  `group_name` VARCHAR(45) NULL,
  `operation`TINYINT(4) DEFAULT 0,
  `topic` VARCHAR(45) NOT NULL,
  `created_at` BIGINT(20) NULL,
  `deleted_at` BIGINT(20) NULL,
  `updated_at` BIGINT(20) NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `app` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_key` VARCHAR(45) NOT NULL,
  `app_secret` VARCHAR(64) NOT  NULL,
  `cluster` VARCHAR(45) NOT NULL,
  `created_at` BIGINT(20) NULL,
  `deleted_at` BIGINT(20) NULL,
  `updated_at` BIGINT(20) NULL,
  PRIMARY KEY (`id`));
```
* 录入app认证数据
```sql
INSERT INTO `app` (`id`,`app_key`,`app_secret`,`cluster`) VALUES (1,'app_key1','app_secret1','cluster1');
```
* 录入group相关数据
```sql
INSERT INTO `auth` (`id`,`app_id`,`group_name`,`operation`,`topic`) VALUES (1,1,'group_name1',3,'test156');
```
### kafka集群 
* 无特别处理

## 使用
### 启动 example_svr/cmd 服务
* 注意目录下的toml文件，需要配置上cluster、addr、mysql(or redis)配置节不能少
* 需要在example_svr/cmd目录运行go run main.go测试，否则运行需要带上-conf=配置文件地址

### 用example_cli测试
* example_cli目录有测试配置，配上认证数据即可，需要在example_cli目录运行go run main.go测试，否则运行需要带上-conf=配置文件地址
* 运行后会创建生产者，并且发布一个消息；然后创建消费者并且一直消费

### 可以用redis-cli测试
* 连接(比如在本地运行) e.g. redis-cli -h 127.0.0.1 -p 6205
* 用生产者身份认证 auth命令 e.g. auth app_key1:app_secret1@group_name1/topic=test156&role=pub
* 测试连通性：ping命令
* 生产消息：set命令 e.g. set x "xx"  注：第一个参数x是key,(kafka多partition时非常有用)，用引号的原因是json反序列化所用
* 用消费者身份认证：auth命令 e.g. app_key1:app_secret1@group_name1/topic=test156&role=sub
* 消费消息 mget命令 e.g. mget pb  注：如果没有消息会block一段时间，pb意思是返回protobuf序列化的数据
* 消费者执行set命令则是提交offset e.g. set 1 1 表示针对partion=1的设置offset=1

#### 附：redis 协议说明
Redis服务器与客户端通过RESP（REdis Serialization Protocol）协议通信。RESP协议支持的数据类型：
* Simple String第一个字节以+开头，随后紧跟内容字符串（不能包含CR LF），最后以CRLF结束。很多Redis命令执行成功时会返回"OK"，"OK"就是一个Simple String：
"+OK\\r\\n"

* Error的结构与Simple String很像，但是第一个字节以-开头：
"-ERR unknown command 'foobar'"

* -符号后的第一个单词代表错误类型，ERR代表一般错误，WRONGTYPE代表在某种数据结构上执行了不支持的操作。

* Integer第一个字节以:开头，随后紧跟数字，以CRLF结束：
":1000\\r\\n"
很多Redis命令会返回Integer,例如INCR LLEN等。

* Bulk String是一种二进制安全的字符串结构，整个结构包含两部分。第一部分以$开头，后面紧跟字符串的字节长度，CRLF结尾。第二部分是真正的字符串内容，CRLF结尾，最大长度限制为512MB。一个Bulk String结构的"Hello World!"是：
"$12\\r\\nHello World!\\r\\n"

* 空字符串是：
"$0\\r\\n\\r\\n"  

* nil是：
"$-1\\r\\n"

* Array也可以看成由两部分组成，第一部分以*开头，后面紧跟一个数字代表Array的长度，以CRLF结束。第二部分是每个元素的具体值，可能是Integer，可能是Bulk String。Array结构的["hello", "world"]是：
"*2\\r\\n$5\\r\\nhello\\r\\n$5\\r\\nworld\\r\\n"

* 空Array：
"*-1\\r\\n" 