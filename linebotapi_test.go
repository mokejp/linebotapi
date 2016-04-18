package linebotapi

import (
    "testing"

    "fmt"
    "bytes"
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
