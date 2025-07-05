package project

// Project represents a project entity with its basic attributes.
// Used for creating, retrieving, and displaying project information.
type Project struct {
	Name        string `json:"name"`        // Name of the project
	Description string `json:"description"` // Description of the project
	Color       string `json:"color"`       // Color label for the project (e.g., for UI)
	User        int    `json:"user_id"`     // ID of the user who owns the project
}

// UpdateProjectRequest is used to update an existing project's details.
// Fields are optional and only provided values will be updated.
type UpdateProjectRequest struct {
	ProjectID   int    `json:"projectId"`   // ID of the project to update
	Name        string `json:"name"`        // New name for the project (optional)
	Description string `json:"description"` // New description (optional)
	Color       string `json:"color"`       // New color label (optional)
}

// IProjectService defines the interface for project-related business logic.
// Each method should be implemented to handle the corresponding project operation.
type IProjectService interface {
	// CreateProject creates a new project in the system.
	// Should accept project details and persist them to storage.
	CreateProject()
	// DeleteProject deletes a project by its ID.
	// Should remove the project from storage.
	DeleteProject()
	// GetProjectById retrieves a project by its ID.
	// Should return the project details if found.
	GetProjectById()
	// GetProjectsByUserId retrieves all projects for a given user.
	// Should return a list of projects owned by the user.
	GetProjectsByUserId()
	// UpdateProject updates the details of an existing project.
	// Should apply changes to the specified project.
	UpdateProject()
}
