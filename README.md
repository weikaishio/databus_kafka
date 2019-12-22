# 基于kafka封装的消息队列, changed from go-common(b~i~l~i~b~i~l~i)
## 说明~仅供学习和借鉴
* 服务启动后，会监听配置的端口，等待生产或消费client连接;
* client连接上后用redis协议认证auth
* 认证成功后，通过认证串里包含的信息来连接kafka集群（使用库github.com/Shopify/sarama），使用对应的topic和groupname
* 根据role是pub还是sub来执行生产或消费操作
* 生产或消费 使用的redis协议中的set和mget命令

## 环境
### mysql(databus连接认证所用)，认证有进程内缓存，每分钟拉数据库数据更新一次，最多一分钟脏数据
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
* 注意目录下的toml文件，需要配置上cluster、addr、mysql配置节不能少
* 需要在example_svr/cmd目录运行go run main.go测试，否则运行需要带上-conf=配置文件地址

### 可以用redis-cli测试
* 连接(比如在本地运行) e.g. redis-cli -h 127.0.0.1 -p 6205
* 用生产者身份认证 auth命令 e.g. auth app_key1:app_secret1@group_name1/topic=test156&role=pub
* 测试连通性：ping命令
* 生产消息：set命令 e.g. set x "xx"  注：第一个参数x是key,(kafka多partition比较有用)，用引号的原因是json反序列化所用
* 用消费者身份认证：auth命令 e.g. app_key1:app_secret1@group_name1/topic=test156&role=sub
* 消费消息 mget命令 e.g. mget pb  注：如果没有消息会block一段时间，pb意思是返回protobuf序列化的数据
* 消费者执行set命令则是提交offset e.g. set 1 1 表示针对partion=1的设置offset=1

### 用example_cli测试
* example_cli目录有测试配置，配上认证数据即可，需要在example_cli目录运行go run main.go测试，否则运行需要带上-conf=配置文件地址
* 运行后会创建生产者，并且发布一个消息；然后创建消费者并且一直消费