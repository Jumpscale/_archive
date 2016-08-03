package rest

import "github.com/gin-gonic/gin"

func Router(authPath string, restSrv *Service) *gin.Engine {
	r := gin.Default()
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	if authPath != "" {
		log.Info("Basic authentification enabled")
		r.Use(BasicAuth(authPath))
	}

	store := r.Group("/store")
	store.POST("/:namespace", restSrv.handlerPut)
	store.GET("/:namespace/:hash", restSrv.handlerGet)
	store.POST("/:namespace/exists", restSrv.handlerExistsList)
	store.HEAD("/:namespace/:hash", restSrv.handlerExists)
	store.DELETE("/:namespace/:hash", restSrv.handlerDelete)
	store.GET("/:namespace", restSrv.handlerNamespaceList)

	fixed := r.Group("/static")
	fixed.POST("/:name", restSrv.handlerPutWithName)
	fixed.GET("/:name", restSrv.handlerGettWithName)

	return r
}
