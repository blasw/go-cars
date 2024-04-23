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
        "/add": {
            "post": {
                "description": "Adds cars from source server by regNums",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Add cars",
                "parameters": [
                    {
                        "description": "Cars' regNums",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.addCarsDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/delete": {
            "delete": {
                "description": "Deletes car selected by regNum",
                "tags": [
                    "cars"
                ],
                "summary": "Delete car",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Car's regNum",
                        "name": "regNum",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/get": {
            "post": {
                "description": "Retrieve cars based on various filters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Get cars",
                "parameters": [
                    {
                        "description": "Car filter parameters",
                        "name": "req",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/controllers.carDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Cars"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/update": {
            "patch": {
                "description": "Updates provided fields of car selected by regNum",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "cars"
                ],
                "summary": "Update car",
                "parameters": [
                    {
                        "description": "Updated information",
                        "name": "req",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.updateCarDto"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.addCarsDto": {
            "type": "object",
            "required": [
                "regNums"
            ],
            "properties": {
                "regNums": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "controllers.carDto": {
            "type": "object",
            "properties": {
                "amount": {
                    "description": "Amount of cars per page",
                    "type": "integer"
                },
                "mark": {
                    "description": "Car mark (optional)",
                    "type": "string"
                },
                "model": {
                    "description": "Car model (optional)",
                    "type": "string"
                },
                "owner": {
                    "description": "Car owner filters (optional)",
                    "allOf": [
                        {
                            "$ref": "#/definitions/controllers.ownerDto"
                        }
                    ]
                },
                "page": {
                    "description": "Page number",
                    "type": "integer"
                },
                "regNum": {
                    "description": "Car registration number (optional)",
                    "type": "string"
                },
                "year": {
                    "description": "Car year (optional)",
                    "type": "integer"
                }
            }
        },
        "controllers.ownerDto": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "patronymic": {
                    "description": "Car owner patronymic",
                    "type": "string"
                },
                "surname": {
                    "description": "Car owner surname",
                    "type": "string"
                }
            }
        },
        "controllers.updateCarDto": {
            "type": "object",
            "required": [
                "regNum"
            ],
            "properties": {
                "mark": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "regNum": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "models.Cars": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "mark": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "owner": {
                    "$ref": "#/definitions/models.Owners"
                },
                "ownerID": {
                    "type": "integer"
                },
                "regNum": {
                    "type": "string"
                },
                "year": {
                    "type": "integer"
                }
            }
        },
        "models.Owners": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "patronymic": {
                    "type": "string"
                },
                "surname": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}