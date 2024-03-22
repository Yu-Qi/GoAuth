package util

import "testing"

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Valid password",
			password: "Password123!",
			want:     true,
		},
		{
			name:     "Password too short",
			password: "Pass1",
			want:     false,
		},
		{
			name:     "Password too long",
			password: "Password123!Password123!",
			want:     false,
		},
		{
			name:     "Password missing uppercase",
			password: "password123!",
			want:     false,
		},
		{
			name:     "Password missing lowercase",
			password: "PASSWORD123!",
			want:     false,
		},
		{
			name:     "Password missing special character",
			password: "Password123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePassword(tt.password); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}

}
