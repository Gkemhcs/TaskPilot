package project

import (
	"context"
	"database/sql"
	"errors"

	customErrors "github.com/Gkemhcs/taskpilot/internal/errors"
	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	"github.com/lib/pq"
)

func NewProjectService(projectRepository projectdb.Querier) *ProjectService {
	return &ProjectService{
		projectRepository: projectRepository,
	}
}

type ProjectService struct {
	projectRepository projectdb.Querier
}

func (p *ProjectService) CreateProject(ctx context.Context, project Project) (*projectdb.Project, error) {

	color := mapColor(project.Color)
	projectParams := projectdb.CreateProjectParams{
		UserID:      int32(project.User),
		Name:        project.Name,
		Description: sql.NullString{String: project.Description, Valid: project.Description != ""},
		Color:       projectdb.NullProjectColor{ProjectColor: color, Valid: true},
	}

	proj, err := p.projectRepository.CreateProject(ctx, projectParams)
	if  IsErrorCode(err, customErrors.UniqueViolationErr) {
		
				return nil,customErrors.ErrProjectAlreadyExists
			}
	if err != nil {
		return nil, err
	}
	return &proj, nil

}

func (p *ProjectService) DeleteProject(ctx context.Context, projectId int) error {

	return p.projectRepository.DeleteProject(ctx, int64(projectId))

}

func (p *ProjectService) GetProjectById(ctx context.Context, projectId int) (*projectdb.Project, error) {
	project, err := p.projectRepository.GetProjectById(ctx, int64(projectId))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrProjectIDNotExist
	}
	if err != nil {
		return nil, err
	}
	return &project, nil

}
func (p *ProjectService) GetProjectsByUserId(ctx context.Context, userId int) ([]projectdb.Project, error) {

	project, err := p.projectRepository.GetProjectsByUserId(ctx, int32(userId))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrProjectIDNotExist
	}
	if err != nil {
		return nil, err
	}
	return project, nil

}

func (p *ProjectService) GetProjectByName(ctx context.Context, name string, userID int) (*projectdb.Project, error) {
	params := projectdb.GetProjectByNameParams{
		Name:   name,
		UserID: int32(userID),
	}
	project, err := p.projectRepository.GetProjectByName(ctx, params)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, customErrors.ErrProjectNotExist
	}
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *ProjectService) UpdateProject(ctx context.Context, projectUpdateRequest UpdateProjectRequest) (*projectdb.Project, error) {
	var updateRequestParams projectdb.UpdateProjectParams
	updateRequestParams.ID = int64(projectUpdateRequest.ProjectID)
	if projectUpdateRequest.Name != "" {
		updateRequestParams.Name = projectUpdateRequest.Name
	}
	if projectUpdateRequest.Description != "" {
		updateRequestParams.Description = sql.NullString{String: projectUpdateRequest.Description, Valid: true}
	}
	if projectUpdateRequest.Color != "" {

		updateRequestParams.Color = projectdb.NullProjectColor{ProjectColor: mapColor(projectUpdateRequest.Color), Valid: true}
	}

	project, err := p.projectRepository.UpdateProject(ctx, updateRequestParams)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func mapColor(color string) projectdb.ProjectColor {
	switch color {
	case "green":
		return projectdb.ProjectColorGREEN
	case "yellow":
		return projectdb.ProjectColorYELLOW
	default:
		return projectdb.ProjectColorRED
	}
}

func IsErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
