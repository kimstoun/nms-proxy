package netConfigServe
 const Swagger ={
  "swagger": "2.0",
  "info": {
    "title": "netServer.proto",
    "version": "version not set"
  },
  "schemes": [
    "http",
    "https"
  ],
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "paths": {
    "/links/echo": {
      "get": {
        "operationId": "QueryAllLinks",
        "responses": {
          "200": {
            "description": "",
            "schema": {
              "$ref": "#/definitions/netConfigServeLinks"
            }
          }
        },
        "parameters": [
          {
            "name": "rioId",
            "in": "query",
            "required": false,
            "type": "integer",
            "format": "int32"
          }
        ],
        "tags": [
          "NetConfig"
        ]
      }
    }
  },
  "definitions": {
    "netConfigServeLinkParameter": {
      "type": "object",
      "properties": {
        "sendPort": {
          "$ref": "#/definitions/netConfigServePortParameter"
        },
        "recvPort": {
          "$ref": "#/definitions/netConfigServePortParameter"
        },
        "linkState": {
          "type": "integer",
          "format": "int64"
        }
      }
    },
    "netConfigServeLinks": {
      "type": "object",
      "properties": {
        "lp": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/netConfigServeLinkParameter"
          }
        }
      }
    },
    "netConfigServePortParameter": {
      "type": "object",
      "properties": {
        "rioId": {
          "type": "integer",
          "format": "int32"
        },
        "appName": {
          "type": "string"
        },
        "portName": {
          "type": "string"
        },
        "slotSize": {
          "type": "integer",
          "format": "int32"
        },
        "type": {
          "$ref": "#/definitions/netConfigServeportType"
        },
        "remoteAppName": {
          "type": "string"
        },
        "remotePortName": {
          "type": "string"
        }
      }
    },
    "netConfigServeportType": {
      "type": "string",
      "enum": [
        "SENDPORT",
        "RECVPORT",
        "DEPIPORT"
      ],
      "default": "SENDPORT"
    }
  }
}
