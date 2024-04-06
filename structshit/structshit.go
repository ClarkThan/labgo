package structshit

import (
	"log"

	"github.com/fatih/structs"
)

type SchedulerReq struct {
	EntStr                   string         `json:"ent_id"` // TODO(ranya)  use_ent_token
	EntID                    int64          `json:"-"`
	TrackID                  string         `json:"track_id"`
	ClientID                 string         `json:"-"`
	UserName                 string         `json:"user_name"` // 前端应该不会传吧
	AgentToken               string         `json:"agent_token"`
	GroupToken               string         `json:"group_token"`
	Captcha                  string         `json:"captcha"`
	FromType                 string         `json:"from_type"`
	ConvInitiateType         string         `json:"conv_initiate_type"`
	Content                  string         `json:"content"`
	ConvCreateType           int8           `json:"conv_create_type"`
	Fallback                 int8           `json:"fallback"`
	Queueing                 bool           `json:"queueing"`
	ReturningCustomerEnabled bool           `json:"returning_customer_enabled"`
	OnlyMpushRobot           bool           `json:"only_mpush_robot"`
	MpushRedirect            bool           `json:"mpush_redirect"`
	NoQueueing               bool           `json:"-"`
	Trusteeship              bool           `json:"-"`
	UseReserveScheduler      bool           `json:"-"`
	UseReserveTicket         bool           `json:"-"`
	FirstMatchMpushRobot     bool           `json:"-"`
	InConv                   bool           `json:"-"`
	MpushServability         bool           `json:"-"`
	TriggerType              string         `json:"-"`
	BaseRule                 string         `json:"-"`
	ClientFirstMsg           string         `json:"-"`
	ClientMessage            string         `json:"-"`
	ReserveToken             string         `json:"-"`
	Source                   string         `json:"source"`
	SubSource                string         `json:"-"`
	SourceToken              string         `json:"source_token"`
	SubSourceName            string         `json:"sub_source_name"`
	Title                    string         `json:"title"`
	URL                      string         `json:"url"`
	ReferrerURL              string         `json:"referrer_url"`
	UAString                 string         `json:"us_string"`
	VisitID                  string         `json:"visit_id"`
	AllocateType             string         `json:"-"`
	MpushBotID               int64          `json:"-"`
	TargetAgentIDs           []int64        `json:"-"`
	TargetGroupIDs           []int64        `json:"-"`
	StickyAgentID            int64          `json:"-"`
	FallbackGroupID          int64          `json:"-"`
	ExcludeAgentTokens       []string       `json:"exclude_agent_tokens"`
	AllocChain               []string       `json:"-"`
	Exclusions               []int64        `json:"exclusions"`
	MpushTemplate            map[string]any `json:"mpush_template"`
	Extra                    struct {
		CaptchaToken        string
		NoMsgFilterRefactor bool
	} `json:"-"`
}

func (r *SchedulerReq) ToMap() map[string]any {
	st := structs.New(r)
	st.TagName = "json"
	return st.Map()
}

func demo1() {
	r := SchedulerReq{
		ClientMessage: "hello",
		Source:        "weixin",
		VisitID:       "fuckyou",
		EntID:         23,
	}

	log.Printf("map = %#v\n", r.ToMap())
}

func Main() {
	demo1()
}
