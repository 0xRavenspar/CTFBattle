// main.go
package main

import (
	// "CTFBattle/middleware"
	// "CTFBattle/services/auth"
	"CTFBattle/services/rooms"
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

	app.Run()
}
