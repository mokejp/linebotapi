## example server
### echo server on GAE

``` go
package linebotapi_gae

import (
    "github.com/mokejp/linebotapi_go"

    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "net/http/httputil"

    "google.golang.org/appengine"
    "google.golang.org/appengine/log"
    "google.golang.org/appengine/urlfetch"
)

func init() {
    http.HandleFunc("/callback", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Parse request body
    dump, err := httputil.DumpRequest(r, true)
    log.Debugf(c, string(dump))
    result, err := linebotapi.ParseRequest(r)
    if err != nil {
        panic(err)
    }

    // logging
    b, err := json.Marshal(result)
    log.Debugf(c, string(b))

    // initialize http client
    client := &http.Client{
        Transport: &urlfetch.Transport{
            Context: c,
        },
    }

    // LINE Bot credential
    cred := linebotapi.Credential{
        ChannelId: ****,        // Your Channel ID
        ChannelSecret: "****",  // Your Channel Secret
        Mid: "****",            // Your MID
    }

    // Process received messages
    var messages = make(map[string][]linebotapi.MessageContent)
    for _, m := range result.Result {
        _, exists := messages[m.Content.From]
        if !exists {
            messages[m.Content.From] = make([]linebotapi.MessageContent, 0)
        }
        if m.Content.OpType == linebotapi.OpTypeAdded {
            // Added friend
            messages[m.Content.From] = append(messages[m.Content.From], linebotapi.MessageContent{
                ContentType: linebotapi.ContentTypeText,
                ToType: linebotapi.ToTypeUser,
                Text: "Thank you!",
            });
        } else if m.Content.ContentType == linebotapi.ContentTypeText {
            // Text
            messages[m.Content.From] = append(messages[m.Content.From], linebotapi.MessageContent{
                ContentType: linebotapi.ContentTypeText,
                ToType: linebotapi.ToTypeUser,
                Text: m.Content.Text,
            });
        } else if m.Content.ContentType == linebotapi.ContentTypeLocation {
            // Location
            messages[m.Content.From] = append(messages[m.Content.From], linebotapi.MessageContent{
                ContentType: linebotapi.ContentTypeLocation,
                ToType: linebotapi.ToTypeUser,
                Text: m.Content.Text,
                Location: &linebotapi.MessageContentLocation{
                    Title: m.Content.Location.Title,
                    Latitude: m.Content.Location.Latitude,
                    Longitude: m.Content.Location.Longitude,
                },
            });
        } else if m.Content.ContentType == linebotapi.ContentTypeSticker {
            // Sticker(preinstall sticker only?)
            messages[m.Content.From] = append(messages[m.Content.From], linebotapi.MessageContent{
                ContentType: linebotapi.ContentTypeSticker,
                ToType: linebotapi.ToTypeUser,
                ContentMetadata: map[string]string{
                    "STKID": m.Content.ContentMetadata["STKID"],
                    "STKPKGID": m.Content.ContentMetadata["STKPKGID"],
                },
            });
        } else if m.Content.ContentType == linebotapi.ContentTypeImage {
            // Image
            // Get content data.
            data, err := linebotapi.GetMessageContentData(client, cred, *m.Content)
            if err != nil {
                panic(err)
            }
            // io.Reader to byte[]
            buf, err := ioutil.ReadAll(data.Reader)
            if err != nil {
                panic(err)
            }
            messages[m.Content.From] = append(messages[m.Content.From], linebotapi.MessageContent{
                ContentType: linebotapi.ContentTypeText,
                ToType: linebotapi.ToTypeUser,
                Text: fmt.Sprintf("Type: %s Length: %d", data.ContentType, len(buf)),
            });
        }
    }
    for k := range messages {
        contacts, err := linebotapi.GetUserProfiles(client, cred, []string{k})
        if err != nil {
            panic(err)
        }
        messages[k] = append(messages[k], linebotapi.MessageContent{
            ContentType: linebotapi.ContentTypeText,
            ToType: linebotapi.ToTypeUser,
            Text: "Hello, " + contacts.Contacts[0].DisplayName,
        });
        err = linebotapi.SendMessages(client, cred, []string{k}, messages[k], 0)
        if err != nil {
            panic(err)
        }
    }
}


```
