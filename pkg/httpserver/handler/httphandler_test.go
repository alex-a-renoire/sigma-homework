package httphandler

import (
	"bytes"
	"fmt"
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
		//context string
	}

	type resp struct {
		code int
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
			want: resp{code: http.StatusBadRequest},
		},
		{
			name: "GET all /persons",
			fields: fields{
				s: storage.MockStorage{
					MockGetAllPersons: func() ([]model.Person, error) {
						return []model.Person{}, nil
					},
				},
			},
			args: args{
				url:    "/persons",
				method: "GET",
			},
			want: resp{code: http.StatusOK},
		},
		{
			name: "Get /persons/1",
			fields: fields{
				s: storage.MockStorage{
					MockGetPerson: func(_ int) (model.Person, error) {
						return model.Person{}, nil
					},
				},
			},
			args: args{
				url:    "/persons/1",
				method: "GET",
			},
			want: resp{code: http.StatusFound},
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

			if err != nil || r.StatusCode != tt.want.code {
				if err != nil {
					t.Errorf("error: %s", err)
				} else {
					t.Errorf("%s %s = %v, want %v", tt.args.method, tt.args.url, r.StatusCode, tt.want.code)
				}
			}
		})
	}
}
