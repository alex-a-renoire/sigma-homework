package httphandler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/pkg/storage"
	"github.com/alex-a-renoire/sigma-homework/service"
)

func TestHttpHandler(t *testing.T) {
	type fields struct {
		s storage.Storage
	}

	type args struct {
		url    string
		method string
		header string
		body   []byte
	}

	type resp struct {
		code int
		body string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   resp
	}{
		{
			name: "POST /persons OK",
			fields: fields{
				s: storage.MockStorage{
					MockAddPerson: func(_ string) (int, error) {
						return 1, nil
					},
				},
			},
			args: args{
				url:    "/persons",
				method: "POST",
				body:   []byte(`{"name":"John"}`),
			},
			want: resp{code: http.StatusCreated},
		},
		{
			name: "POST /persons no name",
			fields: fields{
				s: storage.MockStorage{
					MockAddPerson: func(_ string) (int, error) {
						return 0, fmt.Errorf("failed to add person to db")
					},
				},
			},
			args: args{
				url:    "/persons",
				method: "POST",
				body:   []byte(`{"blahblah":"John"}`),
			},
			want: resp{
				code: http.StatusBadRequest,
				body: "{\"message\":\"failed to add person to db\"}",
			},
		},
		{
			name: "GET all /persons",
			fields: fields{
				s: storage.MockStorage{
					MockGetAllPersons: func() ([]model.Person, error) {
						return []model.Person{{
							Id:   1,
							Name: "John",
						}, {
							Id:   2,
							Name: "Jane",
						},
						}, nil
					},
				},
			},
			args: args{
				url:    "/persons",
				method: "GET",
			},
			want: resp{
				code: http.StatusOK,
				body: "[{\"id\":1,\"name\":\"John\"},{\"id\":2,\"name\":\"Jane\"}]",
			},
		},
		{
			name: "Get /persons/1",
			fields: fields{
				s: storage.MockStorage{
					MockGetPerson: func(_ int) (model.Person, error) {
						return model.Person{
							Id:   1,
							Name: "John",
						}, nil
					},
				},
			},
			args: args{
				url:    "/persons/1",
				method: "GET",
			},
			want: resp{
				code: http.StatusOK,
				body: "{\"id\":1,\"name\":\"John\"}",
			},
		},
		{
			name: "UPDATE /persons/1 OK",
			fields: fields{
				s: storage.MockStorage{
					MockUpdatePerson: func(_ int, _ string) (model.Person, error) {
						return model.Person{
							Id:   1,
							Name: "Jane",
						}, nil
					},
				},
			},
			args: args{
				url:    "/persons/1",
				method: "PATCH",
			},
			want: resp{
				code: http.StatusOK,
				body: "{\"id\":1,\"name\":\"Jane\"}",
			},
		},
		{
			name: "DELETE /persons/1 OK",
			fields: fields{
				s: storage.MockStorage{
					MockDeletePerson: func(_ int) error {
						return nil
					},
				},
			},
			args: args{
				url:    "/persons/1",
				method: "DELETE",
			},
			want: resp{
				code: http.StatusOK,
			},
		},
		{
			name: "DELETE /persons/1 !OK",
			fields: fields{
				s: storage.MockStorage{
					MockDeletePerson: func(_ int) error {
						return fmt.Errorf("failed to delete")
					},
				},
			},
			args: args{
				url:    "/persons/1",
				method: "DELETE",
			},
			want: resp{
				code: http.StatusInternalServerError,
				body: "{\"message\":\"failed to delete\"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := service.NewDirect(tt.fields.s)

			s := &HTTPHandler{
				service: srv,
			}

			hs := httptest.NewServer(s.GetRouter())
			defer hs.Close()

			cl := hs.Client()
			req, _ := http.NewRequest(tt.args.method, hs.URL+tt.args.url, bytes.NewReader(tt.args.body))

			r, err := cl.Do(req)

			data, _ := ioutil.ReadAll(r.Body)
			respBody := string(data)

			if err != nil || r.StatusCode != tt.want.code || (tt.want.body != "" && respBody != tt.want.body) {
				if err != nil {
					t.Errorf("error: %s", err)
				} else if respBody != tt.want.body {
					t.Errorf("%s %s = %v, want %v", tt.args.method, tt.args.url, respBody, tt.want.body)
				} else {
					t.Errorf("%s %s = %v, want %v", tt.args.method, tt.args.url, r.StatusCode, tt.want.code)
				}
			}
		})
	}
}
