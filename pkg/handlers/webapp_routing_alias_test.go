package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebappRoutingAlias(t *testing.T) {
	// Test that both template names work for the webapp routing template
	
	// Test original name
	appRoutingTemplate, err := GetTemplate("app-routing-ingress", "", ".", nil)
	assert.NoError(t, err)
	assert.NotNil(t, appRoutingTemplate)
	assert.Equal(t, "app-routing-ingress", appRoutingTemplate.Config.TemplateName)
	
	// Test alias name
	webappRoutingTemplate, err := GetTemplate("webapp_routing", "", ".", nil)
	assert.NoError(t, err)
	assert.NotNil(t, webappRoutingTemplate)
	assert.Equal(t, "app-routing-ingress", webappRoutingTemplate.Config.TemplateName)
	
	// Verify they are the same template
	assert.Equal(t, appRoutingTemplate.Config.TemplateName, webappRoutingTemplate.Config.TemplateName)
	assert.Equal(t, appRoutingTemplate.Config.Type, webappRoutingTemplate.Config.Type)
	assert.Equal(t, len(appRoutingTemplate.Config.Variables), len(webappRoutingTemplate.Config.Variables))
}

func TestIsValidTemplateWithAlias(t *testing.T) {
	// Test that IsValidTemplate recognizes both names
	assert.True(t, IsValidTemplate("app-routing-ingress"))
	assert.True(t, IsValidTemplate("webapp_routing"))
	assert.False(t, IsValidTemplate("nonexistent-template"))
}