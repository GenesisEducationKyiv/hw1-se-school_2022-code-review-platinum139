package domain

import (
	"customers-service/internal/common"
)

type CustomersServiceImpl struct {
	logger common.Logger
	repo   CustomersRepo
}

func (s *CustomersServiceImpl) CreateCustomer(customer Customer) error {
	createdCustomer, err := s.repo.CreateCustomer(customer)
	if err != nil {
		s.logger.Errorf("Failed to create customer: %s", err)
		return err
	}

	s.logger.Debugf("Created customer with id %d and email %s",
		createdCustomer.ID, createdCustomer.Email)

	return nil
}

func (s *CustomersServiceImpl) DeleteCustomer(customer Customer) error {
	err := s.repo.DeleteCustomer(customer.ID)
	if err != nil {
		s.logger.Errorf("Failed to delete customer: %s", err)
		return err
	}

	s.logger.Debugf("Deleted customer with id %d",
		customer.ID)

	return nil
}

func (s *CustomersServiceImpl) CreateCustomerWithTransaction(
	customer Customer, transaction ProcessedTransaction) error {

	err := s.repo.CreateCustomerWithTransaction(customer, transaction)
	if err != nil {
		s.logger.Errorf("Failed to create customer with transaction: %s", err)
		return err
	}

	s.logger.Debugf("Created customer with email %s", customer.Email)

	return nil
}

func (s *CustomersServiceImpl) DeleteCustomerWithTransaction(
	customer Customer, transaction ProcessedTransaction) error {

	err := s.repo.DeleteCustomerWithTransaction(customer, transaction)
	if err != nil {
		s.logger.Errorf("Failed to delete customer: %s", err)
		return err
	}

	s.logger.Debugf("Deleted customer with transaction %s",
		transaction.TransactionID)

	return nil
}

func NewCustomersService(logger common.Logger, repo CustomersRepo) *CustomersServiceImpl {
	return &CustomersServiceImpl{
		logger: logger,
		repo:   repo,
	}
}
