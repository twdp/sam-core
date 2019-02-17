package res

import (
	"tianwei.pro/business"
	"tianwei.pro/sam-agent"
)

type LoginDto struct {
	business.Response

	sam_agent.UserInfo
	Token string `json:"token"`
}

