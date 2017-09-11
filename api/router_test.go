package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rkcpi/vell/config"
	"github.com/rkcpi/vell/repos"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type mockStore struct{}
type mockRepo struct{}

func (s *mockStore) ListRepositories() []repos.Repository {
	return []repos.Repository{{Name: "foo"}}
}
func (s *mockStore) Initialize(name string) error        { return nil }
func (s *mockStore) Get(name string) repos.AnyRepository { return &mockRepo{} }

func (s *mockRepo) Add(filename string, f io.Reader) error { return errors.New("terribly sorry") }
func (s *mockRepo) Update() error { return errors.New("terribly sorry") }
func (s *mockRepo) ListPackages() ([]repos.Package, error) { return []repos.Package{}, nil }
func (s *mockRepo) PackageWithNameAndVersion(packagename string, version string) (repos.Package, error) { return repos.Package{}, nil }
func (s *mockRepo) IsValid() bool { return true }

func setup() {
	config.RepoStore = &mockStore{}
}

func TestBasicRepositoryListing(t *testing.T) {
	r := NewRouter()

	// it's safe to ignore error here, because we're manually entering URL
	req, _ := http.NewRequest("GET", "http://localhost/repositories", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Error("Should be OK")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var reps []repos.Repository
	if err = json.Unmarshal(body, &reps); err != nil {
		t.Fatal(err)
	}
	if len(reps) != 1 {
		t.Error("Should contain foo repository")
	}
}

func TestBasicErrorHandling(t *testing.T) {
	r := NewRouter()

	reqBody := bytes.NewBuffer([]byte("asdf"))
	req, _ := mulipartPOSTReq("http://localhost/repositories/foo/packages", "file", reqBody)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Error("Should not be OK")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var apiError struct{Message string}
	if err = json.Unmarshal(body, &apiError); err != nil {
		t.Fatal(err)
	}

	if apiError.Message != "Error adding package to repository" {
		t.Error(body)
	}
}

func mulipartPOSTReq(url string, field string, f io.Reader) (req *http.Request, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile(field, "filename")
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}

	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err = http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	return
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func shutdown() {

}
