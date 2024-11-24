package rooms

import (
	"CTFBattle/db"
	"CTFBattle/services/ctfd"
	"encoding/json"
	"net/http"
	"time"

	"gofr.dev/pkg/gofr"
)

type Room struct {
	RoomID    string    `json:"roomid"`
	Name      string    `json:"name"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Ctfdurl   string    `json:"ctfdurl"`
	Threshold string    `json:"threshold"`
	ExpirationAt time.Time `json:"expiration_at"`
}

// type GenerateChallengesRequest struct {
// 	Level   string    `json:"level"`
// 	CtfdURL string `json:"ctfdurl"`
// }


// func GenerateChallenges(roomID string, level string, ctfdurl string) error {
// 	// Prepare the data to send to the challenge generation API
// 	requestData := GenerateChallengesRequest{
// 		Level:   level,
// 		CtfdURL: ctfdurl,
// 	}

// 	jsonData, err := json.Marshal(requestData)
//     if err != nil {
//         return fmt.Errorf("failed to marshal request body: %w", err)
//     }

//     // Create HTTP client with timeout
//     client := &http.Client{
//         Timeout: 3000000000 * time.Second,
//     }

//     // Make the POST request
//     resp, err := client.Post(
//         "http://localhost:5000/create_challenges",
//         "application/json",
//         bytes.NewBuffer(jsonData),
//     )
//     if err != nil {
//         return fmt.Errorf("failed to call create_challenges API: %w", err)
//     }
//     defer resp.Body.Close()

//     // Read the response body
//     body, err := io.ReadAll(resp.Body)
//     if err != nil {
//         return fmt.Errorf("failed to read response body: %w", err)
//     }

//     // Check status code first
//     if resp.StatusCode != http.StatusOK {
//         return fmt.Errorf("non-200 response from create_challenges API: %s", body)
//     }

//     return nil
// }


// CreateRoom adds a new room to the database
func CreateRoom(c *gofr.Context) (string, error) {
	client := db.GetClient()

	// Parse request body to get roomid and name
	var requestData struct {
		RoomID string `json:"roomid"`
		Name   string `json:"name"`
		Level  string `json:"level"`
	}

	err := c.Bind(&requestData)
	if err != nil {
		c.Error(http.StatusBadRequest, "Invalid request payload")
		return "",err
	}

	// Call the function to get the CTFd URL
	ctfdURL, err := ctfd.CreateCTFdContainer(requestData.RoomID) // Adjust parameters as needed
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to create CTFd container")
		return "",err
	}

	// Get the current time and calculate expiration time (3 hours later)
	createdAt := time.Now()
	expirationTime := createdAt.Add(3 * time.Hour)

	// Construct the Room struct
	room := Room{
		RoomID:       requestData.RoomID,
		Name:         requestData.Name,
		Level:        requestData.Level,
		CreatedAt:    createdAt,
		ExpirationAt: expirationTime, // Added expiration time
		Status:       "active",       // Fixed value
		Ctfdurl:      ctfdURL,        // From the external function
		Threshold:    "10",           // Fixed value
	}

	// Construct the data map
	data := map[string]interface{}{
		"roomid":      room.RoomID,
		"name":        room.Name,
		"level":       room.Level,
		"created_at":  room.CreatedAt,
		"status":      room.Status,
		"ctfdurl":     room.Ctfdurl,
		"threshold":   room.Threshold,
		"expiration_at": room.ExpirationAt, // Include expiration time in the data
	}

	// Insert the data into the "rooms" table
	_, _, err = client.From("rooms").Insert(data, false, "", "", "").Execute()
	if err != nil {
		c.Error(http.StatusInternalServerError, "Failed to insert room into database")
		return "",err
	}

	// err = GenerateChallenges(room.RoomID, room.Level, room.Ctfdurl)
	// if err != nil {
	// 	c.Error(http.StatusInternalServerError, "Failed to generate challenges")
	// 	return "", err
	// }

	// Return response with expiration time
	return room.ExpirationAt.Format(time.RFC3339), nil
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

	expiration_time, err := CreateRoom(ctx)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to create room")
		return nil, err
	}

	return map[string]string{"message": "Room created successfully", "expires_at": expiration_time}, nil
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