// main.go
package main

import (
	// "CTFBattle/middleware"
	// "CTFBattle/services/auth"
	"CTFBattle/services/rooms"
	"CTFBattle/services/user_rooms"
	"CTFBattle/services/users"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// User routes
	// app.POST("/users/add", middleware.AuthMiddleware(users.AddUserHandler))
	// app.GET("/users/{email}", middleware.AuthMiddleware(users.GetUserHandler))
	// app.DELETE("/users/{email}", middleware.AuthMiddleware(users.DeleteUserHandler))

	app.GET("/", func(c *gofr.Context) (interface{}, error) {
		return "Hello, World!", nil
	})
	app.POST("/users/add", users.AddUserHandler)
	app.GET("/users/{email}", users.GetUserHandler)
	app.DELETE("/users/{email}", users.DeleteUserHandler)

	// Protected routes
	app.POST("/rooms/create", rooms.CreateRoomHandler)
	app.GET("/rooms/{id}", rooms.GetRoomDetailsHandler)
	// app.POST("/ctfd/create",ctfd.CreateCTFdInstance)

	app.POST("/rooms/join/{roomid}", user_rooms.JoinRoomHandler)
	app.GET("/rooms/user-stats/{userid}/{roomid}", user_rooms.GetUserStatsHandler)
	app.GET("/rooms/stats/{roomid}", user_rooms.GetRoomStatsHandler)

	app.Run()
}
