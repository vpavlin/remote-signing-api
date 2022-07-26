package handlers

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/sirupsen/logrus"
	"github.com/vpavlin/remote-signing-api/config"
	"github.com/vpavlin/remote-signing-api/internal/wallet"
	"github.com/vpavlin/remote-signing-api/pkg/signer"
)

type SignerHandlers struct{}

func SetupSigner(e *echo.Echo, config *config.Config) {
	g := e.Group("")

	wm, err := wallet.NewWalletManager(nil, config.WalletManager)
	if err != nil {
		logrus.Fatal(err)
	}

	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("WalletManager", wm)
			return next(c)
		}
	})

	g.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(ctx echo.Context) bool {
			if ctx.Request().URL.String() == "/signer/new" {
				return true
			}
			return false
		},
		Validator: func(key string, ctx echo.Context) (bool, error) {
			address := ctx.Param("address")
			wm := ctx.Get("WalletManager").(*wallet.WalletManager)

			w, err := wm.GetByAddress(key, common.HexToAddress(address))
			if err != nil {
				return false, err
			}

			logrus.Debugf("Authenticated signer %s == %s -> %t", w.PublicKey.String(), address, w.PublicKey == common.HexToAddress(address))
			return w.PublicKey == common.HexToAddress(address), nil
		},
	}))

	signer.RegisterHandlers(g, SignerHandlers{})
}

func (sh SignerHandlers) NewSigner(ctx echo.Context) error {
	ns := new(signer.SignerKey)
	err := ctx.Bind(ns)
	if err != nil {
		return err
	}

	wm := ctx.Get("WalletManager").(*wallet.WalletManager)

	w, err := wm.New(*ns.Key)
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := new(signer.NewSigner200)
	pk := w.PublicKey.String()
	resp.PublicKey = &pk

	return ctx.JSON(http.StatusOK, resp)
}

func (sn SignerHandlers) SignBytes(ctx echo.Context, address signer.Address, params signer.SignBytesParams) error {
	bearer := params.Authorization
	apiKey := bearer[7:]

	data := new(signer.SignBytesJSONRequestBody)

	err := ctx.Bind(data)
	if err != nil {
		return err
	}

	publicKey := common.HexToAddress(address)

	wm := ctx.Get("WalletManager").(*wallet.WalletManager)

	wallet, err := wm.GetByAddress(apiKey, publicKey)
	if err != nil {
		return err
	}

	signature, err := wallet.Sign(*data.Bytes)
	if err != nil {
		return err
	}

	rtx := new(signer.SignBytes200)
	signedData := signature
	rtx.SignedData = &signedData

	return ctx.JSON(http.StatusOK, rtx)
}

func (sn SignerHandlers) ReplaceKey(ctx echo.Context, address signer.Address, params signer.ReplaceKeyParams) error {
	bearer := params.Authorization
	apiKey := bearer[7:]

	data := new(signer.ReplaceKeyJSONRequestBody)

	err := ctx.Bind(data)
	if err != nil {
		return err
	}

	wm := ctx.Get("WalletManager").(*wallet.WalletManager)

	publicKey := common.HexToAddress(address)
	logrus.Debugf("Updating address %s", address)

	wallet, err := wm.GetByAddress(apiKey, publicKey)
	if err != nil {
		return err
	}

	err = wallet.ReplaceKey(*data.Key)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, nil)
}

/*func HandleTxSign(ctx echo.Context) error {
	bearer := ctx.Request().Header.Get("Authorization")
	key := bearer[7:]

	log.Println("meh")
	address := ctx.Param("address")
	chainIdS := ctx.Param("chainId")
	tx := new(auth.RestTransaction)
	err := json.NewDecoder(ctx.Request().Body).Decode(tx)
	if err != nil {
		return fmt.Errorf("Could not unmarshal the request body: %s", err)
	}
	tx.Address = common.HexToAddress(address)

	wallet, err := wallet.NewWalletFromStorage(key, tx.Address)
	if err != nil {
		return err
	}

	err = tx.Unmarshal()
	if err != nil {
		return err
	}

	chainId, err := strconv.Atoi(chainIdS)
	if err != nil {
		return err
	}

	signedTx, err := wallet.SignTX(tx.Transaction, big.NewInt(int64(chainId)))
	if err != nil {
		return err
	}

	tx.Transaction = signedTx

	return ctx.JSON(http.StatusOK, signedTx)
}*/
