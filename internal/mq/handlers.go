package mq

import (
	"encoding/json"
	"fmt"

	"goflow/internal/pkg/logger"

	"github.com/ThreeDotsLabs/watermill/message"
)

// HandleEmail 处理邮件发送任务
func HandleEmail(msg *message.Message) error {
	var payload EmailPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("unmarshal email payload: %w", err)
	}
	// TODO: 替换为真实的邮件发送逻辑（如 SMTP / SendGrid）
	logger.Log.Infow("[mq] send email",
		"msg_uuid", msg.UUID,
		"to", payload.To,
		"subject", payload.Subject,
	)
	return nil
}

// HandleSMS 处理短信发送任务
func HandleSMS(msg *message.Message) error {
	var payload SMSPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("unmarshal sms payload: %w", err)
	}
	// TODO: 替换为真实的短信发送逻辑（如阿里云 / 腾讯云短信）
	logger.Log.Infow("[mq] send sms",
		"msg_uuid", msg.UUID,
		"phone", payload.Phone,
		"content", payload.Content,
	)
	return nil
}

// HandleStats 处理数据统计任务
func HandleStats(msg *message.Message) error {
	var payload StatsPayload
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return fmt.Errorf("unmarshal stats payload: %w", err)
	}
	// TODO: 替换为真实的统计逻辑（如写入 ClickHouse / 更新聚合表）
	logger.Log.Infow("[mq] record stats",
		"msg_uuid", msg.UUID,
		"event", payload.Event,
		"user_id", payload.UserID,
		"properties", payload.Properties,
	)
	return nil
}
