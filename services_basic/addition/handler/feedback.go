package handler

import (
	"context"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/addition/global"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/addition/model"
	"github.com/Winszheng/grad_project_ecommerce_microservices_gin/services_basic/addition/proto"
)

func (*UserOpServer) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var rsp proto.MessageListResponse
	var messages []model.Feedback
	var messageList []*proto.MessageResponse

	result := global.DB.Where(&model.Feedback{UserID: req.UserId}).Find(&messages)
	rsp.Total = int32(result.RowsAffected)

	for _, message := range messages {
		messageList = append(messageList, &proto.MessageResponse{
			Id:          message.ID,
			UserId:      message.UserID,
			MessageType: message.Type,
			Subject:     message.Title,
			Message:     message.Msg,
			File:        message.File,
		})
	}

	rsp.Data = messageList
	return &rsp, nil
}

func (*UserOpServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	var message model.Feedback

	message.UserID = req.UserId
	message.Type = req.MessageType
	message.Title = req.Subject
	message.Msg = req.Message
	message.File = req.File

	global.DB.Save(&message)

	return &proto.MessageResponse{Id: message.ID}, nil
}
