package order

import (
	"context"
	"errors"
	"log/slog"

	"github.com/AsakoKabe/gophermart/internal/app/db/models"
	"github.com/AsakoKabe/gophermart/internal/app/db/storage"
)

type Service struct {
	orderStorage storage.OrderStorage
	userStorage  storage.UserStorage
}

func NewService(
	orderStorage storage.OrderStorage,
	userStorage storage.UserStorage,
) *Service {
	return &Service{
		orderStorage: orderStorage,
		userStorage:  userStorage,
	}
}

var ErrAlreadyAddedOtherUser = errors.New("order already added other user")
var ErrAlreadyAdded = errors.New("order already added")
var ErrBadFormat = errors.New("bad format num order")

func (s *Service) Add(ctx context.Context, numOrder string, userLogin string) error {
	user, err := s.userStorage.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.Error("error to get user for add order",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		return err
	}

	existedOrder, err := s.orderStorage.GetOrderByNum(ctx, numOrder)
	if existedOrder != nil {
		if existedOrder.UserID != user.ID {
			return ErrAlreadyAddedOtherUser
		}
		return ErrAlreadyAdded
	}
	if err != nil {
		slog.Error("error to select order", slog.String("err", err.Error()))
		return err
	}

	ok := luhnAlgorithm(numOrder)
	if !ok {
		return ErrBadFormat
	}

	err = s.orderStorage.Add(ctx, &models.Order{
		Num:    numOrder,
		UserID: user.ID,
	})
	if err != nil {
		slog.Error("error to add order", slog.String("err", err.Error()))
		return err
	}

	return nil
}

func (s *Service) GetOrders(ctx context.Context, userLogin string) ([]*models.Order, error) {
	user, err := s.userStorage.GetUserByLogin(ctx, userLogin)
	if err != nil {
		slog.Error("error to get user for get orders",
			slog.String("userLogin", userLogin),
			slog.String("err", err.Error()),
		)
		return nil, err
	}

	orders, err := s.orderStorage.GetOrdersByUserIDSortedByUpdatedAt(ctx, user.ID)
	if err != nil {
		slog.Error("error to get order from storage")
		return nil, err
	}
	return orders, nil
}

func luhnAlgorithm(cardNumber string) bool {
	// this function implements the luhn algorithm
	// it takes as argument a cardnumber of type string
	// and it returns a boolean (true or false) if the
	// card number is valid or not

	// initialise a variable to keep track of the total sum of digits
	total := 0
	// Initialize a flag to track whether the current digit is the second digit from the right.
	isSecondDigit := false

	// iterate through the card number digits in reverse order
	for i := len(cardNumber) - 1; i >= 0; i-- {
		// conver the digit character to an integer
		digit := int(cardNumber[i] - '0')

		if isSecondDigit {
			// double the digit for each second digit from the right
			digit *= 2
			if digit > 9 {
				// If doubling the digit results in a two-digit number,
				// subtract 9 to get the sum of digits.
				digit -= 9
			}
		}

		// Add the current digit to the total sum
		total += digit

		// Toggle the flag for the next iteration.
		isSecondDigit = !isSecondDigit
	}

	// return whether the total sum is divisible by 10
	// making it a valid luhn number
	return total%10 == 0
}
