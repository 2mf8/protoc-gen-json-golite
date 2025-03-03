package onebot

type Action int32

const (
	Action_unknown          Action = 0
	Action_send_group_msg   Action = 1
	Action_send_private_msg Action = 2
)

// Enum value maps for Action.
var (
	Action_name = map[int32]string{
		0: "unknown",
		1: "send_group_msg",
		2: "send_private_msg",
	}
	Action_value = map[string]int32{
		"unknown":          0,
		"send_group_msg":   1,
		"send_private_msg": 2,
	}
)

type Message struct {
	Type string            `json:"type,omitempty"`
	Data map[string]string `json:"data,omitempty"`
}

type GroupMessageEvent struct {
	Time        int64                        `json:"time,omitempty"`
	SelfId      int64                        `json:"self_id,omitempty"`
	PostType    string                       `json:"post_type,omitempty"`
	MessageType string                       `json:"message_type,omitempty"`
	SubType     string                       `json:"sub_type,omitempty"`
	MessageId   int32                        `json:"message_id,omitempty"`
	GroupId     int64                        `json:"group_id,omitempty"`
	UserId      int64                        `json:"user_id,omitempty"`
	Anonymous   *GroupMessageEvent_Anonymous `json:"anonymous,omitempty"`
	Message     []*Message                   `json:"message,omitempty"`
	RawMessage  string                       `json:"raw_message,omitempty"`
	Font        int32                        `json:"font,omitempty"`
	Sender      *GroupMessageEvent_Sender    `json:"sender,omitempty"`
	Extra       map[string]string            `json:"extra,omitempty"`
}

type GroupMessageEvent_Anonymous struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Flag string `json:"flag,omitempty"`
}

type GroupMessageEvent_Sender struct {
	UserId   int64  `json:"user_id,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	Card     string `json:"card,omitempty"`
	Sex      string `json:"sex,omitempty"`
	Age      int32  `json:"age,omitempty"`
	Area     string `json:"area,omitempty"`
	Level    string `json:"level,omitempty"`
	Role     string `json:"role,omitempty"`
	Title    string `json:"title,omitempty"`
}
