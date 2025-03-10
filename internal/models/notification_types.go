package model

// Тип для enum
type NotificationType int

// Константы для типов уведомлений
const (
	Response     NotificationType = iota + 1 // Отклик
	Confirmation                             // Подтверждение
	Rejection                                // Отклонение
)

// Метод для преобразования enum в строку
func (nt NotificationType) String() string {
	switch nt {
	case Response:
		return "отклик"
	case Confirmation:
		return "подтверждение"
	case Rejection:
		return "отклонение"
	default:
		return "неизвестный тип"
	}
}

// Метод для преобразования enum в числовое значение (для базы данных)
func (nt NotificationType) ToInt() int {
	return int(nt)
}

func ConvertToNotidy(id int) NotificationType {
	switch id {
	case 1:
		return Response
	case 2:
		return Confirmation
	case 3:
		return Rejection
	default:
		return Response
	}
}
