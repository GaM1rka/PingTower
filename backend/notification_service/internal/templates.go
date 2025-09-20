package internal

import (
	"fmt"
	"notification_service/models"
)

func GenerateEmailContent(req models.NotificationRequest) models.EmailContent {
	subject := fmt.Sprintf("[CRITICAL] %s DOWN - Action Required", req.Site)
	
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Critical Alert</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1); }
        .header { background-color: #dc3545; color: white; padding: 20px; text-align: center; }
        .header h1 { margin: 0; font-size: 24px; }
        .content { padding: 30px; }
        .alert-icon { font-size: 48px; text-align: center; margin-bottom: 20px; }
        .status-badge { display: inline-block; background-color: #dc3545; color: white; padding: 8px 16px; border-radius: 20px; font-weight: bold; margin: 10px 0; }
        .details { background-color: #f8f9fa; padding: 20px; border-radius: 6px; margin: 20px 0; border-left: 4px solid #dc3545; }
        .details h3 { margin-top: 0; color: #495057; }
        .detail-item { margin: 10px 0; }
        .label { font-weight: bold; color: #495057; }
        .value { color: #212529; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #6c757d; font-size: 14px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚨 КРИТИЧЕСКАЯ ПРОБЛЕМА ОБНАРУЖЕНА</h1>
        </div>
        
        <div class="content">
            <div class="alert-icon">⚠️</div>
            
            <div class="details">
                <h3>📊 Детали инцидента:</h3>
                <div class="detail-item">
                    <span class="label">Сервис:</span> 
                    <span class="value">%s</span>
                </div>
                <div class="detail-item">
                    <span class="label">Статус:</span> 
                    <span class="status-badge">НЕДОСТУПЕН</span>
                </div>
                <div class="detail-item">
                    <span class="label">Время обнаружения:</span> 
                    <span class="value">%s</span>
                </div>
                <div class="detail-item">
                    <span class="label">Тип проблемы:</span> 
                    <span class="value">Сервис не отвечает на запросы</span>
                </div>
            </div>
            
            <div style="background-color: #fff3cd; border: 1px solid #ffeaa7; border-radius: 6px; padding: 15px; margin: 20px 0;">
                <h4 style="margin-top: 0; color: #856404;">⚡ Требуется немедленное действие:</h4>
                <ul style="margin-bottom: 0; color: #856404;">
                    <li>Проверьте статус сервера</li>
                    <li>Проанализируйте логи системы</li>
                    <li>Уведомите команду DevOps</li>
                    <li>Подготовьте план восстановления</li>
                </ul>
            </div>
        </div>
        
        <div class="footer">
            <p>Это автоматическое уведомление от системы мониторинга PingTower</p>
            <p>Время отправки: %s</p>
        </div>
    </div>
</body>
</html>`, req.Site, req.Time, req.Time)

	text := fmt.Sprintf(`
🚨 КРИТИЧЕСКАЯ ПРОБЛЕМА ОБНАРУЖЕНА

Сервис: %s
Статус: НЕДОСТУПЕН
Время: %s

📊 Детали:
- Тип проблемы: Сервис не отвечает на запросы
- Требуется немедленное действие

⚡ Что нужно сделать:
1. Проверьте статус сервера
2. Проанализируйте логи системы
3. Уведомите команду DevOps
4. Подготовьте план восстановления

---
Это автоматическое уведомление от системы мониторинга PingTower
`, req.Site, req.Time)

	return models.EmailContent{
		Subject: subject,
		HTML:    html,
		Text:    text,
	}
}