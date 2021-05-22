package api_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"vladusenko.io/home-torrent/api"
)

func convertFixtureToMultipartFile(path string, fieldName string) (*bytes.Buffer, string) {
	var err error
	var file *os.File
	var formFile io.Writer

	if file, err = os.Open(path); err != nil {
		panic(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	defer writer.Close()

	if formFile, err = writer.CreateFormFile(fieldName, filepath.Base(file.Name())); err != nil {
		panic(err)
	}

	io.Copy(formFile, file)

	return body, writer.FormDataContentType()
}

func performPostRequest(r http.Handler, path string, contentType string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("POST", path, body)
	req.Header.Set("content-type", contentType)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestTorrentsApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Torrents API Suite")
}

var _ = Describe("/api/v1/torrents", func() {
	var (
		router *gin.Engine
	)

	BeforeSuite(func() {
		router = api.SetupRouter()
	})

	Describe("/api/v1/torrents/preview", func() {
		It("should successfully preview a single file torrent", func() {
			singleFileTorrentMultipart, contentType := convertFixtureToMultipartFile("../testfixtures/single-file.torrent", "torrent")
			resRecorder := performPostRequest(router, "/api/v1/torrents/preview", contentType, singleFileTorrentMultipart)
			Expect(resRecorder.Result().StatusCode).To(Equal(200))
		})
	})
})
