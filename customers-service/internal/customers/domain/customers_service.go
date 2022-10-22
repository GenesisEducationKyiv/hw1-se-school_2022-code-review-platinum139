package domain

type CustomersService interface {
	CreateCustomer(customer Customer) error
	DeleteCustomer(customer Customer) error
	CreateCustomerWithTransaction(customer Customer, transaction ProcessedTransaction) error
	DeleteCustomerWithTransaction(customer Customer, transaction ProcessedTransaction) error
}
