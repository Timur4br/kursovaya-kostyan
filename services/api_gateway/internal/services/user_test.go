package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/dto"
	"github.com/Timur4br/kursovaya-kostyan/services/api_gateway/internal/models"
)

type userDAOStub struct {
	calls  int
	values map[string]any
}

func (d *userDAOStub) CreateUser(_ context.Context, v map[string]any) (models.User, error) {
	d.calls++
	d.values = v
	return models.User{ID: 1}, nil
}
func (d *userDAOStub) GetUser(context.Context, int64) (models.User, error) {
	d.calls++
	return models.User{ID: 1}, nil
}
func (d *userDAOStub) UpdateUser(_ context.Context, _ int64, v map[string]any) (models.User, error) {
	d.calls++
	d.values = v
	return models.User{ID: 1}, nil
}
func (d *userDAOStub) DeleteUser(context.Context, int64) error { d.calls++; return nil }

func TestUserServiceValidatesBeforeDAO(t *testing.T) {
	for _, body := range []string{`{}`, `{"username":"x"}`, `{"username":"x","name":"  "}`, `{"username":null,"name":"X"}`} {
		d := &userDAOStub{}
		service := NewUserService(d)
		var input dto.UserInput
		if err := json.Unmarshal([]byte(body), &input); err != nil {
			t.Fatal(err)
		}
		if _, err := service.Create(context.Background(), input); err == nil {
			t.Fatal("invalid create accepted")
		}
		if _, err := service.Update(context.Background(), 1, input); err == nil {
			t.Fatal("partial PUT accepted")
		}
		if d.calls != 0 {
			t.Fatal("invalid data reached DAO")
		}
	}
	d := &userDAOStub{}
	service := NewUserService(d)
	if _, err := service.Get(context.Background(), 0); err == nil {
		t.Fatal("zero ID accepted")
	}
	if err := service.Delete(context.Background(), -1); err == nil {
		t.Fatal("negative ID accepted")
	}
	if d.calls != 0 {
		t.Fatal("invalid key reached DAO")
	}
	var input dto.UserInput
	if err := json.Unmarshal([]byte(`{"username":"owner","name":"Owner"}`), &input); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Update(context.Background(), 1, input); err != nil {
		t.Fatal(err)
	}
	if d.calls != 1 || d.values["username"] != "owner" || d.values["name"] != "Owner" {
		t.Fatal("validated values lost", d.values)
	}
}
