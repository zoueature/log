# log

日志包


## 钉钉的配置信息从哪来的??

[怎么创建一个钉钉机器人](https://open.dingtalk.com/document/robots/custom-robot-access)
1. 首先添加一个响应告警群的机器人
 
 > 群设置 --> 机器人管理 --> 添加机器人

2. 设置机器人:
  1. 开启消息推送, 得到一个Webhook地址, 地址上面有access_token参数, 获取此参数得到access_token
  2. 安全设置修改为加签, 得到一个签名密钥串, 此参数为signSecret

3. 用这两个参数初始化log包
