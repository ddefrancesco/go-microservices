package github

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRepoRequestAsJson(t *testing.T)  {
	request := CreateRepoRequest{
		Name:        "golang intro",
		Description: "Test Repo",
		Homepage:    "https://github.com",
		Private:     true,
		HasIssues:   true,
		HasProjects: true,
		HasWiki:     true,
	}

	// json Marshal prende come input una struct e tenta di costruire un json valido
	bytes, err := json.Marshal(request)
	assert.Nil(t,err)
	assert.NotNil(t,bytes)

	assert.EqualValues(t,`{"name":"golang intro","description":"Test Repo","homepage":"https://github.com","private":true,"has_issues":true,"has_projects":true,"has_wiki":true}`, string(bytes))

	var target CreateRepoRequest
	err = json.Unmarshal(bytes,&target)
	assert.Nil(t,err)
	assert.NotNil(t,target)

	assert.EqualValues(t, Name, Name)
	assert.EqualValues(t, HasIssues, HasIssues)
	
}
