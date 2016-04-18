package linebotapi

import (
    "testing"

    "fmt"
    "bytes"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/json"
    "encoding/base64"
)


func Test_ParseRequest_Success(t *testing.T) {
    cred := Credential{
        ChannelId: 1234567890,
        ChannelSecret: "0123456789abcdef0000000000000000",
        Mid: "0123456789abcdef0000000000000000",
    }

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        result, err := ParseRequest(r, cred)
        if err != nil {
            t.Error(err)
            return
        }
        for _, event := range result {
            content := event.GetEventContent()
            if content.IsMessage {
                if content.ContentType == ContentTypeText {
                    msg, err := content.GetMessageText()
                    if err != nil {
                        t.Error(err)
                        return
                    }
                    if msg.Text != "fawef" {
                        t.Error("invalid text")
                        return
                    }
                }
            }
        }
    }))
    defer server.Close()

    buf := bytes.NewBufferString(`{"result":[{"content":{"toType":1,"createdTime":1460529367936,"from":"abced","location":null,"id":"1","to":["abced"],"text":"fawef","contentMetadata":{"AT_RECV_MODE":"2","EMTVER":"4"},"deliveredTime":0,"contentType":1,"seq":null},"createdTime":1460529367957,"eventType":"138311609000106303","from":"abced","fromChannel":1341301815,"id":"WB1519-3361608589","to":["abced"],"toChannel":2345}]}`)
    mac := hmac.New(sha256.New, []byte(cred.ChannelSecret))
    mac.Write(buf.Bytes())
    messageMAC := mac.Sum(nil)
    sign := base64.StdEncoding.EncodeToString(messageMAC)

    c := http.Client{}
    req, err := http.NewRequest("POST", server.URL, buf)
    if err != nil {
        t.Error(err)
        return
    }
    req.Header.Set("Content-Type", "application/json; charset=UTF-8")
    req.Header.Set("X-LINE-ChannelSignature", sign)
    _, err = c.Do(req)
    if err != nil {
        t.Error(err)
        return
    }
}

func Test_SendMessage_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        decoder := json.NewDecoder(r.Body)
        var event Event;
        err := decoder.Decode(&event)
        if err != nil {
            t.Error(err)
            return
        }
        if event.RawContent["text"] != "message" {
            t.Error("'message'")
            t.Errorf("excepted: 'message', actual: '%s'", event.RawContent["text"])
            return
        }
    }))
    defer server.Close()

    cred := Credential{
        ChannelId: 1234567890,
        ChannelSecret: "abcdefg",
        Mid: "abcdefg",
    }
    client := NewClient(cred)
    client.BaseURL = server.URL
    err := client.SendMessage([]string{"test"}, NewMessageText("message"))
    if err != nil {
        t.Error(err)
        return
    }
}

func Test_SendMessage_Failure(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(400)
        fmt.Fprintf(w, `{"statusCode":"400","statusMessage":"error"}`)
    }))
    defer server.Close()

    cred := Credential{
        ChannelId: 1234567890,
        ChannelSecret: "abcdefg",
        Mid: "abcdefg",
    }
    client := NewClient(cred)
    client.BaseURL = server.URL
    err := client.SendMessage([]string{"test"}, NewMessageText("message"))
    if err == nil {
        t.Error("err is nil")
        return
    }
}

func Test_GetUserProfiles_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        q := r.URL.Query()
        if q["mids"][0] != "u0047556f2e40dba2456887320ba7c76d" {
            t.Errorf("excepted: 'u0047556f2e40dba2456887320ba7c76d', actual: '%s'", q["mids"][0])
        }
        w.Header().Set("Content-type", "application/json")
        w.WriteHeader(200)
        fmt.Fprintf(w, `{"contacts":[{"displayName":"BOT API","mid":"u0047556f2e40dba2456887320ba7c76d","pictureUrl":"http://dl.profile.line.naver.jp/abcdefghijklmn","statusMessage":"Hello, LINE!"}]}`)
    }))
    defer server.Close()

    cred := Credential{
        ChannelId: 1234567890,
        ChannelSecret: "abcdefg",
        Mid: "abcdefg",
    }
    client := NewClient(cred)
    client.BaseURL = server.URL
    contacts, err := client.GetUserProfiles([]string{"u0047556f2e40dba2456887320ba7c76d"})
    if err != nil {
        t.Error(err)
        return
    }
    if contacts.Contacts[0].DisplayName != "BOT API" {
        t.Errorf("excepted: 'BOT API', actual: '%s'", contacts.Contacts[0].DisplayName)
        return
    }
}

func Test_GetMessageContent_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-type", "application/json")
        w.WriteHeader(200)
        fmt.Fprintf(w, `{}`)
    }))
    defer server.Close()

    cred := Credential{
        ChannelId: 1234567890,
        ChannelSecret: "abcdefg",
        Mid: "abcdefg",
    }
    client := NewClient(cred)
    client.BaseURL = server.URL

    data, err := client.GetMessageContent(&EventContent{})
    if err != nil {
        t.Error(err)
        return
    }
    defer data.Reader.Close()
    buf, err := ioutil.ReadAll(data.Reader)
    if err != nil {
        t.Error(err)
        return
    }
    if len(buf) != 2 {
        t.Errorf("excepted: 2, actual: %d", len(buf))
    }
    if data.ContentType != "application/json" {
        t.Errorf("excepted: 'application/json', actual: '%s'", data.ContentType)
    }
}
