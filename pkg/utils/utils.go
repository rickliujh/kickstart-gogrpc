package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"google.golang.org/protobuf/types/known/anypb"
)

func ParseSlogLevel(s string) (slog.Level, error) {
	var level slog.Level
	err := level.UnmarshalText([]byte(s))
	return level, err
}

// packJSONIntoAny converts a JSON-serializable struct into an Any message
func PackJSONIntoAny(v any) (*anypb.Any, error) {
	// Convert the struct to JSON bytes
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	// Create Any message
	anyMsg := &anypb.Any{
		// Use a custom type URL for JSON content
		TypeUrl: "type.googleapis.com/json",
		// Store JSON bytes as the value
		Value: jsonBytes,
	}

	return anyMsg, nil
}

// unpackAnyToJSON extracts JSON data from an Any message and unmarshals it into the target struct
func UnpackAnyToJSON(anyMsg *anypb.Any, target any) error {
	// Verify type URL (optional, but recommended)
	if anyMsg.TypeUrl != "type.googleapis.com/json" {
		return fmt.Errorf("unexpected type URL: %s", anyMsg.TypeUrl)
	}

	// Unmarshal the value bytes into the target struct
	if err := json.Unmarshal(anyMsg.Value, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON from Any: %w", err)
	}

	return nil
}
