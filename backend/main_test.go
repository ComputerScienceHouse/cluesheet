package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	muxtrace "github.com/DataDog/dd-trace-go/contrib/gorilla/mux/v2"
	pgxtrace "github.com/DataDog/dd-trace-go/contrib/jackc/pgx.v5/v2"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"csh/cluesheet/config"
)

type postClueParams struct {
	Description  string
	Creator      string
	Tags         []string
	ParentClueId *uuid.UUID `json:"parent_clue_id"`
}

type postClueResult struct {
	Clue         Clue
	ClueRelation *ClueRelation
}

// newTestServer creates a test HTTP server with all routes registered and returns its URL.
func newTestServer(t *testing.T) string {
	t.Helper()
	ctx := t.Context()
	ctx = config.ContextWithConfig(ctx, config.GetConfig(ctx))
	conn, err := pgxtrace.NewPool(ctx, connStr)
	require.NoError(t, err, "failed to connect to postgres")
	t.Cleanup(conn.Close)
	router := muxtrace.NewRouter()
	registerRoutes(ctx, router, conn)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	return server.URL
}

// createCluesheet POSTs a cluesheet and returns the created resource.
func createCluesheet(t *testing.T, serverURL string, params PostCluesheetParams) Cluesheet {
	t.Helper()
	buf, err := json.Marshal(params)
	require.NoError(t, err)
	resp, err := http.Post(serverURL+"/api/v1/cluesheet", "application/json", bytes.NewReader(buf))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var cluesheet Cluesheet
	require.NoError(t, json.Unmarshal(body, &cluesheet))
	return cluesheet
}

// createClue POSTs a clue onto a cluesheet and returns the created resource.
func createClue(t *testing.T, serverURL string, cluesheetID uuid.UUID, params postClueParams) postClueResult {
	t.Helper()
	buf, err := json.Marshal(params)
	require.NoError(t, err)
	resp, err := http.Post(serverURL+"/api/v1/cluesheet/"+cluesheetID.String()+"/clue", "application/json", bytes.NewReader(buf))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var result postClueResult
	require.NoError(t, json.Unmarshal(body, &result))
	return result
}

// TestPostCluesheet verifies that creating a cluesheet returns the correct fields.
func TestPostCluesheet(t *testing.T) {
	serverURL := newTestServer(t)
	params := PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	}

	buf, err := json.Marshal(params)
	assert.NoError(t, err)

	resp, err := http.Post(serverURL+"/api/v1/cluesheet", "application/json", bytes.NewReader(buf))
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var cluesheet Cluesheet
	assert.NoError(t, json.Unmarshal(body, &cluesheet))

	assert.Equal(t, params.Name, cluesheet.Name)
	assert.Nil(t, cluesheet.Origin_id)
	assert.Equal(t, params.Creator, cluesheet.Created_by)
	assert.Equal(t, params.Owners, cluesheet.Owners)
	assert.Equal(t, params.Groups, cluesheet.Groups)
	assert.Equal(t, "hidden", cluesheet.Visibility)
	assert.Nil(t, cluesheet.Clues)
}

// TestGetCluesheet verifies that a created cluesheet can be retrieved by ID.
func TestGetCluesheet(t *testing.T) {
	serverURL := newTestServer(t)
	params := PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	}
	cs := createCluesheet(t, serverURL, params)

	resp, err := http.Get(serverURL + "/api/v1/cluesheet/" + cs.Id.String())
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var cluesheet Cluesheet
	assert.NoError(t, json.Unmarshal(body, &cluesheet))

	assert.Equal(t, cs.Id, cluesheet.Id)
	assert.Equal(t, params.Name, cluesheet.Name)
	assert.Nil(t, cluesheet.Origin_id)
	assert.Equal(t, params.Creator, cluesheet.Created_by)
	assert.Equal(t, params.Owners, cluesheet.Owners)
	assert.Equal(t, params.Groups, cluesheet.Groups)
	assert.Equal(t, "hidden", cluesheet.Visibility)
	assert.Nil(t, cluesheet.Clues)
}

