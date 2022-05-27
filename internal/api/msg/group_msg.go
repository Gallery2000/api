package msg

import (
	comm2 "github.com/glide-im/api/internal/api/comm"
	"github.com/glide-im/api/internal/api/router"
	"github.com/glide-im/api/internal/dao/common"
	"github.com/glide-im/api/internal/dao/msgdao"
	"github.com/glide-im/api/internal/dao/userdao"
	"github.com/glide-im/glide/pkg/messages"
	"math"
)

type GroupMsgApi struct {
}

func (*GroupMsgApi) GetRecentGroupMessage(ctx *route.Context, request *RecentGroupMessageRequest) error {
	ms, err := msgdao.GroupMsgDaoImpl.GetLatestGroupMessage(request.Gid, 20)
	if err != nil && err != common.ErrNoRecordFound {
		return comm2.NewDbErr(err)
	}
	//goland:noinspection GoPreferNilSlice
	resp := []*GroupMessageResponse{}
	for _, m := range ms {
		resp = append(resp, dbGroupMsg2ResponseMsg(m))
	}
	ctx.Response(messages.NewMessage(ctx.Seq, comm2.ActionSuccess, resp))
	return nil
}

func (*GroupMsgApi) GetGroupMessageHistory(ctx *route.Context, request *GroupMsgHistoryRequest) error {
	before := request.BeforeSeq
	if request.BeforeSeq <= 0 {
		before = math.MaxInt64
	}
	ms, err := msgdao.GroupMsgDaoImpl.GetGroupMessage(request.Gid, before, 20)
	if err != nil {
		return comm2.NewDbErr(err)
	}
	//goland:noinspection GoPreferNilSlice
	resp := []*GroupMessageResponse{}
	for _, m := range ms {
		resp = append(resp, dbGroupMsg2ResponseMsg(m))
	}
	ctx.Response(messages.NewMessage(ctx.Seq, comm2.ActionSuccess, resp))
	return nil
}

func (*GroupMsgApi) GetGroupMessage(ctx *route.Context, request *GroupMessageRequest) error {

	ms, err := msgdao.GroupMsgDaoImpl.GetMessages(request.Mid...)
	if err != nil {
		return comm2.NewDbErr(err)
	}
	resp := make([]*GroupMessageResponse, len(ms))
	for _, m := range ms {
		resp = append(resp, dbGroupMsg2ResponseMsg(m))
	}
	ctx.Response(messages.NewMessage(ctx.Seq, comm2.ActionSuccess, resp))
	return nil
}

func (*GroupMsgApi) GetUserGroupMessageState(ctx *route.Context) error {
	groups, err := userdao.Dao.GetContactsByType(ctx.Uid, 2)
	if err != nil && err != common.ErrNoRecordFound {
		return comm2.NewDbErr(err)
	}
	if len(groups) == 0 {
		ctx.Response(messages.NewMessage(ctx.Seq, comm2.ActionSuccess, []string{}))
		return nil
	}
	//goland:noinspection GoPreferNilSlice
	gid := []int64{}
	for _, group := range groups {
		gid = append(gid, group.Id)
	}
	state, err := msgdao.GroupMsgDaoImpl.GetGroupsMessageState(gid...)
	if err != nil {
		return comm2.NewDbErr(err)
	}
	ctx.Response(messages.NewMessage(ctx.Seq, comm2.ActionSuccess, state))
	return nil
}

func (*GroupMsgApi) GetGroupMessageState(ctx *route.Context, request *GroupMsgStateRequest) error {

	state, err := msgdao.GroupMsgDaoImpl.GetGroupMessageState(request.Gid)
	if err != nil {
		return comm2.NewDbErr(err)
	}
	ctx.Response(messages.NewMessage(ctx.Seq, comm2.ActionSuccess, GroupMessageStateResponse{state}))
	return nil
}

func dbGroupMsg2ResponseMsg(m *msgdao.GroupMessage) *GroupMessageResponse {
	return &GroupMessageResponse{
		Mid:      m.MID,
		Sender:   m.From,
		Seq:      m.Seq,
		Gid:      m.To,
		Type:     m.Type,
		SendAt:   m.SendAt,
		Content:  m.Content,
		Status:   m.Status,
		RecallBy: m.RecallBy,
	}
}
