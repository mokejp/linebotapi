package linebotapi

import (
    "testing"

    "fmt"
    "bytes"
    "net/http"
    "net/http/httptest"
    "encoding/json"
)


func Test_ParseRequest_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cred := Credential{
            ChannelId: 1234567890,
            ChannelSecret: "abcdefg",
            Mid: "abcdefg",
        }
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

    _, err := http.Post(server.URL, "application/json", bytes.NewBufferString(`{"result":[{"content":{"toType":1,"createdTime":1460529367936,"from":"abced","location":null,"id":"1","to":["abced"],"text":"fawef","contentMetadata":{"AT_RECV_MODE":"2","EMTVER":"4"},"deliveredTime":0,"contentType":1,"seq":null},"createdTime":1460529367957,"eventType":"138311609000106303","from":"abced","fromChannel":1341301815,"id":"WB1519-3361608589","to":["abced"],"toChannel":2345}]}`))
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
