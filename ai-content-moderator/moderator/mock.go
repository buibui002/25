package moderator

import "strings"

type MockModerator struct{}

func (m *MockModerator) Check(content string) (*ModerationResult, error) {
	if strings.Contains(content, "badword") {
		return &ModerationResult{
			Action: ActionBlock,
			Reason: "Contains badword",
		}, nil
	}

	if strings.Contains(content, "replace") {
		return &ModerationResult{
			Action:          ActionModify,
			Reason:          "Sensitive word replaced",
			ModifiedContent: strings.ReplaceAll(content, "replace", "***"),
		}, nil
	}

	return &ModerationResult{
		Action: ActionPass,
		Reason: "Clean",
	}, nil
}
