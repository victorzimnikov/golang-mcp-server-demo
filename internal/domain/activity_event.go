package domain

type ActivityEventEntityType string

const (
	ActivityEventEntityTypeProject  ActivityEventEntityType = "project"
	ActivityEventEntityTypeTask     ActivityEventEntityType = "task"
	ActivityEventEntityTypeDecision ActivityEventEntityType = "decision"
	ActivityEventEntityTypeNote     ActivityEventEntityType = "note"
)
