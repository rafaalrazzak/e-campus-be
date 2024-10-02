package ecampus_test

import (
	"net/http/httptest"
	"testing"

	"ecampus-be/bunapp"
	"ecampus-be/ecampus"
	"ecampus-be/testbed"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun/dbfixture"
)

func TestUserHandler(t *testing.T) {
	_, app := testbed.StartApp(t)
	defer app.Stop()

	fixture := loadFixture(t, app)
	testUser := fixture.MustRow("User.test").(*ecampus.User)

	{
		req := testbed.NewRequest("GET", "/api/users", nil)
		resp := httptest.NewRecorder()

		err := app.Router().ServeHTTPError(resp, req)

		require.NoError(t, err)
		require.Contains(t, resp.Body.String(), testUser.Name)
	}

	{
		req := testbed.NewRequest("GET", "/api/users/1", nil)
		resp := httptest.NewRecorder()

		err := app.Router().ServeHTTPError(resp, req)

		require.NoError(t, err)
		require.Contains(t, resp.Body.String(), testUser.Name)
	}
}

func loadFixture(t *testing.T, app *bunapp.App) *dbfixture.Fixture {
	db := app.DB()
	db.RegisterModel((*ecampus.User)(nil))

	fixture := dbfixture.New(db, dbfixture.WithRecreateTables())
	err := fixture.Load(app.Context(), bunapp.FS(), "fixture/fixture.yml")
	require.NoError(t, err)

	return fixture
}
