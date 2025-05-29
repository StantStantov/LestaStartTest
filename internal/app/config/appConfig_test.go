package config_test

import (
	"Stant/LestaGamesInternship/internal/app/config"
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	t.Run("Test reading server port", func(t *testing.T) {
		testReadServerPort(t)
	})
}

func testReadServerPort(t *testing.T) {
	t.Run("Pass if port is correct", func(t *testing.T) {
		wantPort := "9090"

		err := os.Setenv("SERVER_PORT", wantPort)
		if err != nil {
			t.Fatal(err)
		}

		config, err := config.ReadAppConfig()
		if err != nil {
			t.Fatal(err)
		}
		gotPort := config.ServerPort()

		if wantPort != gotPort {
			t.Errorf("Wanted %s, got %s", wantPort, gotPort)
		}
	})
	t.Run("Fail if port not set", func(t *testing.T) {
		os.Unsetenv("SERVER_PORT")

		_, err := config.ReadAppConfig()
		if err == nil {
			t.Fatalf("Wanted err, got nil")
		}
	})
	t.Run("Fail if port is incorrect", func(t *testing.T) {
		ports := []string{"", "2000000", "tseta"}

		for _, port := range ports {
			err := os.Setenv("SERVER_PORT", port)
			if err != nil {
				t.Fatal(err)
			}

			_, err = config.ReadAppConfig()
			if err == nil {
				t.Fatalf("Wanted err, got nil with %s", port)
			}
		}
	})
}
