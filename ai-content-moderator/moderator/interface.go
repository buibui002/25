package moderator

type ModerationAction string

const (
	ActionPass   ModerationAction = "pass"
	ActionBlock  ModerationAction = "block"
	ActionModify ModerationAction = "modify"
)

type ModerationResult struct {
	Action          ModerationAction
	Reason          string
	ModifiedContent string
}

type ContentModerator interface {
	Check(content string) (*ModerationResult, error)
}
