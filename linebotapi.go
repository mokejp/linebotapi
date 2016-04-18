package linebotapi

import (
    "io"
    "fmt"
    "bytes"
    "errors"
    "strconv"
    "strings"
    "net/url"
    "net/http"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/json"
    "encoding/base64"
)

const (
    ContentTypeText     = 1
    ContentTypeImage    = 2
    ContentTypeVideo    = 3
    ContentTypeAudio    = 4
    ContentTypeLocation = 7
    ContentTypeSticker  = 8
    ContentTypeContact  = 10
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

type Event struct {
    Id string `json:"id,omitempty"`
    From string `json:"from,omitempty"`
    FromChannel int `json:"fromChannel,omitempty"`
    To []string `json:"to,omitempty"`
    ToChannel int `json:"toChannel,omitempty"`
    EventType string `json:"eventType,omitempty"`
    RawContent map[string]interface{} `json:"content"`
}
func (c *Event) GetEventContent() *EventContent {
    toArray := c.RawContent["to"].([]interface {})
    to := make([]string, len(toArray))
    for i, item := range toArray {
        to[i] = item.(string)
    }
    content := &EventContent{
        Event: c,
        Id: c.RawContent["id"].(string),
        From: c.RawContent["from"].(string),
        CreatedTime: int(c.RawContent["createdTime"].(float64)),
        To: to,
        ToType: uint8(c.RawContent["toType"].(float64)),
    }
    opType, exists := c.RawContent["opType"]
    if exists {
        content.OpType = uint8(opType.(float64))
        content.IsOperation = true
    }
    contentType, exists := c.RawContent["contentType"]
    if exists {
        content.ContentType = uint8(contentType.(float64))
        content.IsMessage = true
    }
    return content
}

type Mapper interface {
    Map() map[string]interface{}
}

type MessageContent struct {
    ContentType uint8
    Content Mapper
}


type EventContent struct {
    Event *Event
    Id string
    From string
    CreatedTime int
    To []string
    ToType uint8
    IsOperation bool
    IsMessage bool
    OpType uint8
    ContentType uint8
}
func (c *EventContent) GetMessageText() (*MessageText, error) {
    if c.ContentType != ContentTypeText {
        return nil, errors.New("invalid contentType")
    }
    return &MessageText{
        Text: c.Event.RawContent["text"].(string),
    }, nil
}
func (c *EventContent) GetMessageImage() (*MessageImage, error) {
    if c.ContentType != ContentTypeImage {
        return nil, errors.New("invalid contentType")
    }
    return &MessageImage{
    }, nil
}
func (c *EventContent) GetMessageVideo() (*MessageVideo, error) {
    if c.ContentType != ContentTypeVideo {
        return nil, errors.New("invalid contentType")
    }
    return &MessageVideo{
    }, nil
}
func (c *EventContent) GetMessageAudio() (*MessageAudio, error) {
    if c.ContentType != ContentTypeAudio {
        return nil, errors.New("invalid contentType")
    }
    metadata := c.Event.RawContent["contentMetadata"].(map[string]interface{})
    length, err := strconv.Atoi(metadata["AUDLEN"].(string))
    if err != nil {
        return nil, err
    }
    return &MessageAudio{
        AudioLength: length,
    }, nil
}
func (c *EventContent) GetMessageLocation() (*MessageLocation, error) {
    if c.ContentType != ContentTypeLocation {
        return nil, errors.New("invalid contentType")
    }
    location := c.Event.RawContent["location"].(map[string]interface{})
    return &MessageLocation{
        Text: c.Event.RawContent["text"].(string),
        Title: location["title"].(string),
        Latitude: location["latitude"].(float64),
        Longitude: location["latitude"].(float64),
    }, nil
}
func (c *EventContent) GetMessageSticker() (*MessageSticker, error) {
    if c.ContentType != ContentTypeSticker {
        return nil, errors.New("invalid contentType")
    }
    metadata := c.Event.RawContent["contentMetadata"].(map[string]interface{})
    return &MessageSticker{
        StickerId: metadata["STKID"].(string),
        StickerPackageId: metadata["STKPKGID"].(string),
        StickerVersion: metadata["STKVER"].(string),
    }, nil
}
func (c *EventContent) GetMessageContact() (*MessageContact, error) {
    if c.ContentType != ContentTypeContact {
        return nil, errors.New("invalid contentType")
    }
    metadata := c.Event.RawContent["contentMetadata"].(map[string]interface{})
    return &MessageContact{
        Mid: metadata["mid"].(string),
        DisplayName: metadata["displayName"].(string),
    }, nil
}

type MessageText struct {
    Text string
}
func (c *MessageText) Map() map[string]interface{} {
    return map[string]interface{}{
        "contentType": ContentTypeText,
        "toType": ToTypeUser,
        "text": c.Text,
    }
}

type MessageImage struct {
    OriginalContentUrl string
    PreviewImageUrl string
}
func (c *MessageImage) Map() map[string]interface{} {
    return map[string]interface{}{
        "contentType": ContentTypeImage,
        "toType": ToTypeUser,
        "originalContentUrl": c.OriginalContentUrl,
        "previewImageUrl": c.PreviewImageUrl,
    }
}

type MessageVideo struct {
    OriginalContentUrl string
    PreviewImageUrl string
}
func (c *MessageVideo) Map() map[string]interface{} {
    return map[string]interface{}{
        "contentType": ContentTypeVideo,
        "toType": ToTypeUser,
        "originalContentUrl": c.OriginalContentUrl,
        "previewImageUrl": c.PreviewImageUrl,
    }
}

type MessageAudio struct {
    OriginalContentUrl string
    AudioLength int
}
func (c *MessageAudio) Map() map[string]interface{} {
    return map[string]interface{}{
        "contentType": ContentTypeAudio,
        "toType": ToTypeUser,
        "originalContentUrl": c.OriginalContentUrl,
        "contentMetadata": map[string]string{
            "AUDLEN": strconv.Itoa(c.AudioLength),
        },
    }
}

type MessageLocation struct {
    Text string
    Title string
    Latitude float64
    Longitude float64
}
func (c *MessageLocation) Map() map[string]interface{} {
    return map[string]interface{}{
        "contentType": ContentTypeLocation,
        "toType": ToTypeUser,
        "text": c.Text,
        "location": map[string]interface{}{
            "title": c.Title,
            "latitude": c.Latitude,
            "longitude": c.Longitude,
        },
    }
}

type MessageSticker struct {
    StickerId string
    StickerPackageId string
    StickerVersion string
}
func (c *MessageSticker) Map() map[string]interface{} {
    return map[string]interface{}{
        "contentType": ContentTypeSticker,
        "toType": ToTypeUser,
        "contentMetadata": map[string]string{
            "STKID": c.StickerId,
            "STKPKGID": c.StickerPackageId,
            "STKVER": c.StickerVersion,
        },
    }
}

type MessageContact struct {
    Mid string
    DisplayName string
}
func (c *MessageContact) Map() map[string]interface{} {
    return map[string]interface{}{
        "contentType": ContentTypeContact,
        "toType": ToTypeUser,
        "contentMetadata": map[string]string{
            "mid": c.Mid,
            "displayName": c.DisplayName,
        },
    }
}

func NewMessageText(text string) MessageContent {
    return MessageContent{
        ContentType: ContentTypeText,
        Content: &MessageText{
            Text: text,
        },
    }
}
func NewMessageImage(contentURL, previewURL string) MessageContent {
    return MessageContent{
        ContentType: ContentTypeImage,
        Content: &MessageImage{
            OriginalContentUrl: contentURL,
            PreviewImageUrl: previewURL,
        },
    }
}
func NewMessageVideo(contentURL, previewURL string) MessageContent {
    return MessageContent{
        ContentType: ContentTypeVideo,
        Content: &MessageVideo{
            OriginalContentUrl: contentURL,
            PreviewImageUrl: previewURL,
        },
    }
}
func NewMessageAudio(contentURL string, length int) MessageContent {
    return MessageContent{
        ContentType: ContentTypeAudio,
        Content: &MessageAudio{
            OriginalContentUrl: contentURL,
            AudioLength: length,
        },
    }
}
func NewMessageLocation(text, title string, lat, long float64) MessageContent {
    return MessageContent{
        ContentType: ContentTypeLocation,
        Content: &MessageLocation{
            Text: text,
            Title: title,
            Latitude: lat,
            Longitude: long,
        },
    }
}
func NewMessageSticker(packageId, id, ver string) MessageContent {
    return MessageContent{
        ContentType: ContentTypeSticker,
        Content: &MessageSticker{
            StickerPackageId: packageId,
            StickerId: id,
            StickerVersion: ver,
        },
    }
}
/*
func NewMessageContact(mid, name string) MessageContent {
    return MessageContent{
        ContentType: ContentTypeContact,
        Content: &MessageContact{
            Mid: mid,
            DisplayName: name,
        },
    }
}
*/

type MessageContentData struct {
    Reader io.Reader
    ContentType string
}


type Contacts struct {
    Contacts []Contact `json:"contact"`
    Count int `json:"count"`
    Total int `json:"total"`
    Start int `json:"start"`
    Display int `json:"display"`
}

type Contact struct {
    DisplayName string `json:"displayName"`
    Mid string `json:"mid"`
    PictureUrl string `json:"pictureUrl"`
    StatusMessage string `json:"statusMessage"`
}


type ErrorResponse struct {
    StatusCode string `json:"statusCode,omitempty"`
    StatusMessage string `json:"statusMessage,omitempty"`
}


type Client struct {
    BaseURL string
    HttpClient *http.Client
    Credential Credential
}
func (c *Client) newRequest(method, url string, body io.Reader) (*http.Request, error) {
    req, err := http.NewRequest(method, url, body)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json; charset=UTF-8")
    req.Header.Set("X-Line-ChannelID", strconv.Itoa(c.Credential.ChannelId))
    req.Header.Set("X-Line-ChannelSecret", c.Credential.ChannelSecret)
    req.Header.Set("X-Line-Trusted-User-With-ACL", c.Credential.Mid)
    return req, nil
}

func (c *Client) handleError(resp *http.Response) error {
    decoder := json.NewDecoder(resp.Body)
    var e ErrorResponse;
    err := decoder.Decode(&e)
    if err != nil {
        return err
    }
    return errors.New(fmt.Sprintf("%s: %s", e.StatusCode, e.StatusMessage))
}

func (c *Client) postEvents(to []string, event Event) error {
    // Build endpoint URL
    url, err := url.Parse(c.BaseURL)
    if err != nil {
        return err
    }
    url.Path = "/v1/events"

    // JSON encoding
    b, err := json.Marshal(event)
    if err != nil {
        return err
    }

    // POST event
    req, err := c.newRequest("POST", url.String(), bytes.NewBuffer(b))
    if err != nil {
        return err
    }
    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return err
    }
    if resp.StatusCode != http.StatusOK {
        return c.handleError(resp)
    }
    defer resp.Body.Close()
    return nil
}

