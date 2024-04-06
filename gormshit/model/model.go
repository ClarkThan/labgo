package model

import "time"

type Event struct {
	ID            int64     `json:"id" gorm:"column:id"`
	Action        int64     `json:"action" gorm:"column:action"`
	EntID         int64     `json:"enterprise_id" gorm:"column:enterprise_id"`
	TargetID      string    `json:"target_id" gorm:"column:target_id"`
	TargetKind    string    `json:"target_kind" gorm:"column:target_kind"`
	Body          []byte    `json:"body" gorm:"column:body"`
	TrackID       string    `json:"track_id" gorm:"column:track_id"`
	AgentID       int64     `json:"agent_id" gorm:"column:agent_id"`
	AgentNickname string    `json:"agent_nickname" gorm:"column:agent_nickname"`
	CreatedOn     time.Time `json:"created_on" gorm:"column:created_on"`
}

func (*Event) TableName() string {
	return "event"
}

type EventActive Event

func (*EventActive) TableName() string {
	return "event_active"
}

type Conversation struct {
	ID                    int64      `gorm:"column:id"`
	EntID                 int64      `gorm:"column:enterprise_id"`
	TrackID               string     `gorm:"column:track_id"`
	VisitID               string     `gorm:"column:visit_id"`
	AgentID               int64      `gorm:"column:agent_id"`
	AgentType             string     `gorm:"column:agent_type"`
	Assignee              int64      `gorm:"column:assignee"`
	CreatedOn             *time.Time `gorm:"column:created_on"`
	EndedOn               *time.Time `gorm:"column:ended_on"`
	EndedBy               string     `gorm:"column:ended_by"`
	MsgNum                int        `gorm:"column:msg_num"`
	ClientMsgNum          int        `gorm:"column:client_msg_num"`
	AgentMsgNum           int        `gorm:"column:agent_msg_num"`
	AgentEffectiveMsgNum  int        `gorm:"column:agent_effective_msg_num"`
	FirstResponseWaitTime int        `gorm:"column:first_response_wait_time"`
	ClientFirstSendTime   *time.Time `gorm:"column:client_first_send_time"`
	FirstMsgCreatedOn     *time.Time `gorm:"column:first_msg_created_on"`
	LastMsgCreatedOn      *time.Time `gorm:"column:last_msg_created_on"`
	LastMsgContent        string     `gorm:"column:last_msg_content"`
	ConvDuration          int        `gorm:"column:conversation_duration"`
	QualityGrade          string     `gorm:"column:quality_grade"`
	Clues                 string     `gorm:"column:clues"`
	URL                   string     `gorm:"column:url"`
	Title                 string     `gorm:"column:title"`
	LastUpdated           *time.Time `gorm:"column:last_updated"`
	Source                string     `gorm:"column:source"`
	SubSource             string     `gorm:"column:sub_source"`
}

func (Conversation) TableName() string {
	return "conversation_v1"
}

type ConvActive Conversation

func (*ConvActive) TableName() string {
	return "conversation_active"
}
