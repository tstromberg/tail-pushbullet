Tail a file, and send your phone notifications via PushBullet once a line matches.

```shell
go get "github.com/xconstruct/go-pushbullet"
go get github.com/hpcloud/tail/...
go build tail-pushbullet.go
```

Usage:

```shell
tail-pushbullet \
  --regex "usb|keyring" \
  --api_key=YOUR_KEY /var/log/messages
```

You can get your API key from https://docs.pushbullet.com/v1/
