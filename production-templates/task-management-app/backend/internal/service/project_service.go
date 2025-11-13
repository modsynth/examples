package service

import (
	"errors"
	"fmt"

	"task-management-app/internal/domain"
	"task-management-app/internal/repository"
)

type ProjectService interface {
	Create(userID uint, req *domain.CreateProjectRequest) (*domain.Project, error)
	GetByID(projectID, userID uint) (*domain.Project, error)
	Update(projectID, userID uint, req *domain.UpdateProjectRequest) (*domain.Project, error)
	Delete(projectID, userID uint) error
	Archive(projectID, userID uint) error
	Unarchive(projectID, userID uint) error
	ListUserProjects(userID uint) ([]*domain.Project, error)

	AddMember(projectID, userID uint, req *domain.AddMemberRequest) error
	RemoveMember(projectID, memberUserID, requestUserID uint) error
	UpdateMemberRole(projectID, memberUserID, requestUserID uint, req *domain.UpdateMemberRoleRequest) error
	GetMembers(projectID, userID uint) ([]domain.ProjectMember, error)

	CheckAccess(projectID, userID uint, requiredRole domain.ProjectRole) (bool, error)
	GetUserRole(projectID, userID uint) (domain.ProjectRole, error)
}

type projectService struct {
	projectRepo repository.ProjectRepository
	userRepo    repository.UserRepository
}

func NewProjectService(projectRepo repository.ProjectRepository, userRepo repository.UserRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

func (s *projectService) Create(userID uint, req *domain.CreateProjectRequest) (*domain.Project, error) {
	if req.Name == "" {
		return nil, errors.New("project name is required")
	}

	// Verify user exists
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	project := &domain.Project{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Color:       req.Color,
		OwnerID:     userID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	// Add owner as a member with owner role
	member := &domain.ProjectMember{
		ProjectID: project.ID,
		UserID:    userID,
		Role:      domain.ProjectRoleOwner,
	}
	if err := s.projectRepo.AddMember(member); err != nil {
		return nil, fmt.Errorf("failed to add owner as member: %w", err)
	}

	// Reload project with members
	return s.projectRepo.FindByID(project.ID)
}

func (s *projectService) GetByID(projectID, userID uint) (*domain.Project, error) {
	// Check if user has access to this project
	hasAccess, err := s.CheckAccess(projectID, userID, domain.ProjectRoleViewer)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("access denied to this project")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return project, nil
}

func (s *projectService) Update(projectID, userID uint, req *domain.UpdateProjectRequest) (*domain.Project, error) {
	// Only admin and owner can update project
	hasAccess, err := s.CheckAccess(projectID, userID, domain.ProjectRoleAdmin)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("insufficient permissions to update project")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Update fields if provided
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Icon != "" {
		project.Icon = req.Icon
	}
	if req.Color != "" {
		project.Color = req.Color
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

func (s *projectService) Delete(projectID, userID uint) error {
	// Only owner can delete project
	role, err := s.GetUserRole(projectID, userID)
	if err != nil {
		return err
	}
	if role != domain.ProjectRoleOwner {
		return errors.New("only project owner can delete the project")
	}

	if err := s.projectRepo.Delete(projectID); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

func (s *projectService) Archive(projectID, userID uint) error {
	// Only admin and owner can archive project
	hasAccess, err := s.CheckAccess(projectID, userID, domain.ProjectRoleAdmin)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("insufficient permissions to archive project")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	project.IsArchived = true
	if err := s.projectRepo.Update(project); err != nil {
		return fmt.Errorf("failed to archive project: %w", err)
	}

	return nil
}

func (s *projectService) Unarchive(projectID, userID uint) error {
	// Only admin and owner can unarchive project
	hasAccess, err := s.CheckAccess(projectID, userID, domain.ProjectRoleAdmin)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("insufficient permissions to unarchive project")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return fmt.Errorf("failed to get project: %w", err)
	}

	project.IsArchived = false
	if err := s.projectRepo.Update(project); err != nil {
		return fmt.Errorf("failed to unarchive project: %w", err)
	}

	return nil
}

func (s *projectService) ListUserProjects(userID uint) ([]*domain.Project, error) {
	projects, err := s.projectRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user projects: %w", err)
	}
	return projects, nil
}

func (s *projectService) AddMember(projectID, userID uint, req *domain.AddMemberRequest) error {
	// Only admin and owner can add members
	hasAccess, err := s.CheckAccess(projectID, userID, domain.ProjectRoleAdmin)
	if err != nil {
		return err
	}
	if !hasAccess {
		return errors.New("insufficient permissions to add members")
	}

	// Verify the user to be added exists
	_, err = s.userRepo.FindByID(req.UserID)
	if err != nil {
		return fmt.Errorf("user to add not found: %w", err)
	}

	// Check if user is already a member
	existingMember, _ := s.projectRepo.GetMember(projectID, req.UserID)
	if existingMember != nil {
		return errors.New("user is already a member of this project")
	}

	member := &domain.ProjectMember{
		ProjectID: projectID,
		UserID:    req.UserID,
		Role:      req.Role,
	}

	if err := s.projectRepo.AddMember(member); err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

func (s *projectService) RemoveMember(projectID, memberUserID, requestUserID uint) error {
	// Get the role of the user making the request
	requestUserRole, err := s.GetUserRole(projectID, requestUserID)
	if err != nil {
		return err
	}

	// Get the role of the member to be removed
	memberRole, err := s.GetUserRole(projectID, memberUserID)
	if err != nil {
		return err
	}

	// Owner cannot be removed
	if memberRole == domain.ProjectRoleOwner {
		return errors.New("project owner cannot be removed")
	}

	// Only admin and owner can remove members
	if !s.hasPermission(requestUserRole, domain.ProjectRoleAdmin) {
		// Members can remove themselves
		if requestUserID != memberUserID {
			return errors.New("insufficient permissions to remove members")
		}
	}

	if err := s.projectRepo.RemoveMember(projectID, memberUserID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

func (s *projectService) UpdateMemberRole(projectID, memberUserID, requestUserID uint, req *domain.UpdateMemberRoleRequest) error {
	// Only owner can change roles
	requestUserRole, err := s.GetUserRole(projectID, requestUserID)
	if err != nil {
		return err
	}
	if requestUserRole != domain.ProjectRoleOwner {
		return errors.New("only project owner can update member roles")
	}

	// Cannot change owner's role
	memberRole, err := s.GetUserRole(projectID, memberUserID)
	if err != nil {
		return err
	}
	if memberRole == domain.ProjectRoleOwner {
		return errors.New("cannot change project owner's role")
	}

	member, err := s.projectRepo.GetMember(projectID, memberUserID)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	member.Role = req.Role
	if err := s.projectRepo.UpdateMember(member); err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	return nil
}

func (s *projectService) GetMembers(projectID, userID uint) ([]domain.ProjectMember, error) {
	// Check if user has access to this project
	hasAccess, err := s.CheckAccess(projectID, userID, domain.ProjectRoleViewer)
	if err != nil {
		return nil, err
	}
	if !hasAccess {
		return nil, errors.New("access denied to this project")
	}

	members, err := s.projectRepo.GetMembers(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	return members, nil
}

func (s *projectService) CheckAccess(projectID, userID uint, requiredRole domain.ProjectRole) (bool, error) {
	userRole, err := s.GetUserRole(projectID, userID)
	if err != nil {
		return false, err
	}

	return s.hasPermission(userRole, requiredRole), nil
}

func (s *projectService) GetUserRole(projectID, userID uint) (domain.ProjectRole, error) {
	member, err := s.projectRepo.GetMember(projectID, userID)
	if err != nil {
		return "", fmt.Errorf("user is not a member of this project")
	}

	return member.Role, nil
}

// Helper function to check if userRole has at least the requiredRole
func (s *projectService) hasPermission(userRole, requiredRole domain.ProjectRole) bool {
	roleHierarchy := map[domain.ProjectRole]int{
		domain.ProjectRoleViewer: 1,
		domain.ProjectRoleMember: 2,
		domain.ProjectRoleAdmin:  3,
		domain.ProjectRoleOwner:  4,
	}

	return roleHierarchy[userRole] >= roleHierarchy[requiredRole]
}
