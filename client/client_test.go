package shellgame

import (
	"encoding/json"
	"fmt"
	"github.com/taise-hub/shellgame-cli/common"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPostProfile(t *testing.T) {
	//NOTE: APIのレスポンスの仕様が固ってないためとりあえずtext/plainを返す。
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := &common.Profile{}
		err := json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(p.Name) > 20 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "20文字以内で入力してください。")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
		return
	}))
	defer ts.Close()

	type args struct {
		name string
	}
	tests := map[string]struct {
		args        args
		expectedErr error
	}{
		"ステータスコードが200の時、エラーを起こさずに終了できる。": {
			args: args{
				name: "bob",
			},
			expectedErr: nil,
		},
		"ステータスコードが200以外の時、エラーを取得することができる。": {
			args: args{
				name: "012345678901234567890",
			},
			expectedErr: fmt.Errorf("%s", "20文字以内で入力してください。"),
		},
	}

	for tName, tt := range tests {
		t.Run(tName, func(t *testing.T) {
			url, _ := url.Parse(ts.URL)
			profileEndpoint = url
			actual := PostProfile(tt.args.name)
			if actual != nil {
				if actual.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected: %v\n\t\t Actual: %v \n", tt.expectedErr, actual)
				}
			}
		})
	}
}

func TestGetPlayers(t *testing.T) {
	//NOTE: APIのレスポンスの仕様が固ってないためとりあえずtext/plainを返す。
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		players := []*common.MatchingPlayer{
			{
				Profile: &common.Profile{
					ID:   "da3fc9dd-bff1-43ed-b360-91e4f4ee9db1",
					Name: "Bob",
				},
				Status: common.WAITING,
			},
			{
				Profile: &common.Profile{
					ID:   "22166795-397e-4a16-ad7e-f63bc8cc9222",
					Name: "Alice",
				},
				Status: common.WAITING,
			},
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		bodyBytes, err := json.Marshal(players)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal Server Error")
			return
		}
		return
	}))
	defer ts.Close()

	tests := map[string]struct {
		expectedErr error
	}{
		"ステータスコードが200の時、model.Playersの配列を返すことができる。": {
			expectedErr: nil,
		},
	}
	for tName, tt := range tests {
		t.Run(tName, func(t *testing.T) {
			url, _ := url.Parse(ts.URL)
			playersEndpoint = url

			actual, err := GetMatchingPlayers()
			if err != nil {
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("Expected: %v\n\t\t Actual: %v \n", tt.expectedErr, actual)
				}
			}
			for _, v := range actual {
				fmt.Printf("%#v\n%#v\n", v, v.Profile)
			}
		})
	}

}
