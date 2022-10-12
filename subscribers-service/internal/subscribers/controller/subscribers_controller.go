package controller

import (
	"errors"
	"github.com/dtm-labs/dtm/client/dtmcli"
	"github.com/gin-gonic/gin"
	"net/http"
	"subscribers-service/config"
	"subscribers-service/internal/common"
	"subscribers-service/internal/subscribers/domain"
	"subscribers-service/pkg/file_storage"
)

type Controller struct {
	logger             common.Logger
	config             *config.AppConfig
	subscribersService domain.SubscribersService
}

func (c Controller) SubscribeHandler(context *gin.Context) {
	subscriberDTO := SubscriberDTO{}
	if err := context.Bind(&subscriberDTO); err != nil {
		c.logger.Errorf("Unable to get form data from request:", err)
		context.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad request",
		})
		return
	}

	if subscriberDTO.Email == "" {
		c.logger.Errorf("Bad request: email was not provided")
		context.JSON(http.StatusBadRequest, map[string]string{
			"error": "email must be provided",
		})
		return
	}

	dtmAddr := c.config.DTMCoordinatorAddr
	globalTransactionID := dtmcli.MustGenGid(dtmAddr)

	subscribersSvcAddr := c.config.SubscribersServiceAddr
	customersSvcAddr := c.config.CustomersServiceAddr

	registerSubscriberURL := subscribersSvcAddr + "/register-subscriber"
	registerSubscriberCompensateURL := subscribersSvcAddr + "/register-subscriber-compensate"

	registerCustomerURL := customersSvcAddr + "/register-customer"
	registerCustomerCompensateURL := customersSvcAddr + "/register-customer-compensate"

	reqBody, err := common.StructToMap(subscriberDTO)
	if err != nil {
		context.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
		return
	}

	saga := dtmcli.
		NewSaga(dtmAddr, globalTransactionID).
		Add(registerSubscriberURL, registerSubscriberCompensateURL, reqBody).
		Add(registerCustomerURL, registerCustomerCompensateURL, reqBody)

	saga.WaitResult = true
	saga.WithRetryLimit(1)

	if err = saga.Submit(); err != nil {
		code, data := dtmcli.Result2HttpJSON(err)
		dataMap, ok := data.(map[string]string)
		if !ok {
			context.JSON(code, data)
			return
		}
		context.Data(code, "application/json", []byte(dataMap["error"]))
		return
	}

	context.Status(http.StatusOK)
}

func (c Controller) RegisterSubscriber(context *gin.Context) interface{} {
	subscriberDTO := SubscriberDTO{}
	if err := context.Bind(&subscriberDTO); err != nil {
		return err
	}

	subscriber := SubscriberFromDTO(subscriberDTO)
	subscriber.TransactionID = context.Query("gid")

	err := c.subscribersService.Subscribe(subscriber)
	if err != nil {
		if errors.Is(err, file_storage.RecordAlreadyExistsError{}) {
			context.JSON(http.StatusConflict, map[string]string{
				"error": "email already exists",
			})
			return dtmcli.ErrFailure
		}
		context.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
		return dtmcli.ErrFailure
	}

	return nil
}

func (c Controller) RegisterSubscriberCompensate(context *gin.Context) interface{} {
	subscriber := domain.Subscriber{
		TransactionID: context.Query("gid"),
	}

	err := c.subscribersService.Unsubscribe(subscriber)
	if err != nil {
		return err
	}

	return nil
}

func NewSubscribersController(
	logger common.Logger,
	config *config.AppConfig,
	subscribersService domain.SubscribersService,
) *Controller {
	return &Controller{
		logger:             logger,
		config:             config,
		subscribersService: subscribersService,
	}
}
