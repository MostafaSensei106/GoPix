package resume

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type ConversionState struct {
	ProcessedFiles []string  `json:"processed_files"`
	StartTime      time.Time `json:"start_time"`
	InputDir       string    `json:"input_dir"`
	TargetFormat   string    `json:"target_format"`
	TotalFiles     int       `json:"total_files"`
	SessionID      string    `json:"session_id"`
}

// SaveState writes the current conversion state to a JSON file in the user's
// state directory. The file is named "conversion_state.json" and is used to
// resume the conversion process when the user restarts the program.
//
// The state is marshalled to JSON using the json.MarshalIndent function, which
// indents the JSON data with two spaces for readability. The resulting data is
// then written to the file using the os.WriteFile function.
//
// If any error occurs during the writing process, the function returns the error.
func SaveState(state *ConversionState) error {
	stateDir := getStateDir()
	if err := os.MkdirAll(stateDir, 0755); err != nil {
		return err
	}

	statePath := filepath.Join(stateDir, "conversion_state.json")
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statePath, data, 0644)
}

// LoadState reads the current conversion state from a JSON file in the user's
// state directory. The file is named "conversion_state.json" and is used to
// resume the conversion process when the user restarts the program.
//
// If the file does not exist, the function returns nil and no error. If any
// error occurs during the reading or unmarshalling process, the function
// returns the error.
func LoadState() (*ConversionState, error) {
	statePath := filepath.Join(getStateDir(), "conversion_state.json")

	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return nil, nil // No saved state
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil, err
	}

	var state ConversionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

// ClearState removes the saved conversion state file from the user's state directory.
//
// This is typically used after a successful conversion to remove the saved state
// and prevent the user from resuming the conversion again.
//
// If the file does not exist, the function returns nil and no error. If any
// error occurs during the removal process, the function returns the error.
func ClearState() error {
	statePath := filepath.Join(getStateDir(), "conversion_state.json")
	return os.Remove(statePath)
}

// getStateDir returns the path to the state directory where conversion state files are saved.
// The directory is located in the user's home directory and is named ".gopix/state".
func getStateDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".gopix", "state")
}
