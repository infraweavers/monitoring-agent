package web

type OSSpecificRequest struct{}
type OSSpecificResult struct{}

func getResult(osspecificrequest OSSpecificRequest) (OSSpecificResult, error) {
	returnvalue := OSSpecificResult{}

	return returnvalue, nil
}
