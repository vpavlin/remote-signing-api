package handlers

import (
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo"
	"github.com/vpavlin/remote-signing-api/internal/nonce"
)

type NonceResponse struct {
	Nonce   uint64         `json:"nonce"`
	ChainId uint64         `json:"chainId"`
	Address common.Address `json:"address"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func SeuptNonceGroup(e *echo.Echo) {
	g := e.Group("/nonce")

	g.GET("/:chainId/:address", HandleGetNonce)
	g.PUT("/:chainId/:address/:nonce", HandleReturnNonce)
	g.POST("/:chainId/:address/sync", HandleSync)
}

func HandleGetNonce(ctx echo.Context) error {
	return NonceWrapper(GetNonce, ctx)
}

func HandleReturnNonce(ctx echo.Context) error {
	return NonceWrapper(ReturnNonce, ctx)
}

func HandleSync(ctx echo.Context) error {
	return NonceWrapper(SyncNonce, ctx)
}

func NonceWrapper(fn func(ctx echo.Context, nm *nonce.NonceManager, chainId uint64, address common.Address) error, ctx echo.Context) error {
	addressS := ctx.Param("address")
	chainIdS := ctx.Param("chainId")
	nm := ctx.Get("NonceManager").(*nonce.NonceManager)

	address := common.HexToAddress(addressS)

	chainId, err := strconv.Atoi(chainIdS)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	return fn(ctx, nm, uint64(chainId), address)
}

func GetNonce(ctx echo.Context, nm *nonce.NonceManager, chainId uint64, address common.Address) error {

	resp := new(NonceResponse)

	nonce, err := nm.GetNonce(nonce.ChainID(chainId), nonce.Address(address.String()))
	if err != nil {
		ctx.Error(err)
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	resp.Nonce = nonce
	resp.Address = address
	resp.ChainId = chainId

	return ctx.JSON(http.StatusOK, resp)

}

func ReturnNonce(ctx echo.Context, nm *nonce.NonceManager, chainId uint64, address common.Address) error {
	nonceS := ctx.Param("nonce")

	nonceInt, err := strconv.Atoi(nonceS)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	err = nm.ReturnNonce(uint64(nonceInt), nonce.ChainID(chainId), nonce.Address(address.String()))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

func DecreaseNonce(ctx echo.Context, nm *nonce.NonceManager, chainId uint64, address common.Address) error {
	err := nm.DecreaseNonce(nonce.ChainID(chainId), nonce.Address(address.String()))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}

func SyncNonce(ctx echo.Context, nm *nonce.NonceManager, chainId uint64, address common.Address) error {
	err := nm.Sync(nonce.ChainID(chainId), nonce.Address(address.String()))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()})
	}

	return ctx.NoContent(http.StatusOK)
}
