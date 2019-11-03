let's encrypt Certbot hook written in golang
===============================

lets'encrypt的通配符https证书官方推荐使用certbot客户端获取。
certbot验证域名是否属于某人，通过产生一个随机字符串，让用户设置一个TXT记录到此值后，当验证成功后颁发证书。
该证书有效期为3个月，意味着每3个月得更新一次，每次手工设置dns解析及其麻烦，还容易忘记更新证书。
好在certbot提供了一个hook，产生随机后调用hook来自动设置dns记录，官方预制了几家dns服务商的hook，
但是国内的阿里云、腾讯云都不在此例。此插件支持国内阿里云、腾讯云等其他常用dns服务商，
使用golang编写，服务器无需搭建运行环境，直接下载二进制分发程序即可。


支持的DNS服务商
-------
- [x] 阿里云

- [x] 腾讯云


DNS服务商密钥获取
---------------
1. 阿里云进入管理后台，悬浮右上角的头像，点击下拉菜单的***accesskeys***进入用户信息管理的安全信息管理  
  AccessKey ID即为ali_AccessKey_ID  
  Access Key Secret 即为阿里云ACCESS_KEY_SECRET

2. 腾讯云进入管理后台，悬浮右上角账号，点击下拉菜单的***访问管理***，点击左侧的访问密钥，再点击API密钥管理
  SecretId即为qcloud_SecretId
  SecretKey即为qcloud_SecretKey
  
  
获取插件
---------------
1. 下载二进制发行版 [点此下载](https://resource-1251086492.cos.ap-shanghai.myqcloud.com/opensource/certbot-feehi-hook) 下载后请记得***chmod +x certbot-feehi-hook***给予执行权限

2. 源码编译安装  
    ```bash
        $ go get github.com/liufee/certbot-feehi-hook
        $ cd $GOPATH/src/liufee/certbot-feehi-hook
        $ sh build.sh
    ```

使用方式
---------------
>为了方便展示，以下shell命令安装均做了变形处理，如执行失败，请手动去掉换行符，将命令调整为单行
>以下示例命令均为dns解析在阿里云的，如果为腾讯云请将所有--type=aliyun换成--type=qcloud，--ali_AccessKey_ID=阿里云ACCESS_KEY_ID换成--qcloud_SecretId=腾讯云secretId,--qcloud_SecretKey=腾讯云secretKey

1. 直接使用  
    1.1 安装certbot
    ```bash
      $ wget https://dl.eff.org/certbot-auto && mv certbot-auto certbot && chmod +x certbot
      $ mv certbot /usr/bin
    ```
    1.2
      ```bash
      $ certbot certonly 
            -d *.您的域名.com --manual --preferred-challenges dns
            --manual-auth-hook "/存放certbot-feehi-hook的目录(下载的certboot-feehi-hook存放目录)/certbot-feehi-hook --type=aliyun --action=add --ali_AccessKey_ID=阿里云ACCESS_KEY_ID --ali_Access_Key_Secret=阿里云ACCESS_KEY_SECRET" 
            --manual-cleanup-hook "/hook/certbot-feehi-hook --type=aliyun --action=delete --ali_AccessKey_ID=阿里云ACCESS_KEY_ID --ali_Access_Key_Secret=阿里云ACCESS_KEY_SECRET"
      ```

2. 通过docker certbot/certbot官方镜像
    * 申请证书
    ```bash
        $ docker run -it --rm --name certbot 
            -v "/宿主机证书存放目录:/etc/letsencrypt" 
            -v "/tmp:/var/log/letsencrypt" 
            -v "/存放certbot-feehi-hook的目录(下载的certboot-feehi-hook存放目录):/hook" 
            certbot/certbot certonly 
            -d *.您的域名.com --manual --preferred-challenges dns
            --manual-auth-hook "/hook/certbot-feehi-hook --type=aliyun --action=add --ali_AccessKey_ID=阿里云ACCESS_KEY_ID --ali_Access_Key_Secret=阿里云ACCESS_KEY_SECRET" 
            --manual-cleanup-hook "/hook/certbot-feehi-hook --type=aliyun --action=delete --ali_AccessKey_ID=阿里云ACCESS_KEY_ID --ali_Access_Key_Secret=阿里云ACCESS_KEY_SECRET"
     ```
   * 更新证书
   ```bash
       $ /path/to/docker run -it --rm --name certbot \
           -v "/宿主机证书存放目录:/etc/letsencrypt" \
           -v "/tmp:/var/log/letsencrypt"
           -v "/存放certbot-feehi-hook的目录(下载的certboot-feehi-hook存放目录):/hook" 
           certbot/certbot renew 
           --manual --preferred-challenges dns
           --manual-auth-hook "/hook/certbot-feehi-hook.sh --type=aliyun --action=add --ali_AccessKey_ID=阿里云ACCESS_KEY_ID --ali_Access_Key_Secret=阿里云ACCESS_KEY_SECRET" 
           --manual-cleanup-hook "/hook/certbot-feehi-hook.sh --type=aliyun --action=delete --ali_AccessKey_ID=阿里云ACCESS_KEY_ID --ali_Access_Key_Secret=阿里云ACCESS_KEY_SECRET"
           --deploy-hook  "service nginx restart"
    ```
        可以配置定时任务每个月执行一次: 0 0 1 * * 上面的命令

        当且仅当成功更新证书后会执行--deploy-hook,根据自身web服务器情况进行重启web服务器重新加载新证书
       
       
帮助
---------------
1. QQ群 258780872

2. 微信 <br> ![微信](http://img-1251086492.cosgz.myqcloud.com/github/wechat.png)

3. Email job@feehi.com