package errors

// AppError структура ошибки для API
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Список ошибок на русском
var (
	// Общие ошибки
	ErrUnauthorized = &AppError{
		Code:    "UNAUTHORIZED",
		Message: "Не авторизован. Пожалуйста, войдите в систему",
		Status:  401,
	}

	ErrForbidden = &AppError{
		Code:    "FORBIDDEN",
		Message: "Доступ запрещен. Недостаточно прав",
		Status:  403,
	}

	ErrNotFound = &AppError{
		Code:    "NOT_FOUND",
		Message: "Запись не найдена",
		Status:  404,
	}

	ErrBadRequest = &AppError{
		Code:    "BAD_REQUEST",
		Message: "Неверный запрос",
		Status:  400,
	}

	ErrInternalServer = &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "Внутренняя ошибка сервера",
		Status:  500,
	}

	// Ошибки аутентификации
	ErrInvalidCredentials = &AppError{
		Code:    "INVALID_CREDENTIALS",
		Message: "Неверный email или пароль",
		Status:  401,
	}

	ErrEmailAlreadyExists = &AppError{
		Code:    "EMAIL_ALREADY_EXISTS",
		Message: "Пользователь с таким email уже существует",
		Status:  400,
	}

	ErrInvalidToken = &AppError{
		Code:    "INVALID_TOKEN",
		Message: "Недействительный токен авторизации",
		Status:  401,
	}

	// Ошибки мероприятий
	ErrEventNotFound = &AppError{
		Code:    "EVENT_NOT_FOUND",
		Message: "Мероприятие не найдено",
		Status:  404,
	}

	ErrEventCancelled = &AppError{
		Code:    "EVENT_CANCELLED",
		Message: "Мероприятие отменено",
		Status:  400,
	}

	ErrEventCompleted = &AppError{
		Code:    "EVENT_COMPLETED",
		Message: "Мероприятие уже завершено",
		Status:  400,
	}

	ErrRegistrationDeadlinePassed = &AppError{
		Code:    "REGISTRATION_DEADLINE_PASSED",
		Message: "Дедлайн регистрации истек",
		Status:  400,
	}

	ErrEventAlreadyStarted = &AppError{
		Code:    "EVENT_ALREADY_STARTED",
		Message: "Мероприятие уже началось",
		Status:  400,
	}

	// Ошибки участия
	ErrAlreadyRegistered = &AppError{
		Code:    "ALREADY_REGISTERED",
		Message: "Вы уже зарегистрированы на это мероприятие",
		Status:  400,
	}

	ErrNotRegistered = &AppError{
		Code:    "NOT_REGISTERED",
		Message: "Вы не зарегистрированы на это мероприятие",
		Status:  400,
	}

	ErrInvalidQRCode = &AppError{
		Code:    "INVALID_QR_CODE",
		Message: "Недействительный QR-код",
		Status:  400,
	}

	ErrQRCodeAlreadyUsed = &AppError{
		Code:    "QR_CODE_ALREADY_USED",
		Message: "Этот QR-код уже был использован",
		Status:  400,
	}

	ErrParticipationCancelled = &AppError{
		Code:    "PARTICIPATION_CANCELLED",
		Message: "Регистрация на мероприятие отменена",
		Status:  400,
	}

	ErrCannotCancelConfirmed = &AppError{
		Code:    "CANNOT_CANCEL_CONFIRMED",
		Message: "Нельзя отменить участие после подтверждения",
		Status:  400,
	}

	ErrNotOrganizer = &AppError{
		Code:    "NOT_ORGANIZER",
		Message: "Только организатор мероприятия может выполнить это действие",
		Status:  403,
	}

	// Ошибки организаторов
	ErrOnlyOrganizers = &AppError{
		Code:    "ONLY_ORGANIZERS",
		Message: "Только организаторы могут создавать мероприятия",
		Status:  403,
	}

	ErrNotYourEvent = &AppError{
		Code:    "NOT_YOUR_EVENT",
		Message: "Вы можете редактировать только свои мероприятия",
		Status:  403,
	}

	// Ошибки отзывов
	ErrAlreadyReviewed = &AppError{
		Code:    "ALREADY_REVIEWED",
		Message: "Вы уже оставляли отзыв об этом организаторе",
		Status:  400,
	}

	ErrCannotReviewSelf = &AppError{
		Code:    "CANNOT_REVIEW_SELF",
		Message: "Нельзя оставить отзыв о самом себе",
		Status:  400,
	}
)

// New создает новую ошибку с сообщением
func New(message string) *AppError {
	return &AppError{
		Code:    "CUSTOM_ERROR",
		Message: message,
		Status:  400,
	}
}

// Wrap оборачивает ошибку с кодом и сообщением
func Wrap(err error, code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  400,
	}
}
