module mama/cmd/mamasrv

require (
	github.com/gorilla/mux v1.8.0 // indirect
	internal/httpserver v1.0.0
	internal/basicauth v1.0.0
	internal/apiv1 v1.0.0
)

replace internal/httpserver => ../../internal/httpserver
replace internal/basicauth => ../../internal/basicauth
replace internal/apiv1 => ../../internal/apiv1

go 1.15
