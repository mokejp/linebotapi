package linebotapi

import (
    "io"
    "bytes"
    "errors"
    "strconv"
    "strings"
    "net/http"
    "encoding/json"
)

const (
    ContentTypeText     = 1
    ContentTypeImage    = 2
    ContentTypeVideo    = 3
    ContentTypeAudio    = 4
    ContentTypeLocation = 7
    ContentTypeSticker  = 8
)

const (
    ToTypeUser = 1
)

const (
    OpTypeAdded   = 4
    OpTypeBlocked = 8
)

type Credential struct {
    ChannelId int
    ChannelSecret string
    Mid string
}

type Result struct {
    Result []Message `json:"result,omitempty"`
}

type Message struct {
    Id string `json:"id,omitempty"`
    From string `json:"from,omitempty"`
    FromChannel int `json:"fromChannel,omitempty"`
    To []string `json:"to,omitempty"`
    ToChannel int `json:"toChannel,omitempty"`
    EventType string `json:"eventType,omitempty"`
    Content *MessageContent `json:"content,omitempty"`
}
type MessageContent struct {
    Id string `json:"id,omitempty"`
    From string `json:"from,omitempty"`
    CreatedTime int `json:"createdTime,omitempty"`
    To []string `json:"to,omitempty"`
    ToType int `json:"toType,omitempty"`

    // User message
    Location *MessageContentLocation `json:"location,omitempty"`
    ContentType uint8 `json:"contentType,omitempty"`
    ContentMetadata map[string]string `json:"contentMetadata,omitempty"`
    Text string `json:"text,omitempty"`

    // User operation
    OpType int `json:"opType,omitempty"`
    Revision int `json:"revision,omitempty"`
    Params []string `json:"params,omitempty"`

    // For sending messages
    MessageNotified int `json:"messageNotified,omitempty"`
    Messages []MessageContent `json:"messages,omitempty"`

    // For sending image / video / audio
    OriginalContentUrl string `json:"originalContentUrl,omitempty"`
    PreviewImageUrl string `json:"previewImageUrl,omitempty"`
}
type MessageContentData struct {
    Reader io.Reader
    ContentType string
}
type MessageContentLocation struct {
    Title string `json:"title,omitempty"`
    Latitude float64 `json:"latitude,omitempty"`
    Longitude float64 `json:"longitude,omitempty"`
}

type Contacts struct {
    Contacts []Contact `json:"contact,omitempty`
    Count int `json:"count,omitempty`
    Total int `json:"total,omitempty`
    Start int `json:"start,omitempty`
    Display int `json:"display,omitempty`
}
type Contact struct {
    DisplayName string `json:"displayName,omitempty`
    Mid string `json:"mid,omitempty`
    PictureUrl string `json:"pictureUrl,omitempty`
    StatusMessage string `json:"statusMessage,omitempty`
}

type ErrorResponse struct {
    StatusCode string `json:"statusCode,omitempty"`
    StatusMessage string `json:"statusMessage,omitempty"`
}

func postEvent(client *http.Client, cred Credential,  to []string, m Message) error {
    b, err := json.Marshal(m)
    if err != nil {
        return err
    }
    req, err := http.NewRequest("POST", "https://trialbot-api.line.me/v1/events", bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json; charset=UTF-8")
    req.Header.Set("X-Line-ChannelID", strconv.Itoa(cred.ChannelId))
    req.Header.Set("X-Line-ChannelSecret", cred.ChannelSecret)
    req.Header.Set("X-Line-Trusted-User-With-ACL", cred.Mid)

    res, err := client.Do(req)
    if err != nil {
        return err
    }
    if res.StatusCode != http.StatusOK {
        decoder := json.NewDecoder(res.Body)
        var e ErrorResponse;
        err := decoder.Decode(&e)
        if err != nil {
            return err
        }
        return errors.New(e.StatusMessage)
    }
    defer res.Body.Close()
    return nil
}

func SendMessage(client *http.Client, cred Credential, to []string, content MessageContent) error {
    return postEvent(client, cred, to, Message{
        To: to,
        ToChannel: 1383378250,
        EventType: "138311608800106203",
        Content: &content,
    })
}

func SendMessages(client *http.Client, cred Credential, to []string, contents []MessageContent, notified int) error {
    return postEvent(client, cred, to, Message{
        To: to,
        ToChannel: 1383378250,
        EventType: "140177271400161403",
        Content: &MessageContent{
            MessageNotified: notified,
            Messages: contents,
        },
    })
}

func GetMessageContentData(client *http.Client, cred Credential, m MessageContent) (MessageContentData, error) {
    req, err := http.NewRequest("GET", "https://trialbot-api.line.me/v1/bot/message/" + m.Id + "/content", nil)
    req.Header.Set("X-Line-ChannelID", strconv.Itoa(cred.ChannelId))
    req.Header.Set("X-Line-ChannelSecret", cred.ChannelSecret)
    req.Header.Set("X-Line-Trusted-User-With-ACL", cred.Mid)

    res, err := client.Do(req)
    if err != nil {
        return MessageContentData{}, err
    }
    if res.StatusCode != http.StatusOK {
        var e ErrorResponse;
        decoder := json.NewDecoder(res.Body)
        err := decoder.Decode(&e)
        if err != nil {
            return MessageContentData{}, err
        }
        return MessageContentData{}, errors.New(e.StatusMessage)
    }
    return MessageContentData{
        Reader: res.Body,
        ContentType: res.Header.Get("Content-Type"),
    }, nil
}

func GetUserProfiles(client *http.Client, cred Credential, mids []string) (Contacts, error) {
    req, err := http.NewRequest("GET", "https://trialbot-api.line.me/v1/profiles?mids=" + strings.Join(mids[:], ","), nil)
    req.Header.Set("X-Line-ChannelID", strconv.Itoa(cred.ChannelId))
    req.Header.Set("X-Line-ChannelSecret", cred.ChannelSecret)
    req.Header.Set("X-Line-Trusted-User-With-ACL", cred.Mid)

    res, err := client.Do(req)
    if err != nil {
        return Contacts{}, err
    }
    decoder := json.NewDecoder(res.Body)
    if res.StatusCode != http.StatusOK {
        var e ErrorResponse;
        err := decoder.Decode(&e)
        if err != nil {
            return Contacts{}, err
        }
        return Contacts{}, errors.New(e.StatusMessage)
    }
    var contacts Contacts;
    err = decoder.Decode(&contacts)
    if err != nil {
        return Contacts{}, err
    }
    defer res.Body.Close()
    return contacts, nil
}

func ParseRequest(r *http.Request) (Result, error) {
    decoder := json.NewDecoder(r.Body)
    var result Result;
    err := decoder.Decode(&result)
    if err != nil {
        return Result{}, err
    }
    return result, nil
}

