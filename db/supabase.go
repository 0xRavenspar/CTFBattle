// db/supabase.go
package db

import (
	"os"
	"sync"

	"github.com/supabase-community/supabase-go"
)

var (
	client *supabase.Client
	once   sync.Once
)

// GetClient returns a singleton instance of the Supabase client
func GetClient() *supabase.Client {
	once.Do(func() {
		supabaseUrl := os.Getenv("SUPABASE_URL")
		supabaseKey := os.Getenv("SUPABASE_ANON_KEY")
		var err error
		client, err = supabase.NewClient(supabaseUrl, supabaseKey, nil)
		if err != nil {
			panic(err)
		}
	})
	return client
}

// GetAdminClient returns a new client with service role key for admin operations
func GetAdminClient() (*supabase.Client, error) {
	supabaseUrl := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	return supabase.NewClient(supabaseUrl, serviceKey, nil)
}
