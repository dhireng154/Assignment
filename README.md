# Assignment
The Alert Management Server is a simple Go application that provides RESTful endpoints for writing and reading alerts. The server stores data in-memory and supports basic alert operations.

The server features an in-memory data store for handling alerts. It exposes endpoints for writing alerts (/alerts) and reading alerts based on service ID and time range (/alerts/service_id={service_id}&start_ts={alert_ts}&end_ts={alert_end_ts}).

##Usage 
go run main.go

## Endpints
POST /alerts: Writes a new alert.
GET /alerts/service_id={service_id}&start_ts={alert_ts}&end_ts={alert_end_ts}: Reads alerts based on service ID and time range.

## Write Alerts
Users should be able to send requests to this API to write alert data to the chosen data storage.

### Request Body:

{
"alert_id": "b950482e9911ec7e41f7ca5e5d9a424f",
"service_id": "my_test_service_id",
"service_name": "my_test_service",
"model": "my_test_model",
"alert_type": "anomaly",
"alert_ts": "1695644160",
"severity": "warning",
"team_slack": "slack_ch"
}

### Success Response Body:

{
"alert_id": "b950482e9911ec7e41f7ca5e5d9a424f",
"error": ""
}

### Error Response body

{
"alert_id": "b950482e9911ec7e41f7ca5e5d9a424f",
"error": "<error details>"
}

## Read Alerts
Users should be able to query alerts using the service_id and specifying a time period defined by
start_ts and end_ts.

### Success Response Body:

{
"service_id" : "my_test_service_id"
"service_name": "my_test_service",
"alerts" : [
{
"alert_id": "b950482e9911ec7e41f7ca5e5d9a424f",
"model": "my_test_model",
"alert_type": "anomaly",
"alert_ts": "1695644060",
"severity": "warning",
"team_slack": "slack_ch"
},
{
"alert_id": "b950482e9911ecsdfs41f75e5d9az23cv",
"model": "my_test_model",
"alert_type": "anomaly",
"alert_ts": "1695644160",
"severity": "warning",
"team_slack": "slack_ch"
},
]
}


### Error Response Body:

{
"alert_id": "b950482e9911ec7e41f7ca5e5d9a424f",
"error": "<error details>"
}