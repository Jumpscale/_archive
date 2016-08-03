package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/kothar/brotli-go.v0/enc"

	"github.com/Jumpscale/aydostorex/utils"
)

// Store is the structure that holds the logic for storing files
type Store struct {
	Root  string //path to the store folder
	Tmp   string //path to tmp directroty in store folder
	Fixed string //path to directory containing fixed path files
}

// NewStore creates a new instance of a store
func NewStore(root string) *Store {
	if err := os.MkdirAll(root, 0770); err != nil {
		log.Fatalf("Error creating root directroty for store :%v", err)
	}

	// create tmp dir in root
	tmpDir := filepath.Join(root, "tmp")
	if err := os.MkdirAll(tmpDir, 0770); err != nil {
		log.Fatalf("Error creating tmp directroty in root dir for store :%v", err)
	}

	fixDir := filepath.Join(root, "fixed")
	if err := os.MkdirAll(fixDir, 0770); err != nil {
		log.Fatalf("Error creating fixed directroty in root dir for store :%v", err)
	}

	return &Store{
		Root:  root,
		Tmp:   tmpDir,
		Fixed: fixDir,
	}
}

type File struct {
	Hash    string
	Size    int64
	Created time.Time
	Reader  io.Reader
}

func NewFile(info os.FileInfo, r io.ReadCloser) *File {

	return &File{
		Hash:    info.Name(),
		Size:    info.Size(),
		Created: info.ModTime(),
		Reader:  r,
	}
}

func (f *File) String() string {
	return fmt.Sprintf("%s|%d|%d", f.Hash, f.Size, f.Created.Unix())
}

//absolute return the absolute path of file in the store
func (s *Store) absolute(hash, namespace string) string {
	return filepath.Join(s.Root, namespace, string(hash[0]), string(hash[1]), hash)
}

func (s *Store) Put(r io.Reader, namespace string) (string, error) {

	tmpFile, err := ioutil.TempFile(s.Tmp, "storx")
	defer tmpFile.Close()
	if err != nil {
		log.Errorf("Put File, error getting temporary file : %v", err)
		return "", err
	}

	if _, err := io.Copy(tmpFile, r); err != nil {
		log.Errorf("Error writing file to temp destination (%v) : %v", tmpFile.Name(), err)
		return "", err
	}

	hash, err := utils.Hash(tmpFile)
	if err != nil {
		return "", err
	}

	path := s.absolute(hash, namespace)

	if _, err := os.Stat(path); err == nil {
		log.Debug("file %s already exists", hash)
		os.Remove(tmpFile.Name())
		return hash, nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0771); err != nil {
		log.Errorf("Put File, error creating parents directorties :%v", err)
		return "", err
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		log.Errorf("Put file, error renaming temp file %v to %v :%v", tmpFile.Name(), path, err)
		return "", err
	}

	return hash, nil
}

func (s *Store) Get(hash string, namespace string) (io.Reader, int64, error) {
	path := s.absolute(hash, namespace)
	if !s.Exists(hash, namespace) {
		return nil, -1, os.ErrNotExist
	}

	f, err := os.Open(path)
	if err != nil {
		log.Errorf("Error opening file %s: %v", path, err)
		return nil, -1, err
	}

	stat, err := f.Stat()
	if err != nil {
		log.Errorf("Error stat file %s: %v", path, err)
		return nil, -1, err
	}

	return f, stat.Size(), nil
}

func (s *Store) Delete(hash string, namespace string) error {
	path := s.absolute(hash, namespace)

	if !s.Exists(hash, namespace) {
		return os.ErrNotExist
	}

	err := os.Remove(path)
	if err != nil {
		log.Errorf("Error deleting file %s: %v", path, err)
		return err
	}

	return nil
}

func (s *Store) Exists(hash string, namespace string) bool {
	res := s.ExistList([]string{hash}, namespace)
	return res[hash]
}

func (s *Store) ExistList(hashes []string, namespace string) map[string]bool {
	results := make(map[string]bool, len(hashes))
	for _, hash := range hashes {
		path := s.absolute(hash, namespace)
		if _, err := os.Stat(path); err != nil {
			results[hash] = false
		} else {
			results[hash] = true
		}
	}
	return results
}

//List returns a sorted list of all hashes & sizes & creation dates of the files located in namespace
func (s *Store) List(namespace string, compress bool, quality int) ([]string, error) {
	list := make([]string, 0, 100)

	compressParams := enc.NewBrotliParams()
	if quality < 0 || 11 < quality {
		compressParams.SetQuality(6)
	} else {
		compressParams.SetQuality(quality)
	}

	rootNamespace := filepath.Join(s.Root, namespace)

	root, err := os.Open(rootNamespace)
	if err != nil {
		return nil, err
	}
	defer root.Close()

	subDirs, err := root.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	for _, dir := range subDirs {
		dir1Path := filepath.Join(rootNamespace, dir)
		dir1, err := os.Open(dir1Path)
		if err != nil {
			return nil, err
		}
		defer dir1.Close()

		dirs2, err := dir1.Readdirnames(-1)
		if err != nil {
			return nil, err
		}

		for _, dir2Name := range dirs2 {

			dir2Path := filepath.Join(dir1Path, dir2Name)
			dir2, err := os.Open(dir2Path)
			if err != nil {
				return nil, err
			}
			defer dir2.Close()

			for {

				infos, err := dir2.Readdir(2000)
				if err != nil {
					if err == io.EOF {
						//we reach the end for the directory
						break
					}
					log.Errorf("Error readdir %s :%v", dir2.Name(), err)
					return nil, err
				}

				// literal function here so the defer can occurs more often and close file as we read the directory.
				err = func(list *[]string, infos []os.FileInfo) error {
					for _, info := range infos {

						if info.IsDir() || strings.HasSuffix(info.Name(), ".bro") {
							continue
						}

						path := filepath.Join(dir2Path, info.Name())
						f, err := os.Open(path)
						defer f.Close()
						if err != nil {
							return err
						}

						if compress {
							fBro, err := os.OpenFile(path+".bro", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
							defer f.Close()
							if err != nil {
								return err
							}

							brotliWriter := enc.NewBrotliWriter(compressParams, fBro)
							defer brotliWriter.Close()

							log.Debug("start compression of  %v", path)
							_, err = io.Copy(brotliWriter, f)
							if err != nil {
								return err
							}
						}

						file := NewFile(info, f)
						*list = append(*list, file.String())
					}
					return nil
				}(&list, infos)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	sort.StringSlice(list).Sort()

	return list, nil
}

func (s *Store) PutWithName(r io.Reader, name string) error {
	destPath := filepath.Join(s.Fixed, name)
	destFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		log.Errorf("Error opening file (%s) for write: %v", destPath, err)
		return err
	}

	if _, err := io.Copy(destFile, r); err != nil {
		log.Errorf("Error writing file to %v : %v", destFile.Name(), err)
		return err
	}

	return nil
}

func (s *Store) GetWithName(name string) (io.ReadCloser, int64, error) {
	path := filepath.Join(s.Fixed, name)
	f, err := os.Open(path)
	if err != nil {
		log.Debug("Error opening file (%s) for read: %v", path, err)
		return nil, -1, err
	}

	stat, err := f.Stat()
	if err != nil {
		log.Errorf("Error stat file %s: %v", path, err)
		return nil, -1, err
	}

	return f, stat.Size(), nil
}
