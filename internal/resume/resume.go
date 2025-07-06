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

func ClearState() error {
    statePath := filepath.Join(getStateDir(), "conversion_state.json")
    return os.Remove(statePath)
}

func getStateDir() string {
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, ".gopix", "state")
}