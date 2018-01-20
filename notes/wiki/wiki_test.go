package wiki

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testPageTitle = "Samtest"
)

func wikiServer(pageTitle string) *httptest.Server {
	mainResp := `{"warnings":{"info":{"*":"The intoken parameter has been deprecated."}},"query":{"pages":{"1":{"pageid":1,"ns":0,"title":"Main Page","contentmodel":"wikitext","pagelanguage":"en","touched":"2018-01-15T19:22:42Z","lastrevid":1750,"counter":309453,"length":5408,"starttimestamp":"2018-01-20T19:08:22Z","edittoken":"+\\","revisions":[{"revid":1750,"parentid":1678,"user":"172.17.0.1","anon":"","timestamp":"2018-01-15T19:22:42Z","comment":""}]}}}}`
	editResp := `{"edit":{"result":"Success","pageid":252,"title":"%s","contentmodel":"wikitext","oldrevid":1813,"newrevid":1814,"newtimestamp":"2018-01-20T19:08:22Z"}}`
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.Form.Get("action") == "query" && r.Form.Get("titles") == "Main Page" {
			fmt.Fprintln(w, mainResp)
			return
		}
		if r.Form.Get("action") == "edit" && r.Form.Get("title") == pageTitle {
			fmt.Fprintln(w, fmt.Sprintf(editResp, pageTitle))
			return
		}
	}))
}

func TestEditWikiPage(t *testing.T) {
	ts := wikiServer(testPageTitle)
	defer ts.Close()
	client, err := NewWikiClient(ts.URL)
	assert.Equal(t, err, nil, "Error should be nil")
	newPageText := fmt.Sprintf("sam test %d", time.Now().Unix())
	err = client.UpdatePageText(testPageTitle, newPageText)
	assert.Equal(t, err, nil, "Error should be nil")
}