func (c *Client) SendMessage(to []string, content MessageContent) error {
    return c.postEvents(to, Event{
        To: to,
        ToChannel: 1383378250,
        EventType: "138311608800106203",
        RawContent: content.Content.Map(),
    })
}

func (c *Client) SendMessages(to []string, contents []MessageContent, notified int) error {
    messages := make([]map[string]interface{}, len(contents))
    for i, c := range contents {
        messages[i] = c.Content.Map()
    }
    return c.postEvents(to, Event{
        To: to,
        ToChannel: 1383378250,
        EventType: "140177271400161403",
        RawContent: map[string]interface{}{
            "messageNotified": notified,
            "messages": messages,
        },
    })
}

func (c *Client) GetMessageContent(m *EventContent) (*MessageContentData, error) {
    // Build endpoint URL
    url, err := url.Parse(c.BaseURL)
    if err != nil {
        return nil, err
    }
    url.Path = fmt.Sprintf("/v1/bot/message/%s/content", m.Id)
    req, err := c.newRequest("GET", url.String(), nil)
    if err != nil {
        return nil, err
    }
    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        return nil, c.handleError(resp)
    }
    defer resp.Body.Close()
    return &MessageContentData{
        Reader: resp.Body,
        ContentType: resp.Header.Get("Content-Type"),
    }, nil
}

