package msgutil

import (
	"fmt"
)

type RestResp struct {
	Message *string     `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewRestResp(message string, data interface{}) RestResp {
	if message == "" {
		message = "Saved successfully!"
	}
	return RestResp{
		Message: &message,
		Data:    data,
	}
}

func EntityUploadSuccessMsg(entityName string) RestResp {
	return NewRestResp(fmt.Sprintf("%s uploaded successfully", entityName), nil)
}

func EntityGetEmptySuccessMsg(entityName string) RestResp {
	return NewRestResp(fmt.Sprintf("%s has no data", entityName), nil)
}

func EntityUpdateSuccessMsg(entityName string) RestResp {
	return NewRestResp(fmt.Sprintf("%s updated successfully", entityName), nil)
}

func EntityDeleteSuccessMsg(entityName string) RestResp {
	return NewRestResp(fmt.Sprintf("%s deleted successfully", entityName), nil)
}

func EntityNotFoundMsg(entityName string) RestResp {
	return NewRestResp(fmt.Sprintf("%s not found", entityName), nil)
}

func EntityGenericInvalidMsg(entityName string) RestResp {
	return NewRestResp(fmt.Sprintf("invalid %s", entityName), nil)
}
