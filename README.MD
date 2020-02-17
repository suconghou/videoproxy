## youtube share 

youtube stream proxy 


## API

接口前缀均为`/video`

GET `/video/{ID}.json`

GET `/video/{ID}/{ITAG}.mp4` `/video/{ID}/{ITAG}.webm`

GET `/video/{ID}/{ITAG}/{TS}.ts`

GET `/video/{ID}.jpg` `/image/{ID}.webp`

**内容接口**

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

maxResults 取值返回 1-50

## 环境变量

`YOUTUBE_API_KEY`

