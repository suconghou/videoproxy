## videoproxy

youtube stream proxy 


## API

接口前缀均为`/video`

**解析接口**

GET `/video/{ID}.json` 

> 默认便于缓存不输出下载地址,query参数`info=all`输出下载地址 

GET `/video/{ID}.xml`

> 获取字幕, query参数`lang`,可以获取指定语言的字幕,哪些语言可用可以从json接口中取到

GET `/video/{ID}.mpd`

> 默认输出一条音频,一条视频;query参数配置输出的音频和视频质量
>
> `a`控制音频,默认低质量webm`249,250,251,600,140,599`,若要高质量webm使用`251,250,249,600,140,599`
> 对于不支持webm的要使用mp4可写`140,599`
>
> `v`控制视频,默认中等质量webm`247,136,244,135,243,134,242,133,278,160`,若要超高清可以写`248,137,303,399,271`等,若要mp4,可以写`136,135,134,133,160`
> 
> 可参考 https://gist.github.com/AgentOak/34d47c65b1d28829bb17c24c04a0096f

GET `/video/{ID}/{ITAG}.mp4` `/video/{ID}/{ITAG}.webm`

> proxy指定itag的资源,如果发起的是range请求,也支持响应range

GET `/video/{ID}/{ITAG}/{TS}.ts`

> proxy指定itag的指定range片段

GET `/video/{ID}.jpg` `/video/{ID}.webp`

> proxy资源banner图

GET `/video/{ID}.mp4` `/video/{ID}.webm` 

> 默认中等清晰度的音视频流
>
> query参数`prefer`配置清晰度优先级,根据itag列表搜寻可用资源,例如`prefer=18,22`
>
> 参数 `download=1` 开启弹出下载框


**6个内容接口**

GET `/video/api/v3/videos` 

> id : 视频ID
>
> chart=mostPopular&maxResults=20
> 
> regionCode

GET `/video/api/v3/search` 

> q : 搜索词,必填
> 
> type :  类型,可为`video` 选填
> 
> order : 选填
> 
> channelId : 选填
>
> pageToken : 选填
>
> maxResults
>
> relatedToVideoId
>
> regionCode

GET `/video/api/v3/channels` 

> forUsername
> 
> id
>
> categoryId 

GET `/video/api/v3/playlists` 

> channelId
> 
> id

GET `/video/api/v3/playlistItems` 

> playlistId
>
> pageToken
>
> maxResults

GET `/video/api/v3/videoCategories` 

> id 

所有接口支持`callback`参数,可用于JSONP

maxResults 取值范围 1-50

## 环境变量

`YOUTUBE_API_KEY`

> data api 用到的key, data api 数据默认缓存48h


命令行参数

> -p listen port
>
> -h bind address
>

**代理选项**

代理分为三部分,均由环境变量配置

1. image proxy

> `IMAGE_PROXY`
>
> 可访问时,一般不需要使用代理

2. video parser proxy & video stream proxy

> `VIDEO_PROXY`
>
> 此二者的代理是保持一致的

3. data api proxy

> `API_PROXY`
>
> 可以单独配置它的代理,不过一般无需配置

代理格式可写http或者socks5格式,程序自动适配

http格式: `http://proxyIp:proxyPort`

socks5格式: `127.0.0.1:9050`

如果socks5有代理验证,配置对应环境变量的键补充`_USER` `_PASSWORD` 即可,
如`VIDEO_PROXY_USER` `VIDEO_PROXY_PASSWORD`

若都不配置,也可使用GO的代理配置,对所有请求有效

`export http_proxy=http://0.0.0.0:1087;export https_proxy=http://0.0.0.0:1087;`

**ID混淆**

使用环境变量`CODE_PASS`开启ID混淆

值为两个正整数用`,`隔开用于解码,例如 `10,20`


## 数据库与缓存

> DB_AUTH 账户和密码

> DB_ADDR 数据库地址

> DB_NAME 要使用的数据库

`export DB_AUTH="work:123456" DB_ADDR="127.0.0.1:3306" DB_NAME="test"`

如果不配置数据库信息,可以配置上游白名单 `BASE_URL` 

`BASE_URL=http://domain1/video;http://domain2/video`

验证白名单时向上游发出请求,验证白名单;

如果此`BASE_URL`也不配置,则默认全部放行

数据表

本系统使用4张数据表

> whitelist 白名单

> captions 字幕缓存表

> cachempd mpd信息缓存表

> cachejson 播放信息缓存表




## run in docker

```
docker run -d -m 100m \
--name videoproxy \
--restart always \
-p 0.0.0.0:8080:6060 \
-e YOUTUBE_API_KEY="yourkey" \
suconghou/videoproxy
```

## run in k8s


创建deployment
```
kubectl run videoproxy \
--image=suconghou/videoproxy \
--env="YOUTUBE_API_KEY=yourkey" \
--port=6060 \
--limits="cpu=200m,memory=100Mi"
```

创建service

```
kubectl expose deployment/videoproxy \
--type="NodePort" \
--port=6060 \
--name=videoproxy
```

对应删除命令

```
kubectl delete deployment videoproxy 
kubectl delete service videoproxy
```