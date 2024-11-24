package user_rooms

import (
	"CTFBattle/db"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"gofr.dev/pkg/gofr"
)

type User_Rooms struct {
	UserID    string    `json:"userid"`
	RoomID    string    `json:"roomid"`
	JoinedAt  time.Time `json:"joined_at"`
}

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

type User struct {
	UserID    string    `json:"userid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type createUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

type createUserResponse struct {
	Password string `json:"password"`
}

// JoinRoom adds a user to a room
func JoinRoom(c *gofr.Context, userRoom *User_Rooms) error {
	client := db.GetClient()
	// Set the join time
	userRoom.JoinedAt = time.Now()

	// Prepare the data to insert
	data := map[string]interface{}{
		"userid":    userRoom.UserID,
		"roomid":    userRoom.RoomID,
		"joined_at": userRoom.JoinedAt,
	}

	// Insert the data into the "user_rooms" table
	_, _, err := client.From("user_rooms").Insert(data, false, "", "", "").Execute()
	if err != nil {
		return err
	}

	return nil
}

// GetUserStats retrieves stats for a specific user in a room
func GetUserStats(userID string, roomID string) (*User_Rooms, error) {
	client := db.GetClient()
	var userRoom User_Rooms

	// Query the "user_rooms" table for the user in the specified room
	result, _, err := client.From("user_rooms").
		Select("*", "", false).
		Eq("userid", userID).
		Eq("roomid", roomID).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	// Parse the result into User_Rooms struct
	err = json.Unmarshal(result, &userRoom)
	if err != nil {
		return nil, err
	}

	return &userRoom, nil
}

// GetRoomStats retrieves overall stats for a specific room
func GetRoomStats(c *gofr.Context, roomID string) (map[string]interface{}, error) {
	client := db.GetClient()

	// Count the number of users in the room
	countResult, _, err := client.From("user_rooms").
		Select("COUNT(*)", "", false).
		Eq("roomid", roomID).
		Execute()
	if err != nil {
		return nil, err
	}

	var count struct {
		Count int `json:"count"`
	}
	err = json.Unmarshal(countResult, &count)
	if err != nil {
		return nil, err
	}

	// Query the highest score in the room
	highScoreResult, _, err := client.From("user_rooms").
		Select("MAX(score)", "", false).
		Eq("roomid", roomID).
		Execute()
	if err != nil {
		return nil, err
	}

	var highScore struct {
		Score string `json:"max"`
	}
	err = json.Unmarshal(highScoreResult, &highScore)
	if err != nil {
		return nil, err
	}

	// Make an API call for room-specific data (e.g., number of challenges)
	apiURL := "https://external-api.com/rooms/" + roomID + "/stats"
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, errors.New("failed to fetch room data from API")
	}
	defer resp.Body.Close()

	var apiData struct {
		Challenges int `json:"challenges"`
	}
	err = json.NewDecoder(resp.Body).Decode(&apiData)
	if err != nil {
		return nil, errors.New("failed to decode API response")
	}

	// Combine the results
	stats := map[string]interface{}{
		"num_users":      count.Count,
		"highest_score":  highScore.Score,
		"num_challenges": apiData.Challenges,
	}

	return stats, nil
}

// JoinRoomHandler handles the endpoint to join a room
func JoinRoomHandler(ctx *gofr.Context) (interface{}, error) {
	roomID := ctx.PathParam("roomid") // Updated to match new parameter name
	if roomID == "" {
		ctx.Error(http.StatusBadRequest, "Room ID is required")
		return nil, errors.New("Room ID is required")
	}

	var user User_Rooms
	if err := ctx.Bind(&user); err != nil {
		ctx.Error(http.StatusBadRequest, "Invalid request body")
		return nil, err
	}

	client := db.GetClient()

	var room Room
	res, _, err := client.From("rooms").
		Select("ctfdurl", "", false).  // Added level to select if needed
		Eq("roomid", roomID).
		Single().
		Execute()
	if err != nil {
		ctx.Error(http.StatusNotFound, "Room not found")
		return nil, err
	}
	
	// Parse the response into the room struct
	if err := json.Unmarshal(res, &room); err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to parse room data")
		return nil, err
	}
	
	if room.Ctfdurl == "" {
		ctx.Error(http.StatusNotFound, "Room ctfdurl not found")
		return nil, fmt.Errorf("Room ctfdurl not found, %s", room.Level)
	}
	

	// Fetch the user details (username, email) based on the user ID from the users table
	var userDetails User
	rest,_, err := client.From("users").
		Select("name, email", "", false).
		Eq("userid", user.UserID).
		Single().
		Execute()
	if err != nil {
		ctx.Error(http.StatusNotFound, "User not found")
		return nil, err
	}
	
	// Parse the response into the userDetails struct
	if err := json.Unmarshal(rest, &userDetails); err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to parse user data")
		return nil, err
	}
	
	// Optional: Check if required fields are present
	if userDetails.Name == "" || userDetails.Email == "" {
		ctx.Error(http.StatusNotFound, "User details incomplete")
		return nil, fmt.Errorf("User details missing required fields")
	}

	// Call the external API to generate the password using username and email
	password, err := generatePassword(userDetails.Name, userDetails.Email, room.Ctfdurl)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to generate password")
		return nil, err
	}

	// Set user-room properties
	user.RoomID = roomID
	user.JoinedAt = time.Now()

	// Insert data into the "user_rooms" table
	data := map[string]interface{}{
		"userid":    user.UserID,
		"roomid":    user.RoomID,
		"joined_at": user.JoinedAt,
	}

	_, _, err = client.From("user_rooms").Insert(data, false, "", "", "").Execute()
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to join room")
		return nil, err
	}

	// Return the response with the ctfdurl and generated password
	return map[string]interface{}{
		"message":   "User joined room successfully",
		"ctfdurl":   room.Ctfdurl,
		"password":  password, // Include the generated password
	}, nil
}

func generatePassword(name, email, ctfdurl string) (string, error) {
	// Prepare the payload
	payload := createUserPayload{
		Name:  name,
		Email: email,
		URL:   fmt.Sprintf("%s/api/v1/users", ctfdurl), // Append the required endpoint
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Make the POST request to the external API
	resp, err := http.Post("http://localhost:5000/create_user", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call /create_user API: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("non-200 response from /create_user API: %s", body)
	}

	// Parse the API response
	var apiResponse createUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return "", fmt.Errorf("failed to decode /create_user API response: %w", err)
	}

	// Return the password
	return apiResponse.Password, nil
}


// GetUserStatsHandler handles the endpoint to get user stats in a room
func GetUserStatsHandler(ctx *gofr.Context) (interface{}, error) {
	userID := ctx.PathParam("userid")
	roomID := ctx.PathParam("roomid")

	userStats, err := GetUserStats(userID, roomID)
	if err != nil {
		ctx.Error(http.StatusNotFound, "User stats not found")
		return nil, err
	}

	return userStats, nil
}

// GetRoomStatsHandler handles the endpoint to get room stats
func GetRoomStatsHandler(ctx *gofr.Context) (interface{}, error) {
	roomID := ctx.PathParam("roomid") // Updated to match new parameter name
	if roomID == "" {
		ctx.Error(http.StatusBadRequest, "Room ID is required")
		return nil, errors.New("Room ID is required")
	}

	client := db.GetClient()

	// Example query to calculate stats (adjust as per schema)
	stats := struct {
		NumberOfUsers int    `json:"number_of_users"`
		HighestScore  string `json:"highest_score"`
	}{
		NumberOfUsers: 0,
		HighestScore:  "0",
	}

	// Query database for number of users in the room
	countResult, _, err := client.From("user_rooms").
		Select("COUNT(userid)", "count", false).
		Eq("roomid", roomID).
		Single().
		Execute()
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to fetch stats")
		return nil, err
	}
	var count struct {
		Count int `json:"count"`
	}
	err = json.Unmarshal(countResult, &count)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to parse count result")
		return nil, err
	}
	stats.NumberOfUsers = count.Count // Parse count result

	// Query database for highest score in the room
	scoreResult, _, err := client.From("user_rooms").
		Select("MAX(score)", "max", false).
		Eq("roomid", roomID).
		Single().
		Execute()
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to fetch stats")
		return nil, err 
	}
	var highScore struct {
		Score string `json:"max"`
	}
	err = json.Unmarshal(scoreResult, &highScore)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, "Failed to parse score result")
		return nil, err
	}
	stats.HighestScore = highScore.Score

	return stats, nil
}