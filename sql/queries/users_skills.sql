-- name: GetUserSkills :many
SELECT * FROM users_skills
WHERE user_id = $1
ORDER BY skill_name;

-- name: GetUserSkillByName :one
SELECT * FROM users_skills
WHERE user_id = $1 AND skill_name = $2;

-- name: CreateUserSkill :one
INSERT INTO users_skills (
    user_id, skill_name, skill_score
) VALUES (
    $1, $2, $3
)
ON CONFLICT (user_id, skill_name) 
DO UPDATE SET
    skill_score = EXCLUDED.skill_score,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: UpdateUserSkill :one
UPDATE users_skills
SET 
    skill_score = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1 AND skill_name = $2
RETURNING *;

-- name: DeleteUserSkill :exec
DELETE FROM users_skills
WHERE user_id = $1 AND skill_name = $2; 