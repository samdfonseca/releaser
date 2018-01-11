package wiki

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	wikiUrl       = "https://wiki.axialmarket.com/api.php"
	testPageTitle = "Samtest"
)

func TestCreateWikiClient(t *testing.T) {
	client, err := NewWikiClient(wikiUrl)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.NotEqual(t, client, nil, "Client should not be nil")
}

func TestEditWikiPage(t *testing.T) {
	client, err := NewWikiClient(wikiUrl)
	assert.Equal(t, err, nil, "Error should be nil")
	newPageText := fmt.Sprintf("sam test %d", time.Now().Unix())
	err = client.UpdatePageText(testPageTitle, newPageText)
	assert.Equal(t, err, nil, "Error should be nil")
	pageText, err := client.ReadPage(testPageTitle)
	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, newPageText, pageText, "Page should be updated with new text")
}
