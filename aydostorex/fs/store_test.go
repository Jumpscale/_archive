package fs

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Jumpscale/aydostorex/utils"
	"github.com/stretchr/testify/assert"
)

func initStore(t *testing.T) (*Store, func(), error) {
	path, err := ioutil.TempDir("", "aydostorex_testing")
	clean := func() {
		os.RemoveAll(path)
	}

	if err != nil {
		t.Error("Error initializing store:", err)
		return nil, clean, err
	}

	return NewStore(path), clean, err
}

func TestAbsolute(t *testing.T) {
	store, clean, err := initStore(t)
	defer clean()
	if !assert.NoError(t, err) {
		return
	}

	reader := bytes.NewReader([]byte("Hello World"))
	hash, err := utils.Hash(reader)
	if !assert.NoError(t, err) {
		return
	}

	actual := store.absolute(hash, "testing")
	expected := filepath.Join(store.Root, "testing", string(hash[0]), string(hash[1]), hash)
	assert.Equal(t, expected, actual)

}

func TestPutFile(t *testing.T) {
	store, clean, err := initStore(t)
	defer clean()
	if !assert.NoError(t, err) {
		return
	}

	content := "Hello World"

	reader := bytes.NewReader([]byte(content))
	hash, err := store.Put(reader, "testing")
	if !assert.NoError(t, err) {
		return
	}

	targetHash, err := utils.Hash(reader)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, targetHash, hash)

	b, err := ioutil.ReadFile(store.absolute(hash, "testing"))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, content, string(b))
}

func TestGetFile(t *testing.T) {
	store, clean, err := initStore(t)
	defer clean()
	if !assert.NoError(t, err) {
		return
	}

	content := "Hello World"
	reader := bytes.NewReader([]byte(content))
	hash, err := utils.Hash(reader)
	if !assert.NoError(t, err) {
		return
	}

	path := store.absolute(hash, "testing")
	os.MkdirAll(filepath.Dir(path), 0770)
	err = ioutil.WriteFile(path, []byte(content), 0660)
	if !assert.NoError(t, err) {
		return
	}

	resultReader, size, err := store.Get(hash, "testing")
	if !assert.NoError(t, err) {
		return
	}
	b, err := ioutil.ReadAll(resultReader)
	assert.Equal(t, content, string(b))
	assert.EqualValues(t, size, 11)
}

func TestDeleteFile(t *testing.T) {
	store, clean, err := initStore(t)
	defer clean()
	if !assert.NoError(t, err) {
		return
	}

	content := "Hello World"
	reader := bytes.NewReader([]byte(content))
	hash, err := utils.Hash(reader)
	if !assert.NoError(t, err) {
		return
	}

	path := store.absolute(hash, "testing")
	os.MkdirAll(filepath.Dir(path), 0770)
	err = ioutil.WriteFile(path, []byte(content), 0660)
	if !assert.NoError(t, err) {
		return
	}

	err = store.Delete(hash, "testing")
	if !assert.NoError(t, err) {
		return
	}
	_, err = os.Stat(store.absolute(hash, "testing"))
	assert.Error(t, err)
}

func TestExistsFile(t *testing.T) {
	store, clean, err := initStore(t)
	defer clean()
	if !assert.NoError(t, err) {
		return
	}

	reader := bytes.NewReader([]byte("Hello World"))
	hash, err := utils.Hash(reader)
	if assert.NoError(t, err) {
		return
	}
	exists := store.Exists(hash, "testing")
	assert.Equal(t, false, exists)

	path := store.absolute(hash, "testing")
	os.MkdirAll(filepath.Dir(path), 0770)
	ioutil.WriteFile(path, []byte("hello world"), 0660)
	exists = store.Exists(hash, "testing")
	assert.Equal(t, true, exists)
}
