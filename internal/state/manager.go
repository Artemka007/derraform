package state

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type State struct {
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
	// Создаем директорию если нужно
	dir := filepath.Dir(sm.filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sm.filename, data, 0644)
}

func (sm *StateManager) SaveResourceState(resourceType, resourceName string, attributes map[string]interface{}) error {
	state, err := sm.Load()
	if err != nil {
		return err
	}

	resourceID := resourceType + "." + resourceName
	state.Resources[resourceID] = ResourceState{
		Type:       resourceType,
		ID:         attributes["id"].(string),
		Attributes: attributes,
	}

	return sm.Save(state)
}

func (sm *StateManager) Clear() error {
	state := &State{
		Resources: make(map[string]ResourceState),
	}
	return sm.Save(state)
}
