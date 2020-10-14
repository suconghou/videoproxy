## videoproxy

youtube stream proxy 


## API

接口前缀均为`/video`

**解析接口**

GET `/video/{ID}.json`

GET `/video/{ID}/{ITAG}.mp4` `/video/{ID}/{ITAG}.webm`

GET `/video/{ID}/{ITAG}/{TS}.ts`

GET `/video/{ID}.jpg` `/video/{ID}.webp`

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

所有接口支持`callback`参数,可用于JSONP

maxResults 取值范围 1-50

## 环境变量

`YOUTUBE_API_KEY`

> data api 用到的key, data api 数据默认缓存48h

`DISABLE_GZIP`

> data api 默认开启gzip,若前置nginx开启gzip,可关闭此项


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


若都不配置,也可使用GO的代理配置,对所有请求有效

`export http_proxy=http://0.0.0.0:1087;export https_proxy=http://0.0.0.0:1087;`

**ID混淆**

使用环境变量`CODE_PASS`开启ID混淆

值为两个正整数用`,`隔开用于解码,例如 `10,20`

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