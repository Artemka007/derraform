package state

import (
    "encoding/json"
    "os"
)

type State struct {
    Version   string            `json:"version"`
    Resources map[string]ResourceState `json:"resources"`
}

type ResourceState struct {
    Type       string                 `json:"type"`
    ID         string                 `json:"id"`
    Attributes map[string]interface{} `json:"attributes"`
}

type StateManager struct {
    filename string
}

func NewStateManager(filename string) *StateManager {
    return &StateManager{filename: filename}
}

func (sm *StateManager) Load() (*State, error) {
    data, err := os.ReadFile(sm.filename)
    if err != nil {
        if os.IsNotExist(err) {
            return &State{
                Version:   "1.0.0",
                Resources: make(map[string]ResourceState),
            }, nil
        }
        return nil, err
    }

    var state State
    if err := json.Unmarshal(data, &state); err != nil {
        return nil, err
    }

    return &state, nil
}

func (sm *StateManager) Save(state *State) error {
    data, err := json.MarshalIndent(state, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(sm.filename, data, 0644)
}