// TestPostClue verifies that adding a clue to a cluesheet returns the correct fields.
func TestPostClue(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})

	clueParams := postClueParams{
		Description: "find the hidden stash",
		Creator:     "mom",
		Tags:        []string{"hard", "outdoors"},
	}

	buf, err := json.Marshal(clueParams)
	assert.NoError(t, err)

	resp, err := http.Post(serverURL+"/api/v1/cluesheet/"+cs.Id.String()+"/clue", "application/json", bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result postClueResult
	assert.NoError(t, json.Unmarshal(body, &result))

	assert.Equal(t, clueParams.Description, result.Clue.Description)
	assert.Equal(t, clueParams.Creator, result.Clue.Created_by)
	assert.Equal(t, clueParams.Tags, result.Clue.Tags)
	assert.Nil(t, result.ClueRelation)
}

// TestPostChildClue verifies that a clue can be added as a child of another clue.
func TestPostChildClue(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})
	parent := createClue(t, serverURL, cs.Id, postClueParams{
		Description: "find the hidden stash",
		Creator:     "mom",
		Tags:        []string{"hard", "outdoors"},
	})

	childParams := postClueParams{
		Description:  "child clue",
		Creator:      "mom",
		Tags:         []string{},
		ParentClueId: &parent.Clue.Id,
	}

	buf, err := json.Marshal(childParams)
	assert.NoError(t, err)

	resp, err := http.Post(serverURL+"/api/v1/cluesheet/"+cs.Id.String()+"/clue", "application/json", bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var result postClueResult
	assert.NoError(t, json.Unmarshal(body, &result))

	assert.Equal(t, childParams.Description, result.Clue.Description)
	assert.NotNil(t, result.ClueRelation)
	assert.Equal(t, parent.Clue.Id, result.ClueRelation.Parent_id)
	assert.Equal(t, result.Clue.Id, result.ClueRelation.Child_id)
}

// TestGetCluesheetHasClues verifies that clues added to a sheet are returned nested in the GET response.
func TestGetCluesheetHasClues(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})
	root := createClue(t, serverURL, cs.Id, postClueParams{
		Description: "find the hidden stash",
		Creator:     "mom",
		Tags:        []string{"hard", "outdoors"},
	})
	child := createClue(t, serverURL, cs.Id, postClueParams{
		Description:  "child clue",
		Creator:      "mom",
		Tags:         []string{},
		ParentClueId: &root.Clue.Id,
	})

	resp, err := http.Get(serverURL + "/api/v1/cluesheet/" + cs.Id.String())
	assert.NoError(t, err)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var cluesheet Cluesheet
	assert.NoError(t, json.Unmarshal(body, &cluesheet))

	assert.NotNil(t, cluesheet.Clues)
	assert.Equal(t, 1, len(*cluesheet.Clues), "expected exactly one root clue")
	assert.Equal(t, root.Clue.Id, (*cluesheet.Clues)[0].Id)
	assert.Equal(t, 1, len((*cluesheet.Clues)[0].Children), "expected one child clue")
	assert.Equal(t, child.Clue.Id, (*cluesheet.Clues)[0].Children[0].Id)
}

// TestGetClueProgressDefault verifies that progress defaults to zero for a user with no recorded completions.
func TestGetClueProgressDefault(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})
	cl := createClue(t, serverURL, cs.Id, postClueParams{
		Description: "find the hidden stash",
		Creator:     "mom",
		Tags:        []string{"hard", "outdoors"},
	})

	url := serverURL + "/api/v1/cluesheet/" + cs.Id.String() + "/clue/" + cl.Clue.Id.String() + "/progress/testuser"
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var progress UserProgress
	assert.NoError(t, json.Unmarshal(body, &progress))

	assert.Equal(t, 0, progress.Completions)
	assert.Equal(t, cl.Clue.Id, progress.Clue_id)
	assert.Equal(t, "testuser", progress.Ipa_uid)
}

