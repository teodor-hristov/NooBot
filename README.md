# Discord bot made for Introduction to Go course at Sofia University

## Supported slash commands:
/play `<channel name> <url>` - Play song in the given channel (bot can play youtube/soundcloud/web streams)

/stop - Stops music from the bot and disconnects the bot

/talk-history `<user>` - Returns the user voice channel connection history.

/voice-chat-top - All people talked since bot connected to the server (sorted)

## Example 
```
/play General https://soundcloud.com/lil-jairmy/alaska
```
```
/stop
```
```
/talk-history tedo3637
```
```
/voice-chat-top
```

### Build with
```
go build main.go stats.go utils.go music.go commands.go
```
```
./main.exe -token= <YOUR DISCORD APP TOKEN>
```

### Feature ideas
* Add playlist queue
* Add download song/playlist command
* Add next/pause/prev song command
* Add search song by name
