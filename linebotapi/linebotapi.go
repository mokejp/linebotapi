package linebotapi

import (
    "io"
    "bytes"
    "errors"
    "strconv"
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

type Credential struct {
    ChannelId int
    ChannelSecret string
    ChannelMid string
}

type Result struct {
  Result []Message `json:"result,omitempty"`
}

type Message struct {
    From string `json:"from,omitempty"`
    FromChannel int `json:"fromChannel,omitempty"`
    To []string `json:"to,omitempty"`
    ToChannel int `json:"toChannel,omitempty"`
    EventType string `json:"eventType,omitempty"`
    Id string `json:"id,omitempty"`
    Content *MessageContent `json:"content,omitempty"`
}
type MessageContent struct {
    Id string `json:"id,omitempty"`
    ContentType uint8 `json:"contentType,omitempty"`
    From string `json:"from,omitempty"`
    CreatedTime int `json:"createdTime,omitempty"`
    To []string `json:"to,omitempty"`
    ToType int `json:"toType,omitempty"`
    // Text
    Text string `json:"text,omitempty"`
    // Image / Video / Audio
    OriginalContentUrl string `json:"originalContentUrl,omitempty`
    PreviewImageUrl string `json:"previewImageUrl,omitempty`
    // Location
    Location *MessageContentLocation `json:"location,omitempty"`
    // Audio / Sticker
    ContentMetadata map[string]string `json:"contentMetadata,omitempty"`
}
type MessageContentLocation struct {
    Title string `json:"title,omitempty"`
    Latitude float64 `json:"latitude,omitempty"`
    Longitude float64 `json:"longitude,omitempty"`
}
type MessageContentSticker struct {
    Id string            // STKID
    PackageId string     // STKPKGID
    Version string       // STKVER
}
type MessageContentImage struct {
    OriginalContentUrl string
    PreviewImageUrl string
}
type MessageContentVideo struct {
    OriginalContentUrl string
    PreviewImageUrl string
}
type MessageContentAudio struct {
    OriginalContentUrl string
    Length int
}

type ErrorResponse struct {
    StatusCode string `json:"statusCode,omitempty"`
    StatusMessage string `json:"statusMessage,omitempty"`
}

func SendMessage(client http.Client, cred Credential, to []string, content MessageContent) error {
    m := Message{
        To: to,
        ToChannel: 1383378250,
        EventType: "138311608800106203",
        Content: &content,
    }
    b, err := json.Marshal(m)
    if err != nil {
        return err
    }
    req, err := http.NewRequest("POST", "https://trialbot-api.line.me/v1/events", bytes.NewBuffer(b))
    req.Header.Set("Content-Type", "application/json; charset=UTF-8")
    req.Header.Set("X-Line-ChannelID", strconv.Itoa(cred.ChannelId))
    req.Header.Set("X-Line-ChannelSecret", cred.ChannelSecret)
    req.Header.Set("X-Line-Trusted-User-With-ACL", cred.ChannelMid)

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

func SendMessageText(client http.Client, cred Credential, to []string, text string) error {
    return SendMessage(client, cred, to, MessageContent{
        ContentType: ContentTypeText,
        ToType: ToTypeUser,
        Text: text,
    })
}

func SendMessageImage(client http.Client, cred Credential, to []string, image MessageContentImage) error {
    return SendMessage(client, cred, to, MessageContent{
        ContentType: ContentTypeImage,
        ToType: ToTypeUser,
        OriginalContentUrl: image.OriginalContentUrl,
        PreviewImageUrl: image.PreviewImageUrl,
    })
}

func SendMessageVideo(client http.Client, cred Credential, to []string, video MessageContentVideo) error {
    return SendMessage(client, cred, to, MessageContent{
        ContentType: ContentTypeVideo,
        ToType: ToTypeUser,
        OriginalContentUrl: video.OriginalContentUrl,
        PreviewImageUrl: video.PreviewImageUrl,
    })
}

func SendMessageAudio(client http.Client, cred Credential, to []string, audio MessageContentAudio) error {
    return SendMessage(client, cred, to, MessageContent{
        ContentType: ContentTypeAudio,
        ToType: ToTypeUser,
        OriginalContentUrl: audio.OriginalContentUrl,
        ContentMetadata: map[string]string{
            "AUDLEN": strconv.Itoa(audio.Length),
        },
    })
}

func SendMessageLocation(client http.Client, cred Credential, to []string, location MessageContentLocation) error {
    return SendMessage(client, cred, to, MessageContent{
        ContentType: ContentTypeLocation,
        ToType: ToTypeUser,
        Location: &location,
    })
}

func SendMessageSticker(client http.Client, cred Credential, to []string, sticker MessageContentSticker) error {
    metadata := map[string]string{
        "STKID": sticker.Id,
        "STKPKGID": sticker.PackageId,
    }
    if (sticker.Version != "") {
        metadata["STKVER"] = sticker.Version
    }
    return SendMessage(client, cred, to, MessageContent{
        ContentType: ContentTypeSticker,
        ToType: ToTypeUser,
        ContentMetadata: metadata,
    })
}

func ParseRequest(body io.Reader) (Result, error) {
    decoder := json.NewDecoder(body)
    var result Result;
    err := decoder.Decode(&result)
    if err != nil {
        return Result{}, err
    }
    return result, nil
}

