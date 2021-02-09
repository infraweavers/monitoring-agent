module mama/internal/httpserver

require internal/basicauth v1.0.0
replace internal/basicauth => ../basicauth
require internal/apiv1 v1.0.0
replace internal/apiv1 => ../apiv1

go 1.15
