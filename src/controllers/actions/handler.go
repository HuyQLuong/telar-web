package function

import (
	"context"
	"fmt"
	"net/http"

	coreServer "github.com/red-gold/telar-core/server"
	"github.com/red-gold/telar-web/src/controllers"
	cf "github.com/red-gold/telar-web/src/controllers/actions/config"
	"github.com/red-gold/telar-web/src/controllers/actions/handlers"
)

func init() {

	cf.InitConfig()
}

// Cache state
var server *coreServer.ServerRouter
var db interface{}

// Handler function
func Handle(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	// Start
	if db == nil {
		var startErr error
		db, startErr = controllers.Start(ctx)
		if startErr != nil {
			fmt.Printf("Error startup: %s", startErr.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(startErr.Error()))
		}
	}

	// Server Routing
	if server == nil {
		server = coreServer.NewServerRouter()
		server.POST("/room", handlers.CreateActionRoomHandle(db), coreServer.RouteProtectionHMAC)
		server.POST("/dispatch/:roomId", handlers.DispatchHandle(db), coreServer.RouteProtectionHMAC)
		server.PUT("/room", handlers.UpdateActionRoomHandle(db), coreServer.RouteProtectionCookie)
		server.PUT("/room/access-key", handlers.SetAccessKeyHandle(db), coreServer.RouteProtectionCookie)
		server.DELETE("/room/:roomId", handlers.DeleteActionRoomHandle(db), coreServer.RouteProtectionHMAC)
		server.GET("/room/access-key", handlers.GetAccessKeyHandle(db), coreServer.RouteProtectionCookie)
		server.GET("/room/verify/:accessKey", handlers.VerifyAccessKeyHandle(db), coreServer.RouteProtectionCookie)
	}
	server.ServeHTTP(w, r)
}
