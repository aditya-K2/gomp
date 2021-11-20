package cache

import "testing"

func TestLoadCache(t *testing.T) {
	expectedResult := [2]string{"hello/wer.jpg", "hello/iwer.jpg"}
	LoadCache("./testdata/cache.txt")
	var i int = 0
	for _, v := range CACHE_LIST {
		if v != expectedResult[i] {
			if v != expectedResult[i+1] {
				t.Errorf("Didn't Get The Expected Value receieved %s", v)
			}
		}
	}
}
