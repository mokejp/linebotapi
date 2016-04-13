package linebotapi

import (
    "testing"
    "fmt"
    "bytes"
    "net/url"
    "net/http"
    "net/http/httptest"
    "net/http/httputil"
)


func Test_ParseRequest_Success(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, err := ParseRequest(r)
        if err != nil {
            t.Error(err)
            return
        }
    }))
    defer server.Close()

    _, err := http.Post(server.URL, "application/json", bytes.NewBufferString(`{"result":[{"content":{"toType":1,"createdTime":1460529367936,"from":"abced","location":null,"id":"1","to":["abced"],"text":"fawef","contentMetadata":{"AT_RECV_MODE":"2","EMTVER":"4"},"deliveredTime":0,"contentType":1,"seq":null},"createdTime":1460529367957,"eventType":"138311609000106303","from":"abced","fromChannel":1341301815,"id":"WB1519-3361608589","to":["abced"],"toChannel":2345}]}`))
    if err != nil {
        t.Error(err)
        return
    }
}

