package domain

type CustomersRepo interface {
	CreateCustomer(customer Customer) (*Customer, error)
	DeleteCustomer(customerID int64) error
	CreateCustomerWithTransaction(customer Customer, transaction ProcessedTransaction) error
	DeleteCustomerWithTransaction(customer Customer, transaction ProcessedTransaction) error
}