func (c *Client) GetUserProfiles(mids []string) (*Contacts, error) {
    // Build endpoint URL
    url, err := url.Parse(c.BaseURL)
    if err != nil {
        return nil, err
    }
    url.Path = fmt.Sprintf("/v1/profiles")
    url.RawQuery = fmt.Sprintf("mids=%s", strings.Join(mids[:], ","))
    req, err := c.newRequest("GET", url.String(), nil)
    if err != nil {
        return nil, err
    }
    resp, err := c.HttpClient.Do(req)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        return nil, c.handleError(resp)
    }
    decoder := json.NewDecoder(resp.Body)
    var contacts Contacts;
    err = decoder.Decode(&contacts)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return &contacts, nil
}

func NewClient(cred Credential) *Client {
    return &Client{
        BaseURL: "https://trialbot-api.line.me",
        HttpClient: &http.Client{},
        Credential: cred,
    }
}

type callbackRequest struct {
    Result []Event
}

func ParseRequest(r *http.Request, cred Credential) ([]Event, error) {
    // Get request body
    buf := new(bytes.Buffer)
    buf.ReadFrom(r.Body)

    // Get request signature
    sign := r.Header.Get("X-LINE-ChannelSignature")
    if sign == "" {
        return nil, errors.New("Not found HTTP header : 'X-LINE-ChannelSignature'.")
    }
    expectedMAC, err := base64.StdEncoding.DecodeString(sign)
    if err != nil {
        return nil, err
    }

    // Validate body
    mac := hmac.New(sha256.New, []byte(cred.ChannelSecret))
    mac.Write(buf.Bytes())
    messageMAC := mac.Sum(nil)
    if !hmac.Equal(messageMAC, expectedMAC) {
        return nil, errors.New("Invalid signature.")
    }

    // Decode json
    decoder := json.NewDecoder(buf)
    var result callbackRequest;
    err = decoder.Decode(&result)
    if err != nil {
        return nil, err
    }
    return result.Result, nil
}