// TestPostClueProgress verifies that posting a completion count persists and returns the updated progress.
func TestPostClueProgress(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})
	cl := createClue(t, serverURL, cs.Id, postClueParams{
		Description: "find the hidden stash",
		Creator:     "mom",
		Tags:        []string{"hard", "outdoors"},
	})

	url := serverURL + "/api/v1/cluesheet/" + cs.Id.String() + "/clue/" + cl.Clue.Id.String() + "/progress/testuser"
	resp, err := http.Post(url, "text/plain", strings.NewReader("3"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var progress UserProgress
	assert.NoError(t, json.Unmarshal(body, &progress))

	assert.Equal(t, 3, progress.Completions)
	assert.Equal(t, cl.Clue.Id, progress.Clue_id)
	assert.Equal(t, "testuser", progress.Ipa_uid)
}

// TestGetClueProgressAfterPost verifies that posted progress is readable back via GET.
func TestGetClueProgressAfterPost(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})
	cl := createClue(t, serverURL, cs.Id, postClueParams{
		Description: "find the hidden stash",
		Creator:     "mom",
		Tags:        []string{"hard", "outdoors"},
	})

	url := serverURL + "/api/v1/cluesheet/" + cs.Id.String() + "/clue/" + cl.Clue.Id.String() + "/progress/testuser"

	setupResp, err := http.Post(url, "text/plain", strings.NewReader("3"))
	require.NoError(t, err, "setup: failed to post clue progress")
	require.Equal(t, http.StatusOK, setupResp.StatusCode, "setup: unexpected status posting clue progress")

	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var progress UserProgress
	assert.NoError(t, json.Unmarshal(body, &progress))

	assert.Equal(t, 3, progress.Completions)
}

// TestGetParticipationDefault verifies that participation defaults to not hidden for a new user.
func TestGetParticipationDefault(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})

	url := serverURL + "/api/v1/cluesheet/" + cs.Id.String() + "/participation/testuser"
	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var participation UserParticipation
	assert.NoError(t, json.Unmarshal(body, &participation))

	assert.False(t, participation.Hidden)
	assert.Equal(t, cs.Id, participation.Cluesheet_id)
	assert.Equal(t, "testuser", participation.Ipa_uid)
}

// TestPostParticipation verifies that a user's hidden status can be set via POST.
func TestPostParticipation(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})

	url := serverURL + "/api/v1/cluesheet/" + cs.Id.String() + "/participation/testuser"
	buf, err := json.Marshal(struct{ Hidden bool }{Hidden: true})
	assert.NoError(t, err)

	resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var participation UserParticipation
	assert.NoError(t, json.Unmarshal(body, &participation))

	assert.True(t, participation.Hidden)
	assert.Equal(t, cs.Id, participation.Cluesheet_id)
	assert.Equal(t, "testuser", participation.Ipa_uid)
}

// TestGetParticipationAfterPost verifies that a posted hidden status is readable back via GET.
func TestGetParticipationAfterPost(t *testing.T) {
	serverURL := newTestServer(t)
	cs := createCluesheet(t, serverURL, PostCluesheetParams{
		Name:    "foo",
		Creator: "mom",
		Owners:  []string{"mom", "willard"},
		Groups:  []string{"rtp"},
	})

	url := serverURL + "/api/v1/cluesheet/" + cs.Id.String() + "/participation/testuser"

	buf, err := json.Marshal(struct{ Hidden bool }{Hidden: true})
	require.NoError(t, err)
	setupResp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	require.NoError(t, err, "setup: failed to post participation")
	require.Equal(t, http.StatusOK, setupResp.StatusCode, "setup: unexpected status posting participation")

	resp, err := http.Get(url)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var participation UserParticipation
	assert.NoError(t, json.Unmarshal(body, &participation))

	assert.True(t, participation.Hidden)
}
