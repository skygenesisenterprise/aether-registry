package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/skygenesisenterprise/aether-registry/server/src/services"

	"github.com/gin-gonic/gin"
)

type ServiceKeyMiddleware struct {
	serviceKeyService *services.ServiceKeyService
}

func NewServiceKeyMiddleware(serviceKeyService *services.ServiceKeyService) *ServiceKeyMiddleware {
	return &ServiceKeyMiddleware{
		serviceKeyService: serviceKeyService,
	}
}

func (m *ServiceKeyMiddleware) RequireServiceKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		serviceKeyHeader := ctx.GetHeader("X-Service-Key")

		var serviceKey string

		if serviceKeyHeader != "" {
			serviceKey = serviceKeyHeader
		} else if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				serviceKey = strings.TrimPrefix(authHeader, "Bearer ")
			} else if strings.HasPrefix(authHeader, "sk_") {
				serviceKey = authHeader
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Invalid authorization format. Use 'Bearer sk_xxx' or 'sk_xxx' or X-Service-Key header",
					"code":    "INVALID_AUTH_FORMAT",
				})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header or X-Service-Key required",
				"code":    "MISSING_AUTH_HEADER",
			})
			ctx.Abort()
			return
		}

		key, err := m.serviceKeyService.ValidateKey(serviceKey)
		if err != nil {
			fmt.Printf("[service key middleware] Validation error: %v\n", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired service key",
				"code":    "INVALID_SERVICE_KEY",
			})
			ctx.Abort()
			return
		}

		ctx.Set("serviceKeyID", key.ID)
		ctx.Set("serviceKeyName", key.Name)
		ctx.Set("serviceKeyScope", key.Scope)

		ctx.Next()
	}
}

func (m *ServiceKeyMiddleware) RequireServiceKeyWithScope(requiredScope string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		serviceKeyHeader := ctx.GetHeader("X-Service-Key")

		var serviceKey string

		if serviceKeyHeader != "" {
			serviceKey = serviceKeyHeader
		} else if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				serviceKey = strings.TrimPrefix(authHeader, "Bearer ")
			} else if strings.HasPrefix(authHeader, "sk_") {
				serviceKey = authHeader
			} else {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Invalid authorization format",
					"code":    "INVALID_AUTH_FORMAT",
				})
				ctx.Abort()
				return
			}
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header or X-Service-Key required",
				"code":    "MISSING_AUTH_HEADER",
			})
			ctx.Abort()
			return
		}

		key, err := m.serviceKeyService.ValidateKey(serviceKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired service key",
				"code":    "INVALID_SERVICE_KEY",
			})
			ctx.Abort()
			return
		}

		scopes := strings.Split(key.Scope, ",")
		hasScope := false
		for _, scope := range scopes {
			if strings.TrimSpace(scope) == requiredScope || strings.TrimSpace(scope) == "admin" {
				hasScope = true
				break
			}
		}

		if !hasScope {
			ctx.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   fmt.Sprintf("Service key does not have required scope: %s", requiredScope),
				"code":    "INSUFFICIENT_SCOPE",
			})
			ctx.Abort()
			return
		}

		ctx.Set("serviceKeyID", key.ID)
		ctx.Set("serviceKeyName", key.Name)
		ctx.Set("serviceKeyScope", key.Scope)

		ctx.Next()
	}
}

func (m *ServiceKeyMiddleware) OptionalServiceKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}

		var serviceKey string
		if strings.HasPrefix(authHeader, "Bearer ") {
			serviceKey = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.HasPrefix(authHeader, "sk_") {
			serviceKey = authHeader
		} else {
			ctx.Next()
			return
		}

		key, err := m.serviceKeyService.ValidateKey(serviceKey)
		if err != nil {
			ctx.Next()
			return
		}

		ctx.Set("serviceKeyID", key.ID)
		ctx.Set("serviceKeyName", key.Name)
		ctx.Set("serviceKeyScope", key.Scope)

		ctx.Next()
	}
}
