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
        ChannelId: ****,        // Your Channel ID
        ChannelSecret: "****",  // Your Channel Secret
        Mid: "****",            // Your MID
    }

    var messages = make([]linebotapi.MessageContent, len(result.Result))
    for i, m := range result.Result {
        if m.Content.ContentType == linebotapi.ContentTypeText {
            // Text
            messages[i] = linebotapi.MessageContent{
                ContentType: linebotapi.ContentTypeText,
                ToType: linebotapi.ToTypeUser,
                Text: m.Content.Text,
            }
        }
        if m.Content.ContentType == linebotapi.ContentTypeSticker {
            // Sticker(preinstall sticker only?)
            messages[i] = linebotapi.MessageContent{
                ContentType: linebotapi.ContentTypeSticker,
                ToType: linebotapi.ToTypeUser,
                ContentMetadata: map[string]string{
                    "STKID": m.Content.ContentMetadata["STKID"],
                    "STKPKGID": m.Content.ContentMetadata["STKPKGID"],
                },
            }
        }
        err = linebotapi.SendMessages(client, cred, []string{m.Content.From}, messages, 0)
        if err != nil {
            panic(err)
        }
    }
}


```
