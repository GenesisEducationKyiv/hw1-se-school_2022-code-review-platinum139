package controller

import "customers-service/internal/customers/domain"

func CustomerToDTO(customer domain.Customer) CustomerDTO {
	return CustomerDTO{
		ID:    customer.ID,
		Email: customer.Email,
	}
}

func CustomerFromDTO(customerDTO CustomerDTO) domain.Customer {
	return domain.Customer{
		ID:    customerDTO.ID,
		Email: customerDTO.Email,
	}
}
