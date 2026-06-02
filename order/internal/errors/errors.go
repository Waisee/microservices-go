package errs

import "errors"

var (
	ErrOrderNotFound    = errors.New("заказ не найден")
	ErrOrderAlreadyPaid = errors.New("заказ уже оплачен")
	ErrOrderCancelled   = errors.New("заказ отменён")
	ErrPartNotFound     = errors.New("деталь не найдена")
	ErrOutOfStock       = errors.New("деталь отсутствует на складе")
	ErrInvalidPartUUID  = errors.New("неверный формат UUID детали")
)
