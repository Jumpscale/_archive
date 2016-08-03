package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
)

func (r *Service) handlerPut(c *gin.Context) {

	namespace := c.Param("namespace")
	if namespace == "" {
		err := fmt.Errorf("No namespace specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer c.Request.Body.Close()
	hash, err := r.store.Put(c.Request.Body, namespace)
	if err != nil {
		log.Errorf("Error saving file to store: %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	data := struct{ Hash string }{hash}
	c.JSON(http.StatusCreated, data)
}

func (r *Service) handlerGet(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		err := fmt.Errorf("No namespace specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hash := c.Param("hash")
	if hash == "" {
		err := fmt.Errorf("No hash specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	reader, size, err := r.store.Get(hash, namespace)
	if err != nil {
		ctx, cancel := context.WithCancel(context.Background())
		reader, err = r.tryBackendStores(ctx, c.Request)
		if handlerError(err, c) != nil {
			return
		}
		cancel()

		if rc, ok := reader.(io.Closer); ok {
			defer rc.Close()
		}

		hash, err := r.store.Put(reader, namespace)
		if err != nil {
			log.Errorf("Error saving file to store: %v\n", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		reader, size, err = r.store.Get(hash, namespace)
		if handlerError(err, c) != nil {
			return
		}
	}

	if rc, ok := reader.(io.Closer); ok {
		defer rc.Close()
	}

	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	c.Writer.WriteHeader(http.StatusOK)
	io.Copy(c.Writer, reader)
}

func (r *Service) tryBackendStores(ctx context.Context, req *http.Request) (io.Reader, error) {
	result := make(chan io.Reader)
	wait := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(r.backendStores))

	for _, u := range r.backendStores {
		go func(u *url.URL, out chan io.Reader) {
			newURL := rewriteURL(req.URL, u)
			log.Debug("Try finding file at %s", newURL.String())

			var err error
			var resp *http.Response
			req.URL, err = url.Parse(newURL.String())
			if err == nil {
				req.RequestURI = "" // RequestURI has to be empty when passing request to client. data is in URL anyway
				resp, err = http.DefaultClient.Do(req)
				if err == nil && resp.StatusCode != http.StatusOK {
					err = fmt.Errorf("response status code not valid %d", resp.StatusCode)
				}
			}

			if err == nil {
				select {
				case out <- resp.Body:
				default:
					f, ok := resp.Body.(io.ReadCloser)
					if ok {
						log.Debug("Closing unused file")
						f.Close()
					}
				}
			}
			wg.Done()
		}(u, result)
	}
	go func() {
		wg.Wait()
		wait <- 1
	}()

	select {
	case r := <-result:
		if r == nil {
			return nil, os.ErrNotExist
		}
		return r, nil
	case <-wait:
		//all exited with no response.
		return nil, os.ErrNotExist
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func rewriteURL(src, dest *url.URL) *url.URL {
	src.Scheme = dest.Scheme
	src.Host = dest.Host
	i := strings.Index(src.Path, "store")
	if i < 0 {
		i = 0
	}
	src.Path = path.Join(dest.Path, src.Path[i:])

	return src
}

func (r *Service) handlerExists(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		err := fmt.Errorf("No namespace specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hash := c.Param("hash")
	if hash == "" {
		err := fmt.Errorf("No hash specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	found := r.store.Exists(hash, namespace)
	if found {
		c.Writer.WriteHeader(http.StatusOK)
		return
	}
	c.Writer.WriteHeader(http.StatusNotFound)
}

func (r *Service) handlerExistsList(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		err := fmt.Errorf("No namespace specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	decoder := json.NewDecoder(c.Request.Body)
	var post struct {
		Hashes []string
	}

	err := decoder.Decode(&post)
	if err != nil || post.Hashes == nil {
		err := fmt.Errorf("No hashlist found\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	response := make(map[string]bool)

	for _, item := range post.Hashes {
		log.Info(item)
		response[item] = r.store.Exists(item, namespace)
	}

	c.JSON(http.StatusOK, response)
}

func (r *Service) handlerDelete(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		err := fmt.Errorf("No namespace specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	hash := c.Param("hash")
	if hash == "" {
		err := fmt.Errorf("No hash specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	err := r.store.Delete(hash, namespace)
	if handlerError(err, c) != nil {
		log.Errorf("Error deleting file from store: %v\n", err)
		return
	}

	data := struct{ deleted bool }{true}
	c.JSON(http.StatusOK, data)
}

func (r *Service) handlerNamespaceList(c *gin.Context) {
	namespace := c.Param("namespace")
	if namespace == "" {
		err := fmt.Errorf("No namespace specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	compress := false
	val := c.Request.URL.Query().Get("compress")
	val = strings.ToLower(val)
	if val == "true" {
		compress = true
	}

	val = c.Request.URL.Query().Get("quality")
	quality, err := strconv.Atoi(val)
	if err != nil {
		quality = -1
	}

	log.Debugf("List Namespace. Compression:%v quality:%v", compress, quality)

	list, err := r.store.List(namespace, compress, quality)
	if handlerError(err, c) != nil {
		return
	}

	c.JSON(http.StatusOK, list)
}

func (r *Service) handlerPutWithName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		err := fmt.Errorf("No name specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer c.Request.Body.Close()
	err := r.store.PutWithName(c.Request.Body, name)
	if err != nil {
		log.Errorf("Error saving file to store: %v\n", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.WriteHeaderNow()
}

func (r *Service) handlerGettWithName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		err := fmt.Errorf("No name specified\n")
		log.Error(err.Error())
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	defer c.Request.Body.Close()
	rc, size, err := r.store.GetWithName(name)
	if handlerError(err, c) != nil {
		return
	}
	defer rc.Close()

	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	c.Writer.WriteHeader(http.StatusOK)
	io.Copy(c.Writer, rc)
}

func handlerError(err error, c *gin.Context) error {
	if err != nil {
		status := http.StatusInternalServerError
		if os.IsNotExist(err) {
			status = http.StatusNotFound
		}
		c.AbortWithError(status, err)
	}
	return err
}
