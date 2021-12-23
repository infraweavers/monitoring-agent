package web

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/lxn/win"
)

func getResult(counter string) CounterResult {
	counterValue, _ := ReadPerformanceCounter(counter)

	return counterValue
}

func ReadPerformanceCounter(counter string) (CounterResult, error) {

	returnvalue := CounterResult{}

	var queryHandle win.PDH_HQUERY
	var counterHandle win.PDH_HCOUNTER

	ret := win.PdhOpenQuery(0, 0, &queryHandle)
	if ret != win.ERROR_SUCCESS {
		return returnvalue, errors.New("unable to open query through dll call")
	}

	// test path
	ret = win.PdhValidatePath(counter)
	if ret == win.PDH_CSTATUS_BAD_COUNTERNAME {
		return returnvalue, errors.New("unable to fetch counter (this is unexpected)")
	}

	ret = win.PdhAddEnglishCounter(queryHandle, counter, 0, &counterHandle)
	if ret != win.ERROR_SUCCESS {
		return returnvalue, fmt.Errorf("unable to add process counter. error code is %x", ret)
	}

	ret = win.PdhCollectQueryData(queryHandle)
	if ret != win.ERROR_SUCCESS {
		return returnvalue, fmt.Errorf("got an error: 0x%x", ret)
	}

	ret = win.PdhCollectQueryData(queryHandle)
	if ret == win.ERROR_SUCCESS {

		var bufSize uint32
		var bufCount uint32
		var size uint32 = uint32(unsafe.Sizeof(win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE{}))
		var emptyBuf [1]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.

		ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &emptyBuf[0])
		if ret == win.PDH_MORE_DATA {
			filledBuf := make([]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE, bufCount*size)
			ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &filledBuf[0])
			if ret == win.ERROR_SUCCESS {
				for i := 0; i < int(bufCount); i++ {
					c := filledBuf[i]
					s := win.UTF16PtrToString(c.SzName)

					metricName := counter
					if len(s) > 0 {
						metricName = fmt.Sprintf("%s.%s", counter, s)
					}

					returnvalue.Name = metricName
					returnvalue.Value = fmt.Sprintf("%f", c.FmtValue.DoubleValue)
				}
			}
		}
	}

	return returnvalue, nil

}
