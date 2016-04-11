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
    Location *MessageContentLocation `json:"location,omitempty"`
    Id string `json:"id,omitempty"`
    ContentType uint8 `json:"contentType,omitempty"`
    From string `json:"from,omitempty"`
    CreatedTime int `json:"createdTime,omitempty"`
    To []string `json:"to,omitempty"`
    ToType int `json:"toType,omitempty"`
    ContentMetadata map[string]string `json:"contentMetadata,omitempty"`
    Text string `json:"text,omitempty"`
}
type MessageContentLocation struct {
    Title string `json:"title,omitempty"`
    Latitude float64 `json:"latitude,omitempty"`
    Longitude float64 `json:"longitude,omitempty"`
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

func ParseRequest(body io.Reader) (Result, error) {
    decoder := json.NewDecoder(body)
    var result Result;
    err := decoder.Decode(&result)
    if err != nil {
        return Result{}, err
    }
    return result, nil
}

