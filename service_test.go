package main

// PLEASE BE AWARE:
// This file is only intended to be using under the assumption of
// the service has been deployed at [http://localhost:80] and not
// intended to be testing right away.

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"testing"
	"time"
)

var paths = map[string]string{
	"住宿信息列表": "/accommodations",
	"心理咨询机构列表": "/platforms/psychological",
	"线上医疗平台列表": "/platforms/medical",
	//"/hospital/supplies",
	"第二版的医院需求列表": "/hospital/supplies/v2",
	"武汉在外人员住宿信息": "/people/accommodations",
	"信息看板": "/wiki/stream",
	"社区物资需求列表": "/community/supplies",
}

var client = http.Client{
	Timeout: time.Minute,
}

func TestAll(t *testing.T) {
	for name, path := range paths {
		t.Logf("> now testing %v (%v)", name, path)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost%v", path), nil)
		if err != nil {
			t.Errorf("	failed to initialize request %v (%v): %v", path, name, err)
		}
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("	failed to fetch request %v (%v): %v", path, name, err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("	response not ok for %v (%v): %v", path, name, resp.StatusCode)
		} else {
			t.Logf("	> PASSED %v (%v) is ok (statuscode=%v)", path, name, resp.StatusCode)
		}
	}
}
