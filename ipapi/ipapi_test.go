package ipapi

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

var loc string

func setup() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	// loc equals Cyprus
	loc = os.Getenv("LOCATION")
}

func TestIsExpectedLocation(t *testing.T) {
	setup()

	type args struct {
		locationName string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name:  "success: location is the same",
			args:  args{locationName: loc},
			want:  loc,
			want1: true,
		},
		{
			name:  "fail: location is not the same, doesn't equal Cyprus",
			args:  args{locationName: "Some-Random-Country"},
			want:  loc,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IsExpectedLocation(tt.args.locationName)
			if got != tt.want {
				t.Errorf("IsExpectedLocation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IsExpectedLocation() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestLocation(t *testing.T) {
	setup()

	t.Run("success: got location Cyprus", func(t *testing.T) {
		got, err := Location()
		if (err != nil) != false {
			t.Errorf("Location() error = %v, wantErr %v", err, false)
			return
		}
		if got != loc {
			t.Errorf("Location() got = %v, want %v", got, loc)
		}
	})
}
