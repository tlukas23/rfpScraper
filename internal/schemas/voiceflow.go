package schemas

const (
	ActionLaunch = "launch"
	ActionText   = "text"
)

type VoiceFlowQuestion struct {
	Action VoiceFlowAction `json:"action"`
}

type VoiceFlowAction struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type VoiceFlowResponse struct {
	Type    string           `json:"type"`
	Payload VoiceFlowPayload `json:"payload"`
}

type VoiceFlowPayload struct {
	Message string `json:"message"`
}

type VoiceFlowExcelRow struct {
	Question string
	Response string
}
