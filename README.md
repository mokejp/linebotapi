## example server
### echo server on GAE

``` go
package linebotapi_gae

import (
    "linebotapi"

    "net/http"

    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
)

func init() {
    http.HandleFunc("/callback", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    result, err := linebotapi.ParseRequest(r.Body)
    if err != nil {
        panic(err)
    }
    client := http.Client{
        Transport: &urlfetch.Transport{
            Context: c,
        },
    }

    cred := linebotapi.Credential{
        ChannelId: "****",
        ChannelSecret: "****",
        ChannelMid: "****",
    }

    m := result.Result[0]
    if m.Content.ContentType == linebotapi.ContentTypeText {
        // Text
        mc := linebotapi.MessageContent{
            ContentType: linebotapi.ContentTypeText,
            ToType: linebotapi.ToTypeUser,
            Text: m.Content.Text,
        }
        err = linebotapi.SendMessage(client, cred, []string{m.Content.From}, mc)
    }
    if m.Content.ContentType == linebotapi.ContentTypeSticker {
        // Sticker(preinstall sticker only?)
        mc := linebotapi.MessageContent{
            ContentType: linebotapi.ContentTypeSticker,
            ToType: linebotapi.ToTypeUser,
            ContentMetadata: map[string]string{
                "STKID": m.Content.ContentMetadata["STKID"],
                "STKPKGID": m.Content.ContentMetadata["STKPKGID"],
            },
        }
        err = linebotapi.SendMessage(client, cred, []string{m.Content.From}, mc)
    }
    
    if err != nil {
        panic(err)
    }
}

```
