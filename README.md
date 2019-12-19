# databus_kafka, changed from go-common(b~i~l~i~b~i~l~i)
## 说明
* 服务启动后，会监听配置的端口，等待client连接;
* client连接上后用redis协议认证auth
* 认证成功后，通过认证串里包含的role是pub还是sub来执行不同的操作（生产或消费）
 `
## 环境
### mysql(databus连接认证所用)
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
  PRIMARY KEY (`id`));

CREATE TABLE `app` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `app_key` VARCHAR(45) NOT NULL,
  `app_secret` VARCHAR(64) NOT  NULL,
  `cluster` VARCHAR(45) NOT NULL,
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
* 无特别处理，就自己搭建即可

## 使用
### 启动 databus_kafka/cmd 服务
* 注意目录下的toml文件，需要配置上cluster、addr、mysql配置节不能少
### 用redis-cli测试
* 连接(比如在本地运行):redis-cli -h 127.0.0.1 -p 6205
* 用生产者身份认证：auth app_key1:app_secret1@group_name1/topic=test156&role=pub
* 测试连通性：ping
* 生产消息：set 1 "xx"  - 用引号的原因是json反序列化所用
* 用消费者身份认证：auth app_key1:app_secret1@group_name1/topic=test156&role=sub
* 消费消息 mget pb  - 如果没有消息会block一段时间，pb意思是返回protobuf序列化的数据
* 注意：消费者执行set命令则是提交offset
### 用example_cli测试
* example_cli目录有测试配置，配上认证数据即可
* 运行后会创建生产者，并且发布一个消息；然后创建消费者并且一直消费