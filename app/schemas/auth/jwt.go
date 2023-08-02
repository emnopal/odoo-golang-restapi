package schemas

type JWTAccessClaims struct {
	ID       string
	Username string
}

type JWTAccessClaimsJSON struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
