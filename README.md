# MixProvider
MixProvider 是按照Open Falcon监控系统短信报警接口规范开发的Provider程序。

## 主要功能：
* HTTP Server 接收请求
* 发送告警消息到企业微信接口
* 发送告警短信到容联云通讯接口
* 发送告警电话到容联云通讯接口
* 低优先级报警收敛

## 文件说明:
* config.json              -- 应用配置文件
* config/config.go         -- 配置处理包
* control                  -- 控制脚本
* delayer/delayer.go       -- 报警延迟处理包
* email/email.go           -- 邮件处理包
* filter/filter.go         -- 过滤器包
* main.go                  -- 程序入口
* server/server.go         -- HTTP Server包
* utils/message.go         -- 报警消息
* utils/set.go             -- Map实现的Set
* wechat/wechat.go         -- 微信处理包
* wechat.yaml              -- falcon/微信用户映射配置
* yuntongxun/yuntongxun.go -- 容联云通讯处理包
* yuntongxun.yaml          -- 容联云通讯媒体文件配置

## 使用方法
1. git clone https://github.com/op-y/MixProvider.git mix-provider
2. cd mix-provider
3. go build -o mix-provider
4. 根据实际情况修改config.json/wechat.yaml/yuntongxun.yaml配置
5. ./control start

## 修改历史
* 2017/9/27 根据业务需要，将短信和微信接口分离，去掉了邮件接口
* 2018/5/8 添加了电话(语言)报警接口，微信接口修改为alarm接口根据不同优先级执行不同操作，恢复邮件接口
