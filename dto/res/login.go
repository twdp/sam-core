package res

import "tianwei.pro/sam-agent"

type LoginDto struct {
	sam_agent.UserInfo
	Token string `json:"token"`
}
