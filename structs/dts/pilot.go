package dts

type (
	QuizResponse struct {
		Base
		UserID               uint   `gorm:"column:user_id;index" json:"user_id"`
		User                 *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
		Skills               string `gorm:"type:jsonb" json:"skills"`                  // JSON: {"python": 4, "js": 2, ...}
		WorkingStyle         string `gorm:"column:working_style" json:"working_style"` // e.g., "independent", "collaborative"
		PreferredRoles       string `gorm:"type:jsonb" json:"preferred_roles"`         // JSON list: ["leader", "designer"]
		ScheduleAvailability string `gorm:"type:jsonb" json:"schedule_availability"`   // JSON: ["Mon AM", "Tue PM", ...]
	}

	Team struct {
		Base
		Name    string  `gorm:"column:name" json:"name"`
		Members []*User `gorm:"many2many:team_members;joinForeignKey:TeamID;JoinReferences:UserID" json:"members,omitempty"`
		Project string  `gorm:"column:project" json:"project"` // optional: title of the team’s project
	}

	Feedback struct {
		Base
		FromUserID uint   `gorm:"column:from_user_id;index" json:"from_user_id"`
		FromUser   *User  `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
		ToUserID   uint   `gorm:"column:to_user_id;index" json:"to_user_id"`
		ToUser     *User  `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
		TeamID     uint   `gorm:"column:team_id;index" json:"team_id"`
		Team       *Team  `gorm:"foreignKey:TeamID" json:"team,omitempty"`
		Rating     int    `gorm:"column:rating" json:"rating"`     // 1–5 scale
		Comments   string `gorm:"column:comments" json:"comments"` // open feedback
	}
)
