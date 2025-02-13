<div align="center">
    <h1>R106随身WiFi功能增强</h1>
    <img alt="Golang" src="https://img.shields.io/badge/-Go-00ADD8?style=flat-square&logo=go&logoColor=ffffff">
    <img  src="https://img.shields.io/badge/R106%E9%9A%8F%E8%BA%ABWiFi-8A2BE2">
    <img src="https://img.shields.io/badge/license-MPL2.0-green">
    <a href="#下载"><img alt="下载最新版" src="https://img.shields.io/badge/Download-Latest%20Releases-blue"></a>
</div>


## 功能

- [x] 短信实时转发钉钉通知
- [x] 设备低电量通知提醒
- [x] 搭配frp可通过手机远程控制WiFi
- [x] 远程查看设备当前状态（内存使用情况、WiFi状态、设备电量、充电状态）

## 效果

### 钉钉通知

![image-20250213214525566](assets/image-20250213214525566.png)

### 远程查看与控制

![image-20250213214637408](assets/image-20250213214637408.png)

### 程序大小

![image-20250213205723305](assets/image-20250213205723305.png)

## 使用

1. 下载修改`config.yaml`配置文件

   ![image-20250213220004054](assets/image-20250213220004054.png)

2. 开机自启：在`/etc/init.d/hostname.sh`里添加启动`./script/r106-init.sh`脚本的命令 例如：`sh /home/root/r106-enhance/script/r106-init.sh &`

3. 重启设备即可。如果不想重启也可以直接执行`./script/r106-init.sh`脚本


## `下载`

- [蓝奏云下载](https://c-xzhi.lanzouu.com/b0mbd0wze) 密码：Xzhi

- [GitHub下载](https://github.com/GitHub-Xzhi/go-r106-enhance/releases) 

## 💖支持

 **如果这个程序对您有帮助🤓，请给个Star，也可以请作者吃辣条或喝咖啡😋**

| **微信** | **微信赞赏** | **支付宝** |
| :---: | :---: | :---: |
| ![微信收款码](assets/wx_fkm.png)|![微信赞赏码](assets/wxzsm.png)| ![支付宝收款码](assets/zfb_fkm.png) |

