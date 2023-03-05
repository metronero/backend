package daemon

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func CreateSecret(username string) (secret string, err error) {
	_, secret, err = Config.JwtSecret.Encode(map[string]interface{}{"username": username})
	return
}

func UserRegister(ctx context.Context, username, password, walletAddress string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return pdbExec(ctx,
	    "INSERT INTO accounts(username,password,wallet_address,commission,template)"+
	    "VALUES($1,$2,$3,(SELECT standard_commission FROM settings),$4)",
	    username, hash, walletAddress, defaultTemplate)
}

func UserLogin(ctx context.Context, username, password string) error {
	var hash string
	row, err := pdbQueryRow(ctx, "SELECT password FROM accounts WHERE username=$1", username)
	if err != nil {
		return err
	}
	if err := row.Scan(&hash); err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
