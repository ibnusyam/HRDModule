package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Gunakan struct Claims yang SAMA STRUKTURNYA dengan Auth Module
// agar bisa membaca data user (SiteID, DeptID, dll)
type HRDClaims struct {
    SiteID string `json:"SiteID"`
    DeptID string `json:"DeptID"`
    PosID  string `json:"PosID"`
    jwt.RegisteredClaims
}

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // 1. Ambil Token dari Header Authorization
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "message": "Missing Authorization Header",
            })
        }

        // Format harus "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            return c.JSON(http.StatusUnauthorized, map[string]string{
                "message": "Invalid Authorization Header Format",
            })
        }
        tokenString := parts[1]

        // 2. Parse dan Validasi Token
        token, err := jwt.ParseWithClaims(tokenString, &HRDClaims{}, func(token *jwt.Token) (interface{}, error) {
            // Validasi Algoritma Signing
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            // KUNCI UTAMA: Secret harus sama dengan Auth Module
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        // 3. LOGIC PENOLAKAN (REJECTION)
        // Jika error atau token tidak valid (termasuk expired)
        if err != nil || !token.Valid {
            // Ini response penolakan standar agar Frontend tahu harus minta Refresh Token
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{
                "message": "Unauthorized: Token expired or invalid",
                "error":   err.Error(), // Opsional: tampilkan error detail untuk debug
            })
        }

        // 4. Jika Valid, Set User Data ke Context
        // Agar bisa dipakai di Handler (misal: mau ambil DeptID user yg login)
        if claims, ok := token.Claims.(*HRDClaims); ok {
            c.Set("user", claims)
        }

        return next(c)
    }
}