package service

import (
	"context"
	"database/sql"
	"max-odyssey-backend/internal/db"
	"max-odyssey-backend/internal/models"
)

type UserSkillService struct {
	db *db.Queries
}

func NewUserSkillService(db *db.Queries) *UserSkillService {
	return &UserSkillService{
		db: db,
	}
}

// GetUserSkills gets all skills for a user
func (s *UserSkillService) GetUserSkills(userID int64) ([]models.UserSkill, error) {
	dbSkills, err := s.db.GetUserSkills(context.Background(), int32(userID))
	if err != nil {
		return nil, err
	}

	skills := make([]models.UserSkill, len(dbSkills))
	for i, skill := range dbSkills {
		skills[i] = models.UserSkill{
			ID:         int(skill.ID),
			UserID:     int(skill.UserID),
			SkillName:  skill.SkillName,
			SkillScore: float64(skill.SkillScore),
			CreatedAt:  skill.CreatedAt.Time,
			UpdatedAt:  skill.UpdatedAt.Time,
		}
	}

	return skills, nil
}

// GetUserSkillByName gets a specific skill for a user
func (s *UserSkillService) GetUserSkillByName(userID int64, skillName string) (*models.UserSkill, error) {
	skill, err := s.db.GetUserSkillByName(context.Background(), db.GetUserSkillByNameParams{
		UserID:    int32(userID),
		SkillName: skillName,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Skill not found
		}
		return nil, err
	}

	return &models.UserSkill{
		ID:         int(skill.ID),
		UserID:     int(skill.UserID),
		SkillName:  skill.SkillName,
		SkillScore: float64(skill.SkillScore),
		CreatedAt:  skill.CreatedAt.Time,
		UpdatedAt:  skill.UpdatedAt.Time,
	}, nil
}

// CreateOrUpdateUserSkill creates or updates a skill for a user
func (s *UserSkillService) CreateOrUpdateUserSkill(userID int64, skillName string, skillScore float64) (*models.UserSkill, error) {
	skill, err := s.db.CreateUserSkill(context.Background(), db.CreateUserSkillParams{
		UserID:     int32(userID),
		SkillName:  skillName,
		SkillScore: float32(skillScore),
	})
	if err != nil {
		return nil, err
	}

	return &models.UserSkill{
		ID:         int(skill.ID),
		UserID:     int(skill.UserID),
		SkillName:  skill.SkillName,
		SkillScore: float64(skill.SkillScore),
		CreatedAt:  skill.CreatedAt.Time,
		UpdatedAt:  skill.UpdatedAt.Time,
	}, nil
}

// DeleteUserSkill deletes a skill for a user
func (s *UserSkillService) DeleteUserSkill(userID int64, skillName string) error {
	return s.db.DeleteUserSkill(context.Background(), db.DeleteUserSkillParams{
		UserID:    int32(userID),
		SkillName: skillName,
	})
}
