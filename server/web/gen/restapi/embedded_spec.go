// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Yakit Online Service",
    "title": "Online",
    "version": "0.1.0"
  },
  "basePath": "/api",
  "paths": {
    "/auth/from-github": {
      "get": {
        "security": [],
        "responses": {
          "200": {
            "description": "Fetch Github Oauth URL",
            "schema": {
              "type": "string"
            }
          }
        }
      }
    },
    "/auth/from-github/callback": {
      "get": {
        "security": [],
        "parameters": [
          {
            "type": "string",
            "name": "code",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      }
    },
    "/operation": {
      "get": {
        "parameters": [
          {
            "$ref": "#/parameters/Page"
          },
          {
            "$ref": "#/parameters/Order"
          },
          {
            "$ref": "#/parameters/Limit"
          },
          {
            "$ref": "#/parameters/OrderBy"
          },
          {
            "type": "string",
            "name": "name",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 Operation 记录",
            "schema": {
              "$ref": "#/definitions/OperationsResponse"
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "name": "Data",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/NewOperation"
            }
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      },
      "delete": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      }
    },
    "/user": {
      "get": {
        "security": [
          {
            "trusted": []
          }
        ],
        "parameters": [
          {
            "$ref": "#/parameters/Page"
          },
          {
            "$ref": "#/parameters/Order"
          },
          {
            "$ref": "#/parameters/Limit"
          },
          {
            "$ref": "#/parameters/OrderBy"
          },
          {
            "type": "string",
            "name": "name",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 User 记录",
            "schema": {
              "$ref": "#/definitions/UsersResponse"
            }
          }
        }
      },
      "delete": {
        "security": [
          {
            "trusted": []
          }
        ],
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      }
    },
    "/user/fetch": {
      "get": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 User 记录",
            "schema": {
              "$ref": "#/definitions/User"
            }
          }
        }
      }
    },
    "/user/tags": {
      "get": {
        "responses": {
          "200": {
            "description": "查询 /user  所有的Tags",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "add",
              "set"
            ],
            "type": "string",
            "name": "op",
            "in": "query"
          },
          {
            "type": "array",
            "items": {
              "type": "string"
            },
            "name": "tags",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      }
    },
    "/yakit/plugin": {
      "get": {
        "parameters": [
          {
            "$ref": "#/parameters/Page"
          },
          {
            "$ref": "#/parameters/Order"
          },
          {
            "$ref": "#/parameters/Limit"
          },
          {
            "$ref": "#/parameters/OrderBy"
          },
          {
            "type": "string",
            "name": "name",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 YakitPlugin 记录",
            "schema": {
              "$ref": "#/definitions/YakitPluginsResponse"
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "name": "Data",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/NewYakitPlugin"
            }
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      },
      "delete": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      }
    },
    "/yakit/plugin/fetch": {
      "get": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 YakitPlugin 记录",
            "schema": {
              "$ref": "#/definitions/YakitPlugin"
            }
          }
        }
      }
    },
    "/yakit/plugin/tags": {
      "get": {
        "responses": {
          "200": {
            "description": "查询 /yakit/plugin  所有的Tags",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "add",
              "set"
            ],
            "type": "string",
            "name": "op",
            "in": "query"
          },
          {
            "type": "array",
            "items": {
              "type": "string"
            },
            "name": "tags",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/ActionSucceeded"
          }
        }
      }
    }
  },
  "definitions": {
    "ActionFailed": {
      "type": "object",
      "required": [
        "from",
        "ok",
        "reason"
      ],
      "properties": {
        "from": {
          "description": "来源于哪个 API",
          "type": "string"
        },
        "ok": {
          "description": "执行状态",
          "type": "boolean"
        },
        "reason": {
          "type": "string"
        }
      }
    },
    "ActionSucceeded": {
      "type": "object",
      "required": [
        "from",
        "ok"
      ],
      "properties": {
        "from": {
          "description": "来源于哪个 API",
          "type": "string"
        },
        "ok": {
          "description": "执行状态",
          "type": "boolean"
        }
      }
    },
    "GormBaseModel": {
      "type": "object",
      "required": [
        "id",
        "created_at",
        "updated_at"
      ],
      "properties": {
        "created_at": {
          "type": "integer"
        },
        "id": {
          "type": "integer"
        },
        "updated_at": {
          "type": "integer"
        }
      }
    },
    "NewOperation": {
      "type": "object",
      "required": [
        "type",
        "trigger_user_unique_id",
        "operation_plugin_id"
      ],
      "properties": {
        "extra": {
          "type": "string"
        },
        "operation_plugin_id": {
          "type": "string"
        },
        "trigger_user_unique_id": {
          "type": "string"
        },
        "type": {
          "type": "string"
        }
      }
    },
    "NewUser": {
      "type": "object",
      "required": [
        "uesr_unique_id",
        "user_verbose",
        "from_platform",
        "trusted",
        "tags"
      ],
      "properties": {
        "email": {
          "type": "string"
        },
        "from_platform": {
          "type": "string"
        },
        "tags": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "trusted": {
          "type": "boolean"
        },
        "uesr_unique_id": {
          "type": "string"
        },
        "user_verbose": {
          "type": "string"
        }
      }
    },
    "NewYakitPlugin": {
      "type": "object",
      "required": [
        "type",
        "script_name",
        "authors",
        "content",
        "published_at",
        "tags",
        "is_official",
        "default_open"
      ],
      "properties": {
        "authors": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "content": {
          "type": "string"
        },
        "default_open": {
          "type": "boolean"
        },
        "downloaded_total": {
          "description": "下载次数",
          "type": "integer"
        },
        "forks": {
          "description": "被修改的次数",
          "type": "integer"
        },
        "is_official": {
          "type": "boolean"
        },
        "params": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/YakitPluginParam"
          }
        },
        "published_at": {
          "description": "插件发布的时间",
          "type": "integer"
        },
        "script_name": {
          "type": "string"
        },
        "stars": {
          "description": "获得推荐的次数",
          "type": "integer"
        },
        "tags": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "type": {
          "type": "string"
        }
      }
    },
    "Operation": {
      "allOf": [
        {
          "$ref": "#/definitions/GormBaseModel"
        },
        {
          "$ref": "#/definitions/NewOperation"
        }
      ]
    },
    "OperationsResponse": {
      "required": [
        "data"
      ],
      "allOf": [
        {
          "$ref": "#/definitions/Paging"
        }
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Operation"
          }
        }
      }
    },
    "PageMeta": {
      "description": "描述分页信息的元信息",
      "type": "object",
      "required": [
        "page",
        "limit",
        "total",
        "total_page"
      ],
      "properties": {
        "limit": {
          "description": "页面数据条数限制",
          "type": "integer",
          "default": 20
        },
        "page": {
          "description": "页面索引",
          "type": "integer",
          "default": 1
        },
        "total": {
          "description": "总共数据条数",
          "type": "integer"
        },
        "total_page": {
          "description": "总页数",
          "type": "integer",
          "default": 1
        }
      }
    },
    "Paging": {
      "type": "object",
      "required": [
        "pagemeta"
      ],
      "properties": {
        "pagemeta": {
          "$ref": "#/definitions/PageMeta"
        }
      }
    },
    "Principle": {
      "type": "object",
      "required": [
        "user"
      ],
      "properties": {
        "user": {
          "type": "string"
        }
      }
    },
    "User": {
      "allOf": [
        {
          "$ref": "#/definitions/GormBaseModel"
        },
        {
          "$ref": "#/definitions/NewUser"
        }
      ]
    },
    "UsersResponse": {
      "required": [
        "data"
      ],
      "allOf": [
        {
          "$ref": "#/definitions/Paging"
        }
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/User"
          }
        }
      }
    },
    "YakitPlugin": {
      "allOf": [
        {
          "$ref": "#/definitions/GormBaseModel"
        },
        {
          "$ref": "#/definitions/NewYakitPlugin"
        }
      ]
    },
    "YakitPluginParam": {
      "type": "object",
      "required": [
        "field"
      ],
      "properties": {
        "default_value": {
          "type": "string"
        },
        "field": {
          "type": "string"
        },
        "field_verbose": {
          "type": "string"
        },
        "group": {
          "type": "string"
        },
        "help": {
          "type": "string"
        },
        "required": {
          "type": "boolean"
        },
        "type_verbose": {
          "type": "string"
        }
      }
    },
    "YakitPluginsResponse": {
      "required": [
        "data"
      ],
      "allOf": [
        {
          "$ref": "#/definitions/Paging"
        }
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/YakitPlugin"
          }
        }
      }
    }
  },
  "parameters": {
    "Limit": {
      "type": "integer",
      "default": 20,
      "name": "limit",
      "in": "query"
    },
    "Order": {
      "enum": [
        "asc",
        "desc"
      ],
      "type": "string",
      "default": "desc",
      "name": "order",
      "in": "query"
    },
    "OrderBy": {
      "type": "string",
      "name": "order_by",
      "in": "query"
    },
    "Page": {
      "type": "integer",
      "default": 1,
      "name": "page",
      "in": "query"
    }
  },
  "responses": {
    "ActionFailed": {
      "description": "API 调用失败",
      "schema": {
        "$ref": "#/definitions/ActionFailed"
      }
    },
    "ActionSucceeded": {
      "description": "API 调用成功",
      "schema": {
        "$ref": "#/definitions/ActionSucceeded"
      }
    }
  },
  "securityDefinitions": {
    "trusted": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    },
    "user": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  },
  "security": [
    {
      "user": []
    }
  ]
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Yakit Online Service",
    "title": "Online",
    "version": "0.1.0"
  },
  "basePath": "/api",
  "paths": {
    "/auth/from-github": {
      "get": {
        "security": [],
        "responses": {
          "200": {
            "description": "Fetch Github Oauth URL",
            "schema": {
              "type": "string"
            }
          }
        }
      }
    },
    "/auth/from-github/callback": {
      "get": {
        "security": [],
        "parameters": [
          {
            "type": "string",
            "name": "code",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      }
    },
    "/operation": {
      "get": {
        "parameters": [
          {
            "type": "integer",
            "default": 1,
            "name": "page",
            "in": "query"
          },
          {
            "enum": [
              "asc",
              "desc"
            ],
            "type": "string",
            "default": "desc",
            "name": "order",
            "in": "query"
          },
          {
            "type": "integer",
            "default": 20,
            "name": "limit",
            "in": "query"
          },
          {
            "type": "string",
            "name": "order_by",
            "in": "query"
          },
          {
            "type": "string",
            "name": "name",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 Operation 记录",
            "schema": {
              "$ref": "#/definitions/OperationsResponse"
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "name": "Data",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/NewOperation"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      },
      "delete": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      }
    },
    "/user": {
      "get": {
        "security": [
          {
            "trusted": []
          }
        ],
        "parameters": [
          {
            "type": "integer",
            "default": 1,
            "name": "page",
            "in": "query"
          },
          {
            "enum": [
              "asc",
              "desc"
            ],
            "type": "string",
            "default": "desc",
            "name": "order",
            "in": "query"
          },
          {
            "type": "integer",
            "default": 20,
            "name": "limit",
            "in": "query"
          },
          {
            "type": "string",
            "name": "order_by",
            "in": "query"
          },
          {
            "type": "string",
            "name": "name",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 User 记录",
            "schema": {
              "$ref": "#/definitions/UsersResponse"
            }
          }
        }
      },
      "delete": {
        "security": [
          {
            "trusted": []
          }
        ],
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      }
    },
    "/user/fetch": {
      "get": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 User 记录",
            "schema": {
              "$ref": "#/definitions/User"
            }
          }
        }
      }
    },
    "/user/tags": {
      "get": {
        "responses": {
          "200": {
            "description": "查询 /user  所有的Tags",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "add",
              "set"
            ],
            "type": "string",
            "name": "op",
            "in": "query"
          },
          {
            "type": "array",
            "items": {
              "type": "string"
            },
            "name": "tags",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      }
    },
    "/yakit/plugin": {
      "get": {
        "parameters": [
          {
            "type": "integer",
            "default": 1,
            "name": "page",
            "in": "query"
          },
          {
            "enum": [
              "asc",
              "desc"
            ],
            "type": "string",
            "default": "desc",
            "name": "order",
            "in": "query"
          },
          {
            "type": "integer",
            "default": 20,
            "name": "limit",
            "in": "query"
          },
          {
            "type": "string",
            "name": "order_by",
            "in": "query"
          },
          {
            "type": "string",
            "name": "name",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 YakitPlugin 记录",
            "schema": {
              "$ref": "#/definitions/YakitPluginsResponse"
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "name": "Data",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/NewYakitPlugin"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      },
      "delete": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      }
    },
    "/yakit/plugin/fetch": {
      "get": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "查询 YakitPlugin 记录",
            "schema": {
              "$ref": "#/definitions/YakitPlugin"
            }
          }
        }
      }
    },
    "/yakit/plugin/tags": {
      "get": {
        "responses": {
          "200": {
            "description": "查询 /yakit/plugin  所有的Tags",
            "schema": {
              "type": "array",
              "items": {
                "type": "string"
              }
            }
          }
        }
      },
      "post": {
        "parameters": [
          {
            "type": "integer",
            "name": "id",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "add",
              "set"
            ],
            "type": "string",
            "name": "op",
            "in": "query"
          },
          {
            "type": "array",
            "items": {
              "type": "string"
            },
            "name": "tags",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "API 调用成功",
            "schema": {
              "$ref": "#/definitions/ActionSucceeded"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "ActionFailed": {
      "type": "object",
      "required": [
        "from",
        "ok",
        "reason"
      ],
      "properties": {
        "from": {
          "description": "来源于哪个 API",
          "type": "string"
        },
        "ok": {
          "description": "执行状态",
          "type": "boolean"
        },
        "reason": {
          "type": "string"
        }
      }
    },
    "ActionSucceeded": {
      "type": "object",
      "required": [
        "from",
        "ok"
      ],
      "properties": {
        "from": {
          "description": "来源于哪个 API",
          "type": "string"
        },
        "ok": {
          "description": "执行状态",
          "type": "boolean"
        }
      }
    },
    "GormBaseModel": {
      "type": "object",
      "required": [
        "id",
        "created_at",
        "updated_at"
      ],
      "properties": {
        "created_at": {
          "type": "integer"
        },
        "id": {
          "type": "integer"
        },
        "updated_at": {
          "type": "integer"
        }
      }
    },
    "NewOperation": {
      "type": "object",
      "required": [
        "type",
        "trigger_user_unique_id",
        "operation_plugin_id"
      ],
      "properties": {
        "extra": {
          "type": "string"
        },
        "operation_plugin_id": {
          "type": "string"
        },
        "trigger_user_unique_id": {
          "type": "string"
        },
        "type": {
          "type": "string"
        }
      }
    },
    "NewUser": {
      "type": "object",
      "required": [
        "uesr_unique_id",
        "user_verbose",
        "from_platform",
        "trusted",
        "tags"
      ],
      "properties": {
        "email": {
          "type": "string"
        },
        "from_platform": {
          "type": "string"
        },
        "tags": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "trusted": {
          "type": "boolean"
        },
        "uesr_unique_id": {
          "type": "string"
        },
        "user_verbose": {
          "type": "string"
        }
      }
    },
    "NewYakitPlugin": {
      "type": "object",
      "required": [
        "type",
        "script_name",
        "authors",
        "content",
        "published_at",
        "tags",
        "is_official",
        "default_open"
      ],
      "properties": {
        "authors": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "content": {
          "type": "string"
        },
        "default_open": {
          "type": "boolean"
        },
        "downloaded_total": {
          "description": "下载次数",
          "type": "integer"
        },
        "forks": {
          "description": "被修改的次数",
          "type": "integer"
        },
        "is_official": {
          "type": "boolean"
        },
        "params": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/YakitPluginParam"
          }
        },
        "published_at": {
          "description": "插件发布的时间",
          "type": "integer"
        },
        "script_name": {
          "type": "string"
        },
        "stars": {
          "description": "获得推荐的次数",
          "type": "integer"
        },
        "tags": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "type": {
          "type": "string"
        }
      }
    },
    "Operation": {
      "allOf": [
        {
          "$ref": "#/definitions/GormBaseModel"
        },
        {
          "$ref": "#/definitions/NewOperation"
        }
      ]
    },
    "OperationsResponse": {
      "required": [
        "data"
      ],
      "allOf": [
        {
          "$ref": "#/definitions/Paging"
        }
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Operation"
          }
        }
      }
    },
    "PageMeta": {
      "description": "描述分页信息的元信息",
      "type": "object",
      "required": [
        "page",
        "limit",
        "total",
        "total_page"
      ],
      "properties": {
        "limit": {
          "description": "页面数据条数限制",
          "type": "integer",
          "default": 20
        },
        "page": {
          "description": "页面索引",
          "type": "integer",
          "default": 1
        },
        "total": {
          "description": "总共数据条数",
          "type": "integer"
        },
        "total_page": {
          "description": "总页数",
          "type": "integer",
          "default": 1
        }
      }
    },
    "Paging": {
      "type": "object",
      "required": [
        "pagemeta"
      ],
      "properties": {
        "pagemeta": {
          "$ref": "#/definitions/PageMeta"
        }
      }
    },
    "Principle": {
      "type": "object",
      "required": [
        "user"
      ],
      "properties": {
        "user": {
          "type": "string"
        }
      }
    },
    "User": {
      "allOf": [
        {
          "$ref": "#/definitions/GormBaseModel"
        },
        {
          "$ref": "#/definitions/NewUser"
        }
      ]
    },
    "UsersResponse": {
      "required": [
        "data"
      ],
      "allOf": [
        {
          "$ref": "#/definitions/Paging"
        }
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/User"
          }
        }
      }
    },
    "YakitPlugin": {
      "allOf": [
        {
          "$ref": "#/definitions/GormBaseModel"
        },
        {
          "$ref": "#/definitions/NewYakitPlugin"
        }
      ]
    },
    "YakitPluginParam": {
      "type": "object",
      "required": [
        "field"
      ],
      "properties": {
        "default_value": {
          "type": "string"
        },
        "field": {
          "type": "string"
        },
        "field_verbose": {
          "type": "string"
        },
        "group": {
          "type": "string"
        },
        "help": {
          "type": "string"
        },
        "required": {
          "type": "boolean"
        },
        "type_verbose": {
          "type": "string"
        }
      }
    },
    "YakitPluginsResponse": {
      "required": [
        "data"
      ],
      "allOf": [
        {
          "$ref": "#/definitions/Paging"
        }
      ],
      "properties": {
        "data": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/YakitPlugin"
          }
        }
      }
    }
  },
  "parameters": {
    "Limit": {
      "type": "integer",
      "default": 20,
      "name": "limit",
      "in": "query"
    },
    "Order": {
      "enum": [
        "asc",
        "desc"
      ],
      "type": "string",
      "default": "desc",
      "name": "order",
      "in": "query"
    },
    "OrderBy": {
      "type": "string",
      "name": "order_by",
      "in": "query"
    },
    "Page": {
      "type": "integer",
      "default": 1,
      "name": "page",
      "in": "query"
    }
  },
  "responses": {
    "ActionFailed": {
      "description": "API 调用失败",
      "schema": {
        "$ref": "#/definitions/ActionFailed"
      }
    },
    "ActionSucceeded": {
      "description": "API 调用成功",
      "schema": {
        "$ref": "#/definitions/ActionSucceeded"
      }
    }
  },
  "securityDefinitions": {
    "trusted": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    },
    "user": {
      "type": "apiKey",
      "name": "Authorization",
      "in": "header"
    }
  },
  "security": [
    {
      "user": []
    }
  ]
}`))
}
