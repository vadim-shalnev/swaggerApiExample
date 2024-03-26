// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/address/geocode": {
            "post": {
                "description": "Поиск координат по адресу",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "geocode"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Models.RequestQuery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Models.RequestQuery"
                        }
                    }
                }
            }
        },
        "/api/address/search": {
            "post": {
                "description": "Поиск полного адреса",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "geocode"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Models.RequestQuery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/Models.RequestQuery"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "Вход в систему с указанным email и паролем",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "parameters": [
                    {
                        "description": "Данные пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Models.NewUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный вход в систему",
                        "schema": {
                            "$ref": "#/definitions/Models.NewUserResponse"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "Регистрация нового пользователя с указанным email и паролем",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "parameters": [
                    {
                        "description": "Данные пользователя",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/Models.NewUserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная регистрация",
                        "schema": {
                            "$ref": "#/definitions/Models.NewUserResponse"
                        }
                    }
                }
            }
        },
        "/users/get/{id}": {
            "get": {
                "description": "Получить информацию о пользователе по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Пользователь",
                        "schema": {
                            "$ref": "#/definitions/Models.NewUserResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удалить пользователя по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Succsec",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "Models.NewUserRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "Models.NewUserResponse": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                },
                "token": {
                    "$ref": "#/definitions/Models.TokenString"
                }
            }
        },
        "Models.RequestQuery": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string"
                }
            }
        },
        "Models.TokenString": {
            "type": "object",
            "properties": {
                "authController": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/",
	Schemes:          []string{},
	Title:            "Swagger Example API",
	Description:      "This is a geocode api server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
