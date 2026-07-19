package main

import "testing"

func TestCleanBody(t *testing.T) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "no profanity",
			body: "This is a clean opinion",
			want: "This is a clean opinion",
		},
		{
			name: "lowercase profanity censored",
			body: "This is a kerfuffle opinion",
			want: "This is a **** opinion",
		},
		{
			name: "mixed case profanity censored",
			body: "This is a Kerfuffle, SHARBERT and fornax opinion",
			want: "This is a Kerfuffle, **** and **** opinion",
		},
		{
			name: "punctuation-attached word not censored",
			body: "Sharbert! is fine",
			want: "Sharbert! is fine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanBody(tt.body, badWords)
			if got != tt.want {
				t.Errorf("cleanBody(%q) = %q, want %q", tt.body, got, tt.want)
			}
		})
	}
}
