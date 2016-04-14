# linebotapi for golang
## example server
### echo server on GAE

``` go
package linebotapi_gae

import (
    "github.com/mokejp/linebotapi"

    "fmt"
    "io/ioutil"
    "net/http"
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

    // LINE Bot credential
    cred := linebotapi.Credential{
        ChannelId: ****,        // Your Channel ID
        ChannelSecret: "****",  // Your Channel Secret
        Mid: "****",            // Your MID
    }

    // Parse request body
    dump, err := httputil.DumpRequest(r, true)
    log.Debugf(c, string(dump))

    result, err := linebotapi.ParseRequest(r, cred)
    if err != nil {
        panic(err)
    }

    // initialize bot API client
    client := linebotapi.NewClient(cred)
    client.HttpClient = &http.Client{
        Transport: &urlfetch.Transport{
            Context: c,
        },
    }

    // Process received messages
    var messages = make(map[string][]linebotapi.MessageContent)
    for _, event := range result {
        content := event.GetEventContent()
        _, exists := messages[content.From]
        if !exists {
            messages[content.From] = make([]linebotapi.MessageContent, 0)
        }
        if content.IsOperation {
            // Operation Event
            if content.OpType == linebotapi.OpTypeAdded {
                // Added friend
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageText("Thank you!"));
            }
        } else {
            if content.ContentType == linebotapi.ContentTypeText {
                msg, err := content.GetMessageText()
                if err != nil {
                    panic(err)
                }
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageText(msg.Text));
            } else if (content.ContentType == linebotapi.ContentTypeImage) {
                _, err := content.GetMessageImage()
                if err != nil {
                    panic(err)
                }
                data, err := client.GetMessageContent(content)
                if err != nil {
                    panic(err)
                }
                // io.Reader to byte[]
                buf, err := ioutil.ReadAll(data.Reader)
                if err != nil {
                    panic(err)
                }
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageText(fmt.Sprintf("Type: %s Length: %d", data.ContentType, len(buf))));
            } else if (content.ContentType == linebotapi.ContentTypeVideo) {
                _, err := content.GetMessageVideo()
                if err != nil {
                    panic(err)
                }
                data, err := client.GetMessageContent(content)
                if err != nil {
                    panic(err)
                }
                // io.Reader to byte[]
                buf, err := ioutil.ReadAll(data.Reader)
                if err != nil {
                    panic(err)
                }
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageText(fmt.Sprintf("Type: %s Length: %d", data.ContentType, len(buf))));
            } else if (content.ContentType == linebotapi.ContentTypeAudio) {
                _, err := content.GetMessageAudio()
                if err != nil {
                    panic(err)
                }
                data, err := client.GetMessageContent(content)
                if err != nil {
                    panic(err)
                }
                // io.Reader to byte[]
                buf, err := ioutil.ReadAll(data.Reader)
                if err != nil {
                    panic(err)
                }
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageText(fmt.Sprintf("Type: %s Length: %d", data.ContentType, len(buf))));
            } else if (content.ContentType == linebotapi.ContentTypeLocation) {
                msg, err := content.GetMessageLocation()
                if err != nil {
                    panic(err)
                }
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageLocation(msg.Text, msg.Title, msg.Latitude, msg.Longitude));
            } else if (content.ContentType == linebotapi.ContentTypeSticker) {
                msg, err := content.GetMessageSticker()
                if err != nil {
                    panic(err)
                }
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageSticker(msg.StickerPackageId, msg.StickerId, ""));
            } else if (content.ContentType == linebotapi.ContentTypeContact) {
                msg, err := content.GetMessageContact()
                if err != nil {
                    panic(err)
                }
                messages[content.From] = append(messages[content.From], linebotapi.NewMessageText(fmt.Sprintf("mid: %s displayName: %s", msg.Mid, msg.DisplayName)));
            }
        }
    }
    for k := range messages {
        contacts, err := client.GetUserProfiles([]string{k})
        if err != nil {
            panic(err)
        }
        messages[k] = append(messages[k], linebotapi.NewMessageText(fmt.Sprintf("Hello, %s", contacts.Contacts[0].DisplayName)));
        err = client.SendMessages([]string{k}, messages[k], 0)
        if err != nil {
            panic(err)
        }
    }
}

```
