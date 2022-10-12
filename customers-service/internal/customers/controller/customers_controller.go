package controller

import (
	"customers-service/config"
	"customers-service/internal/common"
	"customers-service/internal/customers/domain"
	"errors"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CustomersController struct {
	config           *config.AppConfig
	logger           common.Logger
	customersService domain.CustomersService
}

func (c *CustomersController) CreateCustomerHandler(context *gin.Context) {
	c.logger.Debugf("Create customer request")

	var customerDTO CustomerDTO
	if err := context.Bind(&customerDTO); err != nil {
		context.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if customerDTO.Email == "" {
		context.JSON(http.StatusBadRequest, map[string]string{
			"error": "email must be provided",
		})
		return
	}

	err := c.customersService.CreateCustomer(CustomerFromDTO(customerDTO))
	if err != nil {
		var emailAlreadyExistsErr domain.EmailAlreadyExistsError
		if errors.As(err, &emailAlreadyExistsErr) {
			context.JSON(http.StatusBadRequest, map[string]string{
				"error": emailAlreadyExistsErr.Error(),
			})
			return
		}
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusCreated)
}

func (c *CustomersController) DeleteCustomerHandler(context *gin.Context) {
	c.logger.Debugf("Delete customer request")

	idStr := context.Param("id")

	const base = 10
	customerID, err := strconv.ParseInt(idStr, base, c.config.RateValueBitSize)
	if err != nil {
		context.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if customerID < 0 {
		context.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
		return
	}

	customer := domain.Customer{
		ID: customerID,
	}
	if err := c.customersService.DeleteCustomer(customer); err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusOK)
}

func (c *CustomersController) RegisterCustomer(context *gin.Context) interface{} {
	processedTransaction := domain.ProcessedTransaction{
		TransactionID: context.Query("gid"),
	}

	var customerDTO CustomerDTO
	if err := context.Bind(&customerDTO); err != nil {
		return dtmcli.ErrorMessage2Error(err.Error(), dtmcli.ErrFailure)
	}

	if customerDTO.Email == "" {
		return dtmcli.ErrorMessage2Error("email is empty", dtmcli.ErrFailure)
	}

	err := c.customersService.CreateCustomerWithTransaction(
		CustomerFromDTO(customerDTO), processedTransaction)
	if err != nil {
		return dtmcli.ErrorMessage2Error(err.Error(), dtmcli.ErrFailure)
	}

	return nil
}

func (c *CustomersController) RegisterCustomerCompensate(context *gin.Context) interface{} {
	processedTransaction := domain.ProcessedTransaction{
		TransactionID: context.Query("gid"),
	}

	var customerDTO CustomerDTO
	if err := context.Bind(&customerDTO); err != nil {
		return dtmcli.ErrorMessage2Error(err.Error(), dtmcli.ErrFailure)
	}

	err := c.customersService.DeleteCustomerWithTransaction(
		CustomerFromDTO(customerDTO), processedTransaction)
	if err != nil {
		return dtmcli.ErrorMessage2Error(err.Error(), dtmcli.ErrFailure)
	}

	return nil
}

func NewCustomersController(
	config *config.AppConfig,
	logger common.Logger,
	customersService domain.CustomersService,
) *CustomersController {
	return &CustomersController{
		config:           config,
		logger:           logger,
		customersService: customersService,
	}
}
