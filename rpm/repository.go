package rpm

import (
	"compress/gzip"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/rkcpi/vell/repos"
)

type yumRepository struct {
	store *yumRepoStore
	name  string
}

func NewRepository(store *yumRepoStore, name string) repos.AnyRepository {
	return &yumRepository{store, name}
}

func (r *yumRepository) ensureExists() (string, error) {
	path := r.path()

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("Creating repository directory %s", path)
		err = os.MkdirAll(path, 0755)
	}

	return path, err
}

func (r *yumRepository) path() string {
	return filepath.Join(r.store.base, r.name)
}

func (r *yumRepository) Initialize() error {
	log.Printf("Initializing repository %s", r.name)

	path, err := r.ensureExists()
	if err != nil {
		return err
	}

	log.Printf("Executing `createrepo --database %s`", path)
	cmd := exec.Command("createrepo", "--database", path)

	return cmd.Run()
}

func (r *yumRepository) Add(filename string, f io.Reader) error {
	log.Printf("Adding %s to repository %s", filename, r.path())
	destinationPath := filepath.Join(r.path(), filename)

	destination, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, f)

	return err
}

func (r *yumRepository) Update() error {
	path := r.path()
	log.Printf("Executing `createrepo --update %s`", path)
	cmd := exec.Command("createrepo", "--update", path)

	return cmd.Run()
}

func (r *yumRepository) repomdPath() string {
	return filepath.Join(r.path(), "repodata", "repomd.xml")
}

func (r *yumRepository) IsValid() bool {
	_, err := os.Stat(r.repomdPath())
	return err == nil
}

func (r *yumRepository) ListPackages() ([]repos.Package, error) {
	files, err := ioutil.ReadDir(r.path())
	if err != nil {
		return nil, err
	}

	packages := make([]repos.Package, 0, len(files))

	for _, file := range files {
		if !file.IsDir() {
			p := repos.Package{
				Name:      file.Name(),
				Timestamp: file.ModTime().Format(time.RFC3339),
				Size:      file.Size(),
			}
			packages = append(packages, p)
		}
	}

	return packages, nil
}

func (r *yumRepository) PackageWithNameAndVersion(name string, version string) (repos.Package, error) {
	// get filelists name from repomd.xml
	repomdXML, err := os.Open(r.repomdPath())
	if err != nil {
		return repos.Package{}, err
	}
	defer repomdXML.Close()

	repomdData, err := ioutil.ReadAll(repomdXML)
	if err != nil {
		return repos.Package{}, err
	}
	var repomd RepoMd
	xml.Unmarshal(repomdData, &repomd)

	filelist := filepath.Join(r.path(), filelist(repomd))

	// read filelists.xml.gz
	f, err := os.Open(filelist)
	if err != nil {
		return repos.Package{}, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return repos.Package{}, err
	}
	defer gz.Close()

	filelistData, err := ioutil.ReadAll(gz)
	if err != nil {
		return repos.Package{}, err
	}

	var filelistsContent Filelists

	xml.Unmarshal(filelistData, &filelistsContent)

	return findPackageWithVersion(filelistsContent, name, version)
}

func findPackageWithVersion(filelist Filelists, name string, version string) (repos.Package, error) {
	for _, pkg := range filelist.Packages {
		v := pkg.Version.Version + "-" + pkg.Version.Rel
		if pkg.Name == name && v == version {
			return repos.Package{
				Name:    pkg.Name,
				Version: v,
				Arch:    pkg.Arch,
			}, nil
		}
	}
	return repos.Package{}, errors.New("no matching package with given version")
}

func filelist(repomd RepoMd) (filelist string) {
	for _, data := range repomd.Data {
		if data.Type == "filelists" {
			return data.Location.Href
		}
	}
	return ""
}

// structs for repomd.xml
type RepoMd struct {
	XMLName xml.Name `xml:"repomd"`
	Data    []Data   `xml:"data"`
}

type Data struct {
	XMLName  xml.Name `xml:"data"`
	Type     string   `xml:"type,attr"`
	Location Location `xml:"location"`
}

type Location struct {
	XMLName xml.Name `xml:"location"`
	Href    string   `xml:"href,attr"`
}

// structs for filelists.xml
type Filelists struct {
	XMLName  xml.Name          `xml:"filelists"`
	Packages []FilelistPackage `xml:"package"`
}

type FilelistPackage struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`
	Arch    string   `xml:"arch,attr"`
	Version Version  `xml:"version"`
}

type Version struct {
	XMLName xml.Name `xml:"version"`
	Rel     string   `xml:"rel,attr"`
	Version string   `xml:"ver,attr"`
}
