package constants

const (
	STATUS_SUCCESS                = "00"
	STATUS_STILL_PROCESS          = "01"
	STATUS_FAILED                 = "05"
	STATUS_INVALID_REQUEST_FORMAT = "400"
	STATUS_UNAUTHORIZED           = "401"
	STATUS_FORBIDDEN              = "403"
	STATUS_NOT_FOUND              = "404"
	STATUS_CONFLICT               = "409"

	OBERON_STATUS_GENERAL_SUCCESS = "0000"
	OBERON_STATUS_GENERAL_ERROR   = "5000"
	OBERON_STATUS_DATA_NOT_FOUND  = "5007"
)

const (
	MESSAGE_SUCCESS                = "Success"
	MESSAGE_STILL_PROCESS          = "Transaction is being process"
	MESSAGE_FAILED                 = "Something went wrong"
	MESSAGE_INVALID_REQUEST_FORMAT = "Invalid Request Format"
	MESSAGE_UNAUTHORIZED           = "Unauthorized"
	MESSAGE_FORBIDDEN              = "Forbidden"
	MESSAGE_NOT_FOUND              = "Not Found"
	MESSAGE_CONFLICT               = "Conflict"
)

const (
	SuccessCode         string = "00"
	ErrorInvalidRequest string = "97"
	ErrorInvalidJson    string = "98"
	GeneralError        string = "99"
)

type DefaultResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
}

type NewDefaultResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data"`
	Errors  []interface{} `json:"errors"`
}

type PaginationData struct {
	Page        uint `json:"page"`
	TotalPages  uint `json:"totalPages"`
	TotalItems  uint `json:"totalItems"`
	Limit       uint `json:"limit"`
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
}

type PaginationResponseData struct {
	Results        interface{} `json:"results"`
	PaginationData `json:"pagination"`
}
