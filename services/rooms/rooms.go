package rooms

import (
	"CTFBattle/db"
	"encoding/json"
	"net/http"
	"time"

	"gofr.dev/pkg/gofr"
)

type Room struct {
	RoomID    string    `json:"roomid"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Ctfdurl   string    `json:"ctfdurl"`
	Threshold string    `json:"threshold"`
}

// CreateRoom adds a new room to the database
func CreateRoom(c *gofr.Context, room *Room) error {
	client := db.GetClient()
	room.CreatedAt = time.Now()

	// Construct the data map
	data := map[string]interface{}{
		"roomid":     room.RoomID,
		"name":       room.Name,
		"created_at": room.CreatedAt,
		"status":     room.Status,
		"ctfdurl":    room.Ctfdurl,
		"threshold":  room.Threshold,
	}

	// Insert the data into the "rooms" table
	_, _, err := client.From("rooms").Insert(data,false, "", "", "" ).Execute()
	if err != nil {
		return err
	}

	return nil
}

// GetRoomDetails retrieves a room by room ID
func GetRoomDetails(roomID string) (*Room, error) {
	client := db.GetClient()
	var room Room

	// Query the "rooms" table by room ID
	result, _, err := client.From("rooms").
		Select("*", "", false).
		Eq("roomid", roomID).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	// Parse the result into the Room struct
	err = json.Unmarshal(result, &room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

// DeleteRoom deletes a room by room ID
func DeleteRoom(roomID string) error {
	client := db.GetClient()

	// Delete the room from the "rooms" table by room ID
	_, _, err := client.From("rooms").
		Delete("*", "").
		Eq("roomid", roomID).
		Execute()

	if err != nil {
		return err
	}

	return nil
}


// CreateRoomHandler handles the endpoint to create a room
func CreateRoomHandler(ctx *gofr.Context) (interface{}, error) {
	var room Room

	// Parse the request body into the room struct
	err := ctx.Bind(&room)
	if err != nil {
		ctx.Error(http.StatusBadRequest, "Invalid request body")
		return nil, err
	}

	err = CreateRoom(ctx, &room)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to create room")
		return nil, err
	}

	return map[string]string{"message": "Room created successfully"}, nil
}

// GetRoomDetailsHandler handles the endpoint to get room details by ID
func GetRoomDetailsHandler(ctx *gofr.Context) (interface{}, error) {
	roomID := ctx.PathParam("id")

	room, err := GetRoomDetails(roomID)
	if err != nil {
		ctx.Error(http.StatusNotFound, "Room not found")
		return nil, err
	}

	return room, nil
}