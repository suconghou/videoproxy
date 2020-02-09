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

GET `/video/api/v3/channels` 

> forUsername
> 
> id

GET `/video/api/v3/playlists` 

> channelId
> 
> id

GET `/video/api/v3/playlistItems` 

> playlistId

> pageToken

## 环境变量

`YOUTUBE_API_KEY`